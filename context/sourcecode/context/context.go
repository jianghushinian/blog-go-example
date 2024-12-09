package context

import (
	"errors"
	"internal/reflectlite"
	"sync"
	"sync/atomic"
	"time"
)

// Context 可以携带一个截止时间、一个取消信号和一对键值
// Context 的实现是并发安全的，所以它的方法可以被多个 goroutine 同时调用
type Context interface {
	// Deadline 返回该 Context 应该被取消的截止日期
	// 如果未设置截止日期，则返回的 ok 值为 false，多次调用返回结果相同
	Deadline() (deadline time.Time, ok bool)

	// Done 返回一个只读的 channel 作为取消信号，当 Context 被取消，此 channel 会被 close 掉
	// 如果 Context 永远无法取消，可能会返回 nil，多次调用返回结果相同
	// WithCancel 构造的 Context 在调用 cancel() 时关闭 Done channel
	// WithDeadline 构造的 Context 在截止时间到期时关闭 Done channel
	// WithTimeout 构造的 Context 在超时时间到期时关闭 Done channel
	Done() <-chan struct{}

	// Err 如果 Done channel 尚未关闭，Err 返回 nil
	// 如果 Done channel 已关闭，Err 返回一个非 nil 的错误，错误原因如下：
	// 如果 Context 被取消，则返回 Canceled 错误；如果 Context 的截止时间已过，则返回 DeadlineExceeded 错误
	// 在 Err 返回非 nil 错误之后，多次调用返回错误相同
	Err() error

	// Value 返回与给定键（key）关联的值（value），如果没有与该 key 关联的 value，则返回 nil
	// 对于相同的 key，多次调用返回结果相同
	Value(key any) any
}

// Canceled 是在取消 Context 时调用 Context.Err() 返回的错误
var Canceled = errors.New("context canceled")

// DeadlineExceeded 是在当 Context 的最后期限过去时，调用 Context.Err() 返回的错误
var DeadlineExceeded error = deadlineExceededError{}

// 超时的 error 单独定义了一个结构体，并且实现方法，方便调用方对 error 的行为进行断言
// Timeout、Temporary 这两个方法在标准库中的 url.Error 也实现了
// https://jianghushinian.cn/2024/10/01/go-error-guidelines-error-handling/#%E8%A1%8C%E4%B8%BA%E6%96%AD%E8%A8%80
type deadlineExceededError struct{}

func (deadlineExceededError) Error() string   { return "context deadline exceeded" }
func (deadlineExceededError) Timeout() bool   { return true }
func (deadlineExceededError) Temporary() bool { return true }

// emptyCtx 是 Context 接口的最小实现，作为 backgroundCtx 和 todoCtx 的基础
// 一个空的 Context，永远不会被取消，没有值，也没有截止日期
type emptyCtx struct{}

func (emptyCtx) Deadline() (deadline time.Time, ok bool) {
	return
}

func (emptyCtx) Done() <-chan struct{} {
	return nil
}

func (emptyCtx) Err() error {
	return nil
}

func (emptyCtx) Value(key any) any {
	return nil
}

// 包装了 emptyCtx，作为最顶层的 Background Context
type backgroundCtx struct{ emptyCtx }

// 实现了 fmt.Stringer 接口
func (backgroundCtx) String() string {
	return "context.Background"
}

// 包装了 emptyCtx，在没想好用哪个 Context 时，使用 TODO Context
type todoCtx struct{ emptyCtx }

// 实现了 fmt.Stringer 接口
func (todoCtx) String() string {
	return "context.TODO"
}

// Background 返回一个非 nil 的空 Context
// 通常由主函数、初始化和测试使用，并作为传入请求的顶层 Context
func Background() Context {
	return backgroundCtx{} // 返回的是结构体
}

// TODO 返回一个非 nil 的空 Context
// 当不清楚要使用哪个 Context 或它还不可用时（相关的函数尚未扩展支持 Context 参数），应使用 context.TODO
func TODO() Context {
	return todoCtx{} // 返回的是结构体
}

// CancelFunc 取消函数，用于通知一个操作取消或中止
// 不会等待作业停止，可以被多个 goroutine 同时调用，第一次调用后，后续重复调用不产生任何作用
type CancelFunc func()

