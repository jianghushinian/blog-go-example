package cron

import (
	"context"
	"sort"
	"sync"
	"time"
)

// Cron 核心结构体，用于调度注册进来的作业
// 记录作业列表（entries），可以启动、停止，并且检查运行中的作业状态
type Cron struct {
	entries   []*Entry          // 作业对象列表
	chain     Chain             // 装饰器链
	stop      chan struct{}     // 停止信号
	add       chan *Entry       // Cron 运行时，增加作业的 channel
	remove    chan EntryID      // 移除指定 ID 作业的 channel
	snapshot  chan chan []Entry // 获取当前作业列表快照的 channel
	running   bool              // 标识 Cron 是否正在运行
	logger    Logger            // 日志对象，Cron 会将运行的日志内容输出到 logger
	runningMu sync.Mutex        // 当 Cron 运行时，保护并发操作的锁
	location  *time.Location    // 本地时区，Cron 根据此时区计算任务执行计划
	parser    ScheduleParser    // 任务执行计划解析器
	nextID    EntryID           // 下一个要执行的作业 ID
	jobWaiter sync.WaitGroup    // 使用 wg 等待作业完成
}

// ScheduleParser 此接口用于解析调度规范（spec）并返回一个 Schedule
type ScheduleParser interface {
	Parse(spec string) (Schedule, error)
}

// Job 定义提交的定时作业接口，所有交由执行器 Cron 执行的作业都需要实现此接口
type Job interface {
	Run()
}

// Schedule 描述一个作业的执行计划
type Schedule interface {
	// Next 返回给定时间之后的下一次执行时间
	// Next 第一次开始时被调用，之后每次作业运行时也都会被调用
	Next(time.Time) time.Time
}

// EntryID 标识 Cron 实例中的一个作业
type EntryID int

// Entry 作业实体对象，表示一个被注册进 Cron 执行器中的作业
// 由一个 schedule 和按照该 schedule 执行的 func 组成
type Entry struct {
	// ID 是作业的唯一 ID，可用于查找快照或将其删除
	ID EntryID

	// Schedule 作业的执行计划，应该按照此计划来执行作业
	Schedule Schedule

	// Next 下次运行作业的时间，如果 Cron 尚未启动或无法满足此作业的执行计划，则为 zero time
	Next time.Time

	// Prev 是此作业的最后一次运行时间，如果从未运行，则为 zero time
	Prev time.Time

	// WrappedJob 作业装饰器，为作业增加新的功能，会在 Schedule 被激活时运行
	WrappedJob Job

	// Job 提交到 Cron 中的作业
	Job Job
}

// Valid 校验作业 ID 是否有效，如果不为 0 返回 true
func (e Entry) Valid() bool { return e.ID != 0 }

// 用于按时间对作业列表进行排序的包装器（zero time 会排在末尾）
type byTime []*Entry

func (s byTime) Len() int      { return len(s) }
func (s byTime) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
func (s byTime) Less(i, j int) bool {
	// Two zero times should return false.
	// Otherwise, zero is "greater" than any other time.
	// (To sort it at the end of the list.)
	if s[i].Next.IsZero() {
		return false
	}
	if s[j].Next.IsZero() {
		return true
	}
	// 按时间由小到大排序
	return s[i].Next.Before(s[j].Next)
}

// New 返回一个新的 Cron 任务执行器，可以通过给定的选项进行修改。
//
// 可用设置
//
//	时区
//	  描述: 用于解析调度计划的时区
//	  默认值: time.Local
//	解析器
//	  描述: 解析器将 cron 规范字符串转换为 cron.Schedules。
//	  默认值: 接受此规范: https://en.wikipedia.org/wiki/Cron
//	链
//	  描述: 对提交的任务进行包装以自定义行为。
//	  默认值: 一个会捕获 panic 并将其记录到 stderr 的链。
//
// 使用 "cron.With*" 来修改默认行为。
func New(opts ...Option) *Cron {
	c := &Cron{
		entries:   nil,               // 开始时作业列表为空
		chain:     NewChain(),        // 创建一个链对象
		add:       make(chan *Entry), // 初始化各个 channel 对象
		stop:      make(chan struct{}),
		snapshot:  make(chan chan []Entry),
		remove:    make(chan EntryID),
		running:   false,          // 未运行
		runningMu: sync.Mutex{},   // 初始化互斥锁对象
		logger:    DefaultLogger,  // 使用默认日志对象
		location:  time.Local,     // 本地区域
		parser:    standardParser, // 使用默认的解析器
	}
	for _, opt := range opts { // 应用选项，替换掉默认值
		opt(c)
	}
	return c
}

// FuncJob 是一个将 func() 转换为 cron.Job 的装饰器
type FuncJob func()

// Run 实现 cron.Job 接口
func (f FuncJob) Run() { f() }

