package workerpool

import (
	"context"
	"sync"
	"sync/atomic"
	"time"

	"github.com/gammazero/deque"
)

const (
	// 工作协程（worker）处于空闲状态的超时时间，超过此时间就会关闭 worker
	idleTimeout = 2 * time.Second
)

// New 创建并启动协程池
// maxWorkers 参数指定可以并发执行任务的最大工作协程数。
// 当没有任务需要执行时，工作协程（worker）会逐渐停止，直至没有剩余的 worker。
func New(maxWorkers int) *WorkerPool {
	// 至少有一个 worker
	if maxWorkers < 1 {
		maxWorkers = 1
	}

	// 实例化协程池对象
	pool := &WorkerPool{
		maxWorkers:  maxWorkers,
		taskQueue:   make(chan func()),
		workerQueue: make(chan func()),
		stopSignal:  make(chan struct{}),
		stoppedChan: make(chan struct{}),
	}

	// 启动任务调度器
	go pool.dispatch()

	return pool
}

// WorkerPool 是 Go 协程的集合池，用于确保同时处理请求的协程数量严格受控于预设的上限值
type WorkerPool struct {
	maxWorkers   int                 // 最大工作协程数
	taskQueue    chan func()         // 任务提交队列
	workerQueue  chan func()         // 工作协程消费队列
	stoppedChan  chan struct{}       // 停止完成通知通道
	stopSignal   chan struct{}       // 停止信号通道
	waitingQueue deque.Deque[func()] // 等待队列（双端队列）
	stopLock     sync.Mutex          // 停止操作互斥锁
	stopOnce     sync.Once           // 控制只停止一次
	stopped      bool                // 是否已经停止
	waiting      int32               // 等待队列中任务计数
	wait         bool                // 协程池退出时是否等待已入队任务执行完成
}

// Size 返回协程池大小
func (p *WorkerPool) Size() int {
	return p.maxWorkers
}

// Stop 停止工作池并仅等待当前正在执行的任务完成。
// 所有尚未开始执行的等待任务将被丢弃。调用 Stop 后
// 禁止继续向工作池提交新任务。
//
// 由于创建工作池时会自动启动至少一个调度协程（dispatcher），
// 当工作池不再需要时，必须调用 Stop() 或 StopWait() 方法
// 以确保正确释放资源。
func (p *WorkerPool) Stop() {
	p.stop(false)
}

// StopWait 停止工作池并等待所有已入队任务执行完毕。
// 调用后禁止提交新任务，但会确保所有队列中的任务在函数返回前由工作协程处理完成。
func (p *WorkerPool) StopWait() {
	p.stop(true)
}

// Stopped 如果工作池已停止，则返回 true
func (p *WorkerPool) Stopped() bool {
	p.stopLock.Lock()
	defer p.stopLock.Unlock()
	return p.stopped
}

// Submit 将任务函数提交到工作池队列等待执行，不会等待任务执行完成
//
// 任务函数所需的外部变量必须通过闭包捕获。
// 需要返回值的任务应当通过闭包内的通道（channel）传递结果。
//
// 无论提交多少个任务，Submit 都不会阻塞调用方。
// 任务会立即分配给可用的 worker 或启动新的 worker 执行任务。
// 如果达到最大 worker 数量限制，没有可用的 worker，则任务将放入等待队列（waitingQueue）中。
//
// 当等待队列（waitingQueue）中存在任务时，所有新提交的任务都会进入队列，
// 等待工作协程（worker）就绪后按先进先出（FIFO）顺序处理。
//
// 只要没有接收到新的任务，系统将每隔一定时间间隔（默认 2 秒）终止一个
// 空闲工作协程，直至所有空闲协程都被回收。基于 Go 协程的轻量级特性，
// 新建协程的耗时开销可忽略不计，因此无需长期维持空闲协程池。
func (p *WorkerPool) Submit(task func()) {
	if task != nil {
		p.taskQueue <- task
	}
}