// WithCancel 根据给定的父 Context 构造一个新的具有取消功能的 Context 并返回
// 返回的 Context 的 Done channel 会在以下两种情况下关闭（以先发生者为准）：
// 调用返回的 cancel() 函数时
// 或者当父 Context 的 Done channel 关闭时
// 取消这个 Context 会释放与其相关的资源，因此应该在该 Context 中的操作完成后尽早调用 cancel() 函数
func WithCancel(parent Context) (ctx Context, cancel CancelFunc) {
	c := withCancel(parent)
	return c, func() { c.cancel(true, Canceled, nil) }
}

// CancelCauseFunc 取消函数，取消时可以传入根因
type CancelCauseFunc func(cause error)

// WithCancelCause 与 WithCancel 类似，但返回 CancelCauseFunc 而不是 CancelFunc
func WithCancelCause(parent Context) (ctx Context, cancel CancelCauseFunc) {
	c := withCancel(parent)
	return c, func(cause error) { c.cancel(true, Canceled, cause) }
}

// 构造带有取消功能的 Context
func withCancel(parent Context) *cancelCtx {
	if parent == nil {
		panic("cannot create context from nil parent")
	}
	c := &cancelCtx{}            // 带取消功能的 Context
	c.propagateCancel(parent, c) // 将新构造的 Context 向上传播挂载到父 Context 的 children 属性中，这样当父 Context 取消时子 Context 对象 c 也会级联取消
	return c
}

// Cause 从 Context 中提取根因，没有则返回 Context 的取消原因 Context.Err()
func Cause(c Context) error {
	if cc, ok := c.Value(&cancelCtxKey).(*cancelCtx); ok {
		cc.mu.Lock()
		defer cc.mu.Unlock()
		return cc.cause
	}
	return c.Err()
}

// AfterFunc 用于在给定的 Context 过期或取消时异步（开启新的 goroutine）执行延迟任务
// 在 Context 的 Done channel 被 close 时，自动执行 f 函数（仅会执行一次）
func AfterFunc(ctx Context, f func()) (stop func() bool) {
	a := &afterFuncCtx{
		f: f,
	}
	// 调用 cancelCtx 的向上传播方法，将 a 的取消功能挂载到父 ctx 的 children 属性中，实现级联取消
	a.cancelCtx.propagateCancel(ctx, a)
	return func() bool { // 返回一个停止函数，用于阻止 f 被执行
		stopped := false
		a.once.Do(func() { // 确保仅执行一次
			stopped = true // 如果此处被执行，则 a.cancel 方法内部的 a.once.Do 就不会重复执行，即阻止 f 被执行
		})
		if stopped { // 第一次调用，取消 Context
			a.cancel(true, Canceled, nil)
		}
		return stopped
	}
}

type afterFuncer interface {
	AfterFunc(func()) func() bool
}

// 能够实现延迟调用给定函数 f 的 Context，包装了 cancelCtx
type afterFuncCtx struct {
	cancelCtx           // “继承”了 cancelCtx
	once      sync.Once // 要么用来开始执行 f，要么用来阻止 f 被执行
	f         func()
}

// afterFuncCtx 的取消函数
func (a *afterFuncCtx) cancel(removeFromParent bool, err, cause error) {
	a.cancelCtx.cancel(false, err, cause) // 取消 cancelCtx
	if removeFromParent {
		removeChild(a.Context, a) // 将当前 *afterFuncCtx 从 cancelCtx 的父 Context 的 children 属性中移除
	}
	a.once.Do(func() { // 确保仅执行一次
		go a.f() // 开启新的 goroutine 执行 f，如果在调用 a.cancel() 之前 stop 函数被调用，stop 函数中的 a.once.Do 优先被执行，则此处就不会执行
	})
}

// stopCtx 是作为 cancelCtx 的父 Context 使用的
// 当一个 AfterFunc 已经在父 Context 中注册（registered）时，它持有用于取消（unregister）AfterFunc 的 stop 函数
type stopCtx struct {
	Context
	stop func() bool
}

// goroutines 统计创建的 goroutine 的数量，用于测试
var goroutines atomic.Int32

// &cancelCtxKey 是 cancelCtx 返回自身时所使用的 key
var cancelCtxKey int