// AddFunc 将作业函数（cmd）添加到执行器 Cron 中，以按给定的调度计划（spec）运行
// 返回一个作业 ID，之后可以使用这个 ID 将作业从执行器中移除
func (c *Cron) AddFunc(spec string, cmd func()) (EntryID, error) {
	// 将 cmd 包装成 cron.Job 后转发给 AddJob 方法
	return c.AddJob(spec, FuncJob(cmd))
}

// AddJob 将一个 cron.Job 添加到执行器 Cron 中，以按给定的执行计划（spec）运行
// 返回一个作业 ID，之后可以使用这个 ID 将作业从执行器中移除
func (c *Cron) AddJob(spec string, cmd Job) (EntryID, error) {
	schedule, err := c.parser.Parse(spec) // 解析任务的执行计划（spec）并将其转换成 Schedule 对象
	if err != nil {
		return 0, err
	}
	// 将作业注册到 Cron
	return c.Schedule(schedule, cmd), nil
}

// Schedule 将 Job 添加到 Cron 中，以按给定的执行计划 schedule 运行
// 会使用配置的 Chain 对作业进行装饰
func (c *Cron) Schedule(schedule Schedule, cmd Job) EntryID {
	c.runningMu.Lock() // 加锁保证并发安全
	defer c.runningMu.Unlock()
	c.nextID++       // 计算作业 ID（由此可见作业 ID 是自增的）
	entry := &Entry{ // 构造一个作业实例对象
		ID:         c.nextID,          // 作业 ID
		Schedule:   schedule,          // 执行计划
		WrappedJob: c.chain.Then(cmd), // 装饰作业，附加可选的功能
		Job:        cmd,               // 作业函数
	}
	if !c.running { // 如果 Cron 未运行
		c.entries = append(c.entries, entry) // 直接追加到 entries 列表
	} else { // 已运行（调用过 Start/Run 方法）
		c.add <- entry // 放入 add channel 中，通知 Cron 的调度器有新的作业被加入
	}
	return entry.ID // 返回作业 ID
}

// Entries 返回执行器 Cron 中当前作业列表的快照
func (c *Cron) Entries() []Entry {
	c.runningMu.Lock() // 加锁保证并发安全
	defer c.runningMu.Unlock()
	// 如果调度器正在运行
	if c.running {
		replyChan := make(chan []Entry, 1) // 构造一个 channel 用来传递作业列表
		c.snapshot <- replyChan            // 通知调度器返回当前作业列表 c.entries 的副本
		return <-replyChan
	}
	// 如果调度器未运行，则直接返回当前作业列表 c.entries 的副本
	return c.entrySnapshot()
}

// Location 获取执行器配置的时区
func (c *Cron) Location() *time.Location {
	return c.location
}

// Entry 返回给定 ID 的作业对象快照，如果找不到，则返回空对象
func (c *Cron) Entry(id EntryID) Entry {
	for _, entry := range c.Entries() {
		if id == entry.ID {
			return entry
		}
	}
	return Entry{}
}

// Remove 移除一个给定 ID 的作业
func (c *Cron) Remove(id EntryID) {
	c.runningMu.Lock() // 加锁保证并发安全
	defer c.runningMu.Unlock()
	if c.running { // 如果调度器正在运行
		c.remove <- id // 通知调度器移除指定作业
	} else { // 如果调度器未运行
		c.removeEntry(id) // 可以直接从作业列表 c.entries 中移除指定作业
	}
}

// Start 在新的 goroutine 中启动执行器 Cron 进行作业调度
// 如果之前已经启动，则什么也不做（no-op）
func (c *Cron) Start() {
	c.runningMu.Lock() // 加锁保证并发安全
	defer c.runningMu.Unlock()
	if c.running { // 如果已经在运行中，no-op
		return
	}
	c.running = true // 标记为运行中
	go c.run()       // 开启新的 goroutine 异步起动执行器
}

// Run 启动执行器 Cron 进行作业调度
// 如果之前已经启动，则什么也不做（no-op）
func (c *Cron) Run() {
	c.runningMu.Lock()
	if c.running {
		c.runningMu.Unlock()
		return
	}
	c.running = true
	c.runningMu.Unlock()
	c.run() // 同步启动执行器
}