// SubmitWait 提交任务函数到队列，并阻塞等待任务执行完成
func (p *WorkerPool) SubmitWait(task func()) {
	if task == nil {
		return
	}
	doneChan := make(chan struct{})
	p.taskQueue <- func() { // 提交任务
		task()
		close(doneChan)
	}
	<-doneChan // 阻塞等待任务执行完成
}

// WaitingQueueSize 返回等待队列中的任务计数
func (p *WorkerPool) WaitingQueueSize() int {
	return int(atomic.LoadInt32(&p.waiting))
}

// Pause 通过 Context 控制协程池的暂停与恢复
// 当所有工作协程都进入等待状态时，Pause 才会返回。
// 任务可以继续被提交到协程池，但在 Context 被取消或超时之前不会执行这些任务。
//
// 当协程池已经处于暂停状态时再次调用 Pause，将导致该调用阻塞直到之前的所有暂停被取消。
// 这允许一个 goroutine 在其他 goroutine 解除暂停后立即接管工协程池池的暂停控制权。
//
// 当协程池被停止时，所有工作协程将被唤醒，已排队任务会在 StopWait 期间执行。
func (p *WorkerPool) Pause(ctx context.Context) {
	p.stopLock.Lock() // 加锁，确保并发安全
	defer p.stopLock.Unlock()
	if p.stopped { // 已经停止，无需处理
		return
	}
	ready := new(sync.WaitGroup)
	ready.Add(p.maxWorkers) // 设置与最大 worker 数匹配的计数器
	for i := 0; i < p.maxWorkers; i++ {
		p.Submit(func() { // 向每个 worker 发送暂停指令
			ready.Done() // 标记暂停指令发送完成
			select {
			case <-ctx.Done(): // 调用方通过 ctx 取消暂停
			case <-p.stopSignal: // 协程池内部通过接收停止信号取消暂停
			}
		})
	}
	// 阻塞等待所有 worker 都进入暂停状态
	ready.Wait()
}

// 任务派发，循环的将下一个排队中的任务发送给可用的工作协程（worker）执行
func (p *WorkerPool) dispatch() {
	defer close(p.stoppedChan)            // 保证调度器退出时关闭停止通知通道
	timeout := time.NewTimer(idleTimeout) // 创建 2 秒周期的空闲检测定时器
	var workerCount int                   // 当前活跃 worker 计数器
	var idle bool                         // 空闲状态标识
	var wg sync.WaitGroup                 // 用于等待所有 worker 完成

Loop:
	for { // 主循环处理任务分发
		// 当等待队列中存在任务时，程序将进入队列优先模式：
		//   1. 新提交的任务自动进入等待队列尾部
		//   2. 工作协程（worker）从队列头部提取任务执行
		// 等待队列完全清空后，程序自动切换回直通模式：
		//   1. 新任务将直接派发给空闲工作协程（worker）处理
		//   2. 如果工作协程数已达上限，将任务提交到等待队列

		// 队列优先模式：优先检测等待队列
		if p.waitingQueue.Len() != 0 {
			if !p.processWaitingQueue() {
				break Loop // 协程池已经停止
			}
			continue // 队列不为空则继续下一轮循环
		}

		// 直通模式：开始处理提交上来的新任务
		select {
		case task, ok := <-p.taskQueue: // 接收到新任务
			if !ok { // 协程池停止时会关闭任务通道，如果 !ok 说明协程池已停止，退出循环
				break Loop
			}

			select {
			case p.workerQueue <- task: // 尝试派发任务
			default: // 没有空闲的 worker，无法立即派发任务
				if workerCount < p.maxWorkers { // 如果协程池中的活跃协程数量小于最大值，那么创建一个新的协程（worker）来执行任务
					wg.Add(1)
					go worker(task, p.workerQueue, &wg) // 创建新的 worker 执行任务
					workerCount++                       // worker 记数加 1
				} else { // 已达协程池容量上限
					p.waitingQueue.PushBack(task)                              // 将任务提交到等待队列
					atomic.StoreInt32(&p.waiting, int32(p.waitingQueue.Len())) // 原子更新等待计数
				}
			}
			idle = false // 标记为非空闲
		case <-timeout.C: // 空闲超时处理
			// 在一个空闲超时周期内，存在空闲的 workers，那么停止一个 worker
			if idle && workerCount > 0 {
				if p.killIdleWorker() { // 回收一个 worker
					workerCount-- // worker 计数减 1
				}
			}
			idle = true                // 标记为空闲
			timeout.Reset(idleTimeout) // 复用定时器
		}
	}

	if p.wait { // 调用了 StopWait() 方法，需要运行等待队列中的任务，直至队列清空
		p.runQueuedTasks()
	}

	// 终止所有 worker
	for workerCount > 0 {
		p.workerQueue <- nil // 发送终止信号给 worker
		workerCount--        // worker 计数减 1，直至为 0 退出循环
	}
	wg.Wait() // 阻塞等待所有 worker 完成

	timeout.Stop() // 停止定时器
}