// parentCancelCtx 用于向上查找父 Context 路径中是否存在 *cancelCtx 并返回查找结果
// 通过查找 parent.Value(&cancelCtxKey) 来找到最内层的 *cancelCtx（从下向上查找整条链路，找到第一个）
// 然后检查 parent.Done() 是否与该 *cancelCtx 的 Done channel 匹配
// 如果不匹配，说明 *cancelCtx 已经被包装在一个自定义实现中，提供了一个不同的 done channel，
// 在这种情况下，为了避免跳过这个自定义的 Context 实现，我们不能直接使用原来的 *cancelCtx
func parentCancelCtx(parent Context) (*cancelCtx, bool) {
	done := parent.Done()
	// 如果父 Context 的 Done() 方法返回 closedchan，说明已经被取消
	// 如果返回 nil，说明父 Context 永远不会被取消
	if done == closedchan || done == nil {
		return nil, false
	}
	// 从 Context 路径中自下而上查找 *cancelCtx，传入特殊的 &cancelCtxKey 可以得到 *cancelCtx
	p, ok := parent.Value(&cancelCtxKey).(*cancelCtx)
	if !ok {
		return nil, false
	}
	// 如果找到 *cancelCtx，对 *cancelCtx 进行进一步的检查
	// 确保返回的 *cancelCtx 的 Done channel 与父 Context 的 Done channel 是匹配的
	// 如果匹配，说明 *cancelCtx 是有效的，并且可以继续使用它
	// 否则，说明 *cancelCtx 已经被包装在一个自定义实现中，返回 nil 和 false，表示找不到合适的 *cancelCtx
	pdone, _ := p.done.Load().(chan struct{})
	if pdone != done {
		return nil, false
	}
	return p, true
}

// 从父 Context 的 children 集合中移除子 Context
// 此函数会在 *cancelCtx/*timerCtx/*afterFuncCtx 的 cancel 方法中被调用，即 Context 被取消时调用
func removeChild(parent Context, child canceler) {
	if s, ok := parent.(stopCtx); ok { // 如果父 Context 是 stopCtx 类型，调用其 stop 方法，并返回
		s.stop() // 父 Context 取消时，需要调用 stopCtx.stop() 来取消子 *cancelCtx
		return
	}
	p, ok := parentCancelCtx(parent) // 判断父 Context 是 *cancelCtx 或者从 *cancelCtx 派生而来
	if !ok {                         // 如果不是 *cancelCtx，直接返回
		return
	}
	// 如果是 *cancelCtx，从父 Context 的 children 集合中删除 child
	p.mu.Lock()
	if p.children != nil {
		delete(p.children, child)
	}
	p.mu.Unlock()
}

// canceler 表示一种可以直接取消的 Context 类型，具体实现是 *cancelCtx 和 *timerCtx
type canceler interface {
	cancel(removeFromParent bool, err, cause error) // 取消函数
	Done() <-chan struct{}                          // 通过返回的 channel 能够知道是否被取消
}

// closedchan 表示一个已关闭的 channel
var closedchan = make(chan struct{})

// 导入 context 包时直接关闭 closedchan
func init() {
	close(closedchan)
}

// 带有取消功能的 Context，实现了 canceler 接口
// 取消后，它还会级联取消实现了 canceler 接口的所有 children
type cancelCtx struct {
	Context // “继承”的父 Context

	mu       sync.Mutex            // 持有锁保护下面这些字段
	done     atomic.Value          // 值为 chan struct{} 类型，会被懒惰创建，在第一次调用取消函数 cancel() 时被关闭，表示 Context 已被取消
	children map[canceler]struct{} // 所有可以被取消的子 Context 集合，它们在第一次调用取消函数 cancel() 时被级联取消，然后置为 nil
	err      error                 // 取消原因，在第一次调用取消函数 cancel() 时被设置值
	cause    error                 // 取消根因，在第一次调用取消函数 cancel() 时被设置值
}

// Value 通过给定的 key 查询 *cancelCtx 中对应的 value
func (c *cancelCtx) Value(key any) any {
	// 使用 &cancelCtxKey 标记需要返回自身
	// 这是一个未导出的（unexported）类型，所以仅作为 context 包内部实现的一个“协议”，对用户不可见
	if key == &cancelCtxKey {
		return c
	}
	// 接着向上遍历父 Context 链路，查询 key
	return value(c.Context, key)
}