// 运行执行器 Cron，开始进行作业调度
// 调度器主逻辑放在私有方法 run 中，目的是为了在运行前标记 running 字段状态
func (c *Cron) run() {
	c.logger.Info("start")

	// 计算每个作业的下一次执行时间
	now := c.now()
	for _, entry := range c.entries {
		entry.Next = entry.Schedule.Next(now)
		c.logger.Info("schedule", "now", now, "entry", entry.ID, "next", entry.Next)
	}

	// 外层 for 循环每轮次会对作业列表 c.entries 进行排序，并且重新计算 timer
	for {
		// NOTE: 外层 for 循环逻辑

		// 按时间排序，确定下一个要执行的作业（距离现在最近的）
		sort.Sort(byTime(c.entries))

		// timer 会作为内层 for-select 其中的一个 case
		var timer *time.Timer
		if len(c.entries) == 0 || c.entries[0].Next.IsZero() {
			// 如果还未注册任何作业，则将 timer 设置一个比较长的时间（这不会影响作业的注册和停止操作）
			timer = time.NewTimer(100000 * time.Hour)
		} else {
			// 如果已注册作业，取下一个要执行作业的时间
			timer = time.NewTimer(c.entries[0].Next.Sub(now))
		}

		// NOTE: 内层 for 循环逻辑

		// 内层 for 循环是作业调度主逻辑，会监听所有调度期间触发的事件
		for {
			select {
			case now = <-timer.C: // 本轮次等待结束
				now = now.In(c.location) // 当前被唤醒的时间，这个 now 是由 timer 返回的，需要确保时区正确
				c.logger.Info("wake", "now", now)

				// 运行下一次执行时间小于当前时间的所有作业
				for _, e := range c.entries {
					if e.Next.After(now) || e.Next.IsZero() {
						break // 还未到执行时间 break（c.entries 已根据时间排过序）
					}
					c.startJob(e.WrappedJob)      // 执行被装饰过的作业，内部会启动新的 goroutine 来执行
					e.Prev = e.Next               // 记录这次执行作业的时间到 Prev
					e.Next = e.Schedule.Next(now) // 计算下一次执行作业的时间并记录到 Next
					c.logger.Info("run", "now", now, "entry", e.ID, "next", e.Next)
				}

			case newEntry := <-c.add: // 有新的作业加入进来
				timer.Stop()                                // 停止当前 timer
				now = c.now()                               // 获取当前时间
				newEntry.Next = newEntry.Schedule.Next(now) // 计算新加入作业的下一次执行时间
				c.entries = append(c.entries, newEntry)     // 将新加入的作业追加到 c.entries 列表
				c.logger.Info("added", "now", now, "entry", newEntry.ID, "next", newEntry.Next)

			case replyChan := <-c.snapshot: // 获取当前作业列表
				replyChan <- c.entrySnapshot() // 传递当前作业列表快照给 replyChan
				continue                       // 没有修改 entries，所以不影响 timer，直接下一轮循环继续 select-case

			case <-c.stop: // 停止执行器 Cron
				timer.Stop() // 停止 timer 并退出程序
				c.logger.Info("stop")
				return

			case id := <-c.remove: // 移除指定 ID 的作业
				timer.Stop()      // 停止当前 timer
				now = c.now()     // 更新当前时间
				c.removeEntry(id) // 移除作业
				c.logger.Info("removed", "entry", id)
			}

			// case 执行完成后会走到这里
			break // 打断内层循环，需要靠外层循环重新对作业列表 c.entries 进行排序，和重新计算 timer
		}
	}
}

// startJob 在新的 goroutine 中运行给定的作业
func (c *Cron) startJob(j Job) {
	c.jobWaiter.Add(1) // 运行作业数 + 1
	go func() {
		defer c.jobWaiter.Done() // 作业完成，wg 计数器 - 1
		j.Run()                  // 运行作业
	}()
}

// 返回执行器 Cron 所配置时区的当前时间
func (c *Cron) now() time.Time {
	// 这里将当前时间转换为 c.location 指定的时区
	return time.Now().In(c.location)
}

// Stop 如果执行器 Cron 的调度器正在运行，则停止它；否则什么也不做（does nothing）
// 返回一个上下文，以便调用方可以等待正在运行的作业完成
func (c *Cron) Stop() context.Context {
	c.runningMu.Lock() // 加锁保证并发安全
	defer c.runningMu.Unlock()
	if c.running { // 如果调度器正在运行
		c.stop <- struct{}{} // 发送停止信号，通知调度器停止
		c.running = false    // 标记已停止
	}
	ctx, cancel := context.WithCancel(context.Background())
	go func() { // 开启新的 goroutine 等待所有正在执行的作业完成
		c.jobWaiter.Wait()
		cancel()
	}()
	return ctx // 返回带有 cancel 功能的 context，等待所有作业完成时 cancel() 会被调用，调用方就能拿到完成信号
}

// 返回当前作业列表 c.entries 的副本
func (c *Cron) entrySnapshot() []Entry {
	var entries = make([]Entry, len(c.entries))
	for i, e := range c.entries {
		entries[i] = *e
	}
	return entries
}

// 从作业列表 c.entries 中移除给定 ID 的作业
func (c *Cron) removeEntry(id EntryID) {
	var entries []*Entry
	for _, e := range c.entries {
		if e.ID != id {
			entries = append(entries, e)
		}
	}
	c.entries = entries
}
