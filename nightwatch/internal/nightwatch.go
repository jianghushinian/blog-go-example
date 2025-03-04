package nightwatch

import (
	"context"
	"log/slog"
	"time"

	"github.com/go-redsync/redsync/v4"
	"github.com/go-redsync/redsync/v4/redis/goredis/v9"
	"github.com/robfig/cron/v3"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/kubernetes"

	"github.com/jianghushinian/blog-go-example/nightwatch/internal/watcher"
	"github.com/jianghushinian/blog-go-example/nightwatch/pkg/db"
	"github.com/jianghushinian/blog-go-example/nightwatch/pkg/store"
	// 触发 init 函数
	_ "github.com/jianghushinian/blog-go-example/nightwatch/internal/watcher/all"
)

var (
	lockName          = "nightwatch-lock"
	jobStopTimeout    = 3 * time.Minute
	extendExpiration  = 5 * time.Second
	defaultExpiration = 2 * extendExpiration
)

// Watcher 管理器
type nightWatch struct {
	runner *cron.Cron     // 执行器
	locker *redsync.Mutex // 分布式锁
	config *watcher.Config
}

// Config 配置信息，用于创建 nightWatch 对象
type Config struct {
	MySQLOptions *db.MySQLOptions
	RedisOptions *db.RedisOptions
	Clientset    kubernetes.Interface
}

// CreateWatcherConfig 创建 nightWatch 需要的配置
func (c *Config) CreateWatcherConfig() (*watcher.Config, error) {
	gormDB, err := db.NewMySQL(c.MySQLOptions)
	if err != nil {
		slog.Error(err.Error(), "Failed to create MySQL client")
		return nil, err
	}
	datastore := store.NewStore(gormDB)

	return &watcher.Config{Store: datastore, Clientset: c.Clientset}, nil
}

// New 通过配置构造一个 nightWatch 对象
func (c *Config) New() (*nightWatch, error) {
	rdb, err := db.NewRedis(c.RedisOptions)
	if err != nil {
		slog.Error(err.Error(), "Failed to create Redis client")
		return nil, err
	}

	logger := newCronLogger()
	runner := cron.New(
		cron.WithSeconds(),
		cron.WithLogger(logger),
		cron.WithChain(cron.SkipIfStillRunning(logger), cron.Recover(logger)),
	)

	pool := goredis.NewPool(rdb)
	lockOpts := []redsync.Option{
		redsync.WithRetryDelay(50 * time.Microsecond),
		redsync.WithTries(3),
		redsync.WithExpiry(defaultExpiration),
	}
	locker := redsync.New(pool).NewMutex(lockName, lockOpts...)

	cfg, err := c.CreateWatcherConfig()
	if err != nil {
		return nil, err
	}

	nw := &nightWatch{runner: runner, locker: locker, config: cfg}
	if err := nw.addWatchers(); err != nil {
		return nil, err
	}

	return nw, nil
}

// 注册所有 Watcher 实例到 nightWatch
func (nw *nightWatch) addWatchers() error {
	for n, w := range watcher.ListWatchers() {
		if err := w.Init(context.Background(), nw.config); err != nil {
			slog.Error(err.Error(), "Failed to construct watcher", "watcher", n)
			return err
		}

		spec := watcher.Every3Seconds
		if obj, ok := w.(watcher.ISpec); ok {
			spec = obj.Spec()
		}

		if _, err := nw.runner.AddJob(spec, w); err != nil {
			slog.Error(err.Error(), "Failed to add job to the cron", "watcher", n)
			return err
		}
	}

	return nil
}

// Run 执行异步任务，此方法会阻塞直到关闭 stopCh
func (nw *nightWatch) Run(stopCh <-chan struct{}) {
	ctx := wait.ContextForChannel(stopCh)

	// 循环加锁，直到加锁成功，再去启动任务
	ticker := time.NewTicker(defaultExpiration + (5 * time.Second))
	defer ticker.Stop()
	for {
		err := nw.locker.LockContext(ctx)
		if err == nil {
			slog.Info("Successfully acquired lock", "lockName", lockName)
			break
		}
		slog.Debug("Failed to acquire lock", "lockName", lockName, "err", err)
		<-ticker.C
	}

	// 看门狗，实现锁自动续约
	ticker = time.NewTicker(extendExpiration)
	defer ticker.Stop()
	go func() {
		for {
			<-ticker.C
			if ok, err := nw.locker.ExtendContext(ctx); !ok || err != nil {
				slog.Debug("Failed to extend lock", "err", err, "status", ok)
			}
		}
	}()

	// 启动定时任务
	nw.runner.Start()
	slog.Info("Successfully started nightwatch server")

	// 阻塞等待退出信号
	<-stopCh

	nw.stop()
}

// 停止异步任务
func (nw *nightWatch) stop() {
	ctx := nw.runner.Stop()
	select {
	case <-ctx.Done():
	case <-time.After(jobStopTimeout):
		slog.Error("Context was not done immediately", "timeout", jobStopTimeout.String())
	}

	if ok, err := nw.locker.Unlock(); !ok || err != nil {
		slog.Debug("Failed to unlock", "err", err, "status", ok)
	}
}