// Done 实现 Context.Done 方法
func (c *cancelCtx) Done() <-chan struct{} {
	// 使用 double-check 来提升性能
	d := c.done.Load() // 原子操作，比互斥锁更加轻量
	if d != nil {      // 如果存在 channel 直接返回
		return d.(chan struct{})
	}
	c.mu.Lock() // 如果不存在 channel，则要先加锁，然后创建 channel 并返回
	defer c.mu.Unlock()
	d = c.done.Load()
	if d == nil { // 为保证并发安全，再做一次检查
		d = make(chan struct{})
		c.done.Store(d)
	}
	return d.(chan struct{})
}

// Err 实现 Context.Err 方法
func (c *cancelCtx) Err() error {
	c.mu.Lock()
	err := c.err
	c.mu.Unlock()
	return err
}

// propagateCancel 将子 Context 对象 child 向上传播挂载到父 Context 的 children 集合中，这样当父 Context 被取消时，子 Context 也会被级联取消
// 此方法沿着父 Context 路径向上查找，直到找到一个 *cancelCtx（实现了 Done 方法）或者为 nil 停止
// 如果找到 *cancelCtx，就把 child 加入到这个 *cancelCtx 的 children 属性中，以便这个 *cancelCtx 被取消时 child 也会被自动取消
// 如果未找到 *cancelCtx，再判断父 Context 是否为 afterFuncer 类型，如果是，就设置当父 Context 延迟时间到期后，取消子 Context
// 否则，会开起一个 goroutine，由它来监听父 Context 是否被取消（Done channel 将被 close）
func (c *cancelCtx) propagateCancel(parent Context, child canceler) {
	c.Context = parent // “继承”父 Context，这里可以是任何实现了 Context 接口的类型

	// NOTE: 父 Context 没有实现取消功能
	done := parent.Done()
	if done == nil { // 如果父 Context 的 Done() 方法返回 nil，说明父 Context 没有取消的功能，那么无需传播子 Context 的 cancel 功能到父 Context
		return
	}

	// NOTE: 父 Context 已经被取消
	select {
	case <-done: // 直接取消子 Context，且取消原因设置为父 Context 的取消原因
		child.cancel(false, parent.Err(), Cause(parent))
		return
	default:
	}

	// NOTE: 父 Context 还未取消
	if p, ok := parentCancelCtx(parent); ok { // 如果父 Context 是 *cancelCtx 或者从 *cancelCtx 派生而来
		p.mu.Lock()
		if p.err != nil {
			// 如果父 Context 的 err 属性有值，说明已经被取消，直接取消子 Context
			child.cancel(false, p.err, p.cause)
		} else {
			if p.children == nil { // 延迟创建父 Context 的 children 属性
				p.children = make(map[canceler]struct{})
			}
			p.children[child] = struct{}{} // 将 child 加入到这个 *cancelCtx 的 children 集合中
		}
		p.mu.Unlock()
		return
	}

	// NOTE: 父 Context 实现了 afterFuncer 接口
	if a, ok := parent.(afterFuncer); ok { // 测试文件 afterfunc_test.go 中 *afterFuncCtx 实现了 afterFuncer 接口
		c.mu.Lock()
		stop := a.AfterFunc(func() { // 注册子 Context 取消功能到父 Context，当父 Context 取消时，能级联取消子 Context
			child.cancel(false, parent.Err(), Cause(parent))
		})
		c.Context = stopCtx{ // 将当前 *cancelCtx 的直接父 Context 设置为 stopCtx
			Context: parent, // stopCtx 的父 Context 设置为 parent
			stop:    stop,
		}
		c.mu.Unlock()
		return
	}

	// NOTE: 父 Context 不是已知类型，但实现了取消功能
	goroutines.Add(1) // 记录下开启了几个 goroutine，用于测试代码
	go func() {       // 开起一个 goroutine，监听父 Context 是否被取消，如果取消则级联取消子 Context
		select {
		case <-parent.Done(): // 父 Context 被取消
			child.cancel(false, parent.Err(), Cause(parent))
		case <-child.Done(): // 自己被取消
		}
	}()
}