// 工作协程，执行任务并在收到 nil 信号时停止
func worker(task func(), workerQueue chan func(), wg *sync.WaitGroup) {
	for task != nil { // 循环执行任务，直至接收到终止信号 nil
		task()               // 执行任务
		task = <-workerQueue // 接收新任务
	}
	wg.Done() // 标记 worker 完成
}

// stop 通知调度协程（dispatcher）退出，wait 参数决定是否等待已入队任务完成
func (p *WorkerPool) stop(wait bool) {
	// 通过 sync.Once 确保停止逻辑仅执行一次
	p.stopOnce.Do(func() {
		// 关闭停止信号通道，用于唤醒所有暂停中的工作协程（worker）
		close(p.stopSignal)
		// 加锁，保证并发安全
		p.stopLock.Lock()
		p.stopped = true // 标记停止
		p.stopLock.Unlock()
		p.wait = wait // 标记是否等待已入队任务执行完成
		// 关闭任务队列通道，停止接收新任务
		close(p.taskQueue)
	})
	<-p.stoppedChan // 阻塞等待调度协程退出
}

// 处理等待队列
// 将接收到的新任务放入等待队列中，并在 worker 可用时从等待队列中删除任务。
// 如果 worker pool 已停止，则返回 false。
func (p *WorkerPool) processWaitingQueue() bool {
	select {
	case task, ok := <-p.taskQueue: // 接收到新任务
		if !ok { // 协程池停止时会关闭任务通道，如果 !ok 说明协程池已停止，返回 false，不再继续处理
			return false
		}
		p.waitingQueue.PushBack(task) // 将新任务加入等待队列队尾
	case p.workerQueue <- p.waitingQueue.Front(): // 从等待队列队头获取任务并放入工作队列
		p.waitingQueue.PopFront() // 任务已经开始处理，所以要从等待队列中移除任务
	}
	atomic.StoreInt32(&p.waiting, int32(p.waitingQueue.Len())) // 原子修改等待队列中任务计数
	return true
}

// 停止一个空闲 worker
func (p *WorkerPool) killIdleWorker() bool {
	select {
	case p.workerQueue <- nil: // 发送终止信号给工作协程（worker）
		// Sent kill signal to worker.
		return true
	default:
		// No ready workers. All, if any, workers are busy.
		return false
	}
}

// 运行等待队列中的任务，直至队列清空
func (p *WorkerPool) runQueuedTasks() {
	for p.waitingQueue.Len() != 0 { // 直至队列清空终止循环
		// 从等待队列中获取队首任务并交给工作队列去执行
		p.workerQueue <- p.waitingQueue.PopFront()
		atomic.StoreInt32(&p.waiting, int32(p.waitingQueue.Len())) // 原子修改等待任务计数
	}
}