type stringer interface {
	String() string
}

// 获取给定 Context 的字符串名称
func contextName(c Context) string {
	if s, ok := c.(stringer); ok {
		return s.String()
	}
	return reflectlite.TypeOf(c).String()
}

func (c *cancelCtx) String() string {
	return contextName(c.Context) + ".WithCancel"
}

// 取消 Context，关闭 c.done channel，会级联取消 c 的每一个子 Context（c.children）
// 如果 removeFromParent 为 true，它会将 c 从父 Context 的 children 集合中移除
// 如果这是 c 第一次被取消（即第一次调用 cancel），会将 c.cause 设置为 cause
func (c *cancelCtx) cancel(removeFromParent bool, err, cause error) {
	if err == nil {
		panic("context: internal error: missing cancel error")
	}
	if cause == nil { // 如果没有设置根因，取 err
		cause = err
	}
	c.mu.Lock()
	if c.err != nil { // 如果 err 不为空，说明已经被取消，直接返回
		c.mu.Unlock()
		return
	}

	// NOTE: 只有第一次调用 cancel 才会执行之后的代码

	// 记录错误和根因
	c.err = err
	c.cause = cause
	d, _ := c.done.Load().(chan struct{})
	if d == nil { // 如果 done 为空，直接设置一个已关闭的 channel
		c.done.Store(closedchan)
	} else { // 如果 done 有值，将其关闭
		close(d)
	}
	// 级联取消所有子 Context
	for child := range c.children {
		// NOTE: 获取子 Context 的锁，同时持有父 Context 的锁
		child.cancel(false, err, cause)
	}
	c.children = nil // 清空子 Context 集合，因为已经完成了 Context 树整个链路的取消操作
	c.mu.Unlock()

	if removeFromParent { // 从父 Context 的 children 集合中移除当前 Context
		removeChild(c.Context, c)
	}
}

// WithoutCancel 返回一个新的 Context 副本，这个 Context 不会随着父 Context 的取消而取消
// 返回的 Context 不具有截止时间（Deadline）、错误信息（Err），并且它的 Done channel 为 nil，这意味着该 Context 没有与取消相关的功能
// 调用 Cause() 函数从返回的 Context 提取被取消的根因时，返回 nil，因为它不会被取消，所以也就没有取消的根因
func WithoutCancel(parent Context) Context {
	if parent == nil {
		panic("cannot create context from nil parent")
	}
	return withoutCancelCtx{parent}
}

// 没有取消功能的 Context，由此派生的 Context 不会在父 Context 取消时被级联取消
// 仅比 emptyCtx 多实现了一个 Value 方法
type withoutCancelCtx struct {
	c Context
}

func (withoutCancelCtx) Deadline() (deadline time.Time, ok bool) {
	return
}

func (withoutCancelCtx) Done() <-chan struct{} {
	return nil
}

func (withoutCancelCtx) Err() error {
	return nil
}

// Value 虽然没有取消功能，但实现了 Value 方法，可以根据 key 查询 value
func (c withoutCancelCtx) Value(key any) any {
	return value(c, key)
}

func (c withoutCancelCtx) String() string {
	return contextName(c.c) + ".WithoutCancel"
}

// WithDeadline 返回一个新的 Context 副本，并将截止时间设置为不晚于 d
// 如果父 Context 的截止时间比传入的时间 d 更早，WithDeadline(parent, d) 的调用则不会改变父 Context 的行为，相当于直接使用父 Context 本身
// 返回的 Context 的 Done channel 将在截止时间到期时关闭，
// 无论是返回的取消函数被调用，还是父 Context 的 Done channel 被关闭，哪个先发生就哪个先关闭该 Context
// 取消这个 Context 会释放与其相关的资源，因此应该在该 Context 中的操作完成后尽早调用 cancel 函数
func WithDeadline(parent Context, d time.Time) (Context, CancelFunc) {
	return WithDeadlineCause(parent, d, nil)
}

// WithDeadlineCause 与 WithDeadline 类似，但是当截止时间到期时，它还会设置返回的 Context 的根因
// 返回的 Context 将会是 *cancelCtx 或 *timerCtx 类型
func WithDeadlineCause(parent Context, d time.Time, cause error) (Context, CancelFunc) {
	if parent == nil {
		panic("cannot create context from nil parent")
	}
	// 如果父 Context 的截止时间已经比传入的 d 更早，直接返回一个 *cancelCtx（无需构造 *timerCtx 等待定时器判断截止时间到了才取消 Context）
	if cur, ok := parent.Deadline(); ok && cur.Before(d) {
		return WithCancel(parent)
	}
	c := &timerCtx{ // 构造一个带有定时器和截止时间功能的 Context
		deadline: d,
	}
	// 这里使用 cancelCtx 结构体默认值，初始化 timerCtx 时没有显式初始化 cancelCtx 字段
	c.cancelCtx.propagateCancel(parent, c) // 向父 Context 传播 cancel 功能，这样当父 Context 取消时当前 Context 也会被级联取消
	dur := time.Until(d)
	if dur <= 0 { // 截止日期已过，直接取消
		c.cancel(true, DeadlineExceeded, cause)
		return c, func() { c.cancel(false, Canceled, nil) }
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.err == nil {
		c.timer = time.AfterFunc(dur, func() { // 等待截止时间到期，自动调用 cancel 取消 Context
			c.cancel(true, DeadlineExceeded, cause)
		})
	}
	return c, func() { c.cancel(true, Canceled, nil) }
}

// timerCtx 此 Context 的实现关联了一个定时器和截止时间
// 它嵌入了一个 cancelCtx 以实现 Done 和 Err
// 它通过停止定时器然后委托给 cancelCtx.cancel 来实现取消操作
type timerCtx struct {
	cancelCtx             // “继承”了 cancelCtx
	timer     *time.Timer // Under cancelCtx.mu.

	deadline time.Time
}

func (c *timerCtx) Deadline() (deadline time.Time, ok bool) {
	return c.deadline, true
}

func (c *timerCtx) String() string {
	return contextName(c.cancelCtx.Context) + ".WithDeadline(" +
		c.deadline.String() + " [" +
		time.Until(c.deadline).String() + "])"
}

// 取消函数
func (c *timerCtx) cancel(removeFromParent bool, err, cause error) {
	c.cancelCtx.cancel(false, err, cause)
	if removeFromParent {
		// 将此 *timerCtx 从其父 *cancelCtx 的 children 集合中删除
		removeChild(c.cancelCtx.Context, c)
	}
	c.mu.Lock()
	if c.timer != nil {
		c.timer.Stop()
		c.timer = nil
	}
	c.mu.Unlock()
}

// WithTimeout 返回一个调用了 WithDeadline(parent, time.Now().Add(timeout)) 的 Context
// 就是 WithDeadline 的变种，WithDeadline 需要传递截止时间 d time.Time，WithTimeout 则需要传递超时时间 timeout time.Duration
func WithTimeout(parent Context, timeout time.Duration) (Context, CancelFunc) {
	return WithDeadline(parent, time.Now().Add(timeout))
}

// WithTimeoutCause 与 WithTimeout 类似，但是超时时间到期时，它还会设置返回的 Context 的根因
func WithTimeoutCause(parent Context, timeout time.Duration, cause error) (Context, CancelFunc) {
	return WithDeadlineCause(parent, time.Now().Add(timeout), cause)
}

// WithValue 返回一个新的 Context，它复制了父 Context，并在其中将指定的键 key 与值 val 关联，这样，新的 Context 可以附加数据
// WithValue 主要用于携带请求相关的数据，例如跨进程或跨 API 边界传递的信息，不要将它用于传递普通的函数参数或可选参数
// 提供的 key 必须是可比较的类型（comparable），且不应使用内置类型（如 string），使用内置类型可能会导致不同包之间的键冲突，
// 因此，建议使用自定义类型作为键，例如，可以使用一个独立的结构体类型来作为键，
// 如果需要避免内存分配，context 的 key 通常会使用像 struct{} 这样的空结构体类型，空结构体没有数据，内存大小为零，适合作为键，
// 如果 context 的 key 是可导出的（exported）类型，它的静态类型应当是指针类型或接口类型，以避免不同包之间的冲突
func WithValue(parent Context, key, val any) Context {
	if parent == nil {
		panic("cannot create context from nil parent")
	}
	if key == nil {
		panic("nil key")
	}
	if !reflectlite.TypeOf(key).Comparable() {
		panic("key is not comparable")
	}
	return &valueCtx{parent, key, val}
}

// valueCtx 携带一个键值对，它为该 key 实现了 Value 方法，这意味着它能够根据给定的键返回关联的值
// 所有其他方法的调用委托给嵌入的 Context（调用 Deadline、Done、Err 都会直接转发到父 Context）
// 可以并发安全的传递数据，对标其他语言中的 thread-local
type valueCtx struct {
	Context      // 持有父 Context
	key, val any // 存储的键值对，注意一个 Context 仅能保存一对 key/value，这样就能实现并发读的安全，copy-on-write
}

// stringify 尝试在不使用 fmt 包的情况下将 v 转换为字符串
// 由于 fmt 包在格式化过程中可能会涉及到 Unicode 字符集的处理（例如输出表情符号等），
// 而 context 包则不希望引入对 Unicode 表的依赖，所以采用了不依赖 fmt 的方式来实现字符串化
// stringify 函数仅在 *valueCtx.String() 方法中被使用
func stringify(v any) string {
	switch s := v.(type) {
	case stringer: // 实现了 String() 方法，就返回 String() 内容
		return s.String()
	case string: // 字符串类型就返回字符串内容
		return s
	case nil: // nil 返回字符串格式
		return "<nil>"
	}
	// 其他类型会返回对象类型名的字符串格式，而不是对象值的字符串形式
	return reflectlite.TypeOf(v).String()
}

// 代码示例：context.WithValue(context.Background(), "a", 1)
// 输出示例：context.Background.WithValue(a, int)
func (c *valueCtx) String() string {
	// 取父 Context 的 string 形式 + .WithValue(k, v)
	return contextName(c.Context) + ".WithValue(" +
		stringify(c.key) + ", " +
		stringify(c.val) + ")"
}

// Value 实现链式查找，优先从自己的键值对中查找，不存在会向父 Context 中查找
func (c *valueCtx) Value(key any) any {
	if c.key == key { // 在自己的键值对中查找
		return c.val
	}
	return value(c.Context, key) // 沿着父 Context 向上查找
}

// 从下往上遍历 Context 树的一条分支，并从中查找与给定 key 相关联的 value
// 不会遍历整棵树，只会从当前节点，沿着一条链路向上查找父节点，直到找到 value 或到根节点终止
func value(c Context, key any) any {
	for {
		switch ctx := c.(type) { // 断言 Context 类型
		case *valueCtx: // 表示一个用于安全传递数据的 Context
			if key == ctx.key { // 与当前 Context 的 key 匹配，直接返回对应的值 val
				return ctx.val
			}
			c = ctx.Context // key 不匹配，继续向上遍历父 Context
		case *cancelCtx: // 表示一个带有取消功能的 Context
			if key == &cancelCtxKey { // 检查 key 是否等于 &cancelCtxKey（这是一个指向 *cancelCtx 的特殊键），如果匹配，就返回自身（即 c 对象）
				return c
			}
			c = ctx.Context // key 不匹配，继续向上遍历父 Context
		case withoutCancelCtx: // 表示一个不带取消功能的 Context（使用 WithoutCancel() 创建出来的 Context 类型）
			if key == &cancelCtxKey { // 检查 key 是否等于 &cancelCtxKey，如果匹配，说明要查找的是取消信号的特殊键，就返回 nil，因为这种 Context 没有取消信号
				return nil
			}
			c = ctx.c // 如果 key 不匹配，则继续向上遍历父 Context
		case *timerCtx: // 表示一个带有定时器的 Context
			if key == &cancelCtxKey { // 检查 key 是否等于 &cancelCtxKey，如果匹配，返回其包装的 *cancelCtx
				return &ctx.cancelCtx
			}
			c = ctx.Context // key 不匹配，继续向上遍历父 Context
		case backgroundCtx, todoCtx: // 这两个类型是无值的 Context（通常这是 Context 树的根），所以直接返回 nil
			return nil
		default: // 如果没有匹配任何已知的 Context 类型，则调用 Context 的 Value 方法去查找 key 对应的值
			return c.Value(key)
		}
	}
}
