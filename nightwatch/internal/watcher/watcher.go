package watcher

import (
	"context"
	"sync"

	"github.com/robfig/cron/v3"

	reflectutil "github.com/jianghushinian/blog-go-example/nightwatch/pkg/util/reflect"
)

const (
	Every3Seconds = "@every 3s"
)

type Watcher interface {
	Init(ctx context.Context, config *Config) error
	cron.Job
}

type ISpec interface {
	// Spec 返回任务的定时周期，支持两种格式
	// 标准 Cron 格式: https://en.wikipedia.org/wiki/Cron
	// Quartz 格式: http://www.quartz-scheduler.org/documentation/quartz-2.3.0/tutorials/tutorial-lesson-06.html
	Spec() string
}

var (
	registryLock = new(sync.Mutex)
	registry     = make(map[string]Watcher)
)

func Register(watcher Watcher) {
	registryLock.Lock()
	defer registryLock.Unlock()

	name := reflectutil.StructName(watcher)
	if _, ok := registry[name]; ok {
		panic("duplicate watcher entry: " + name)
	}

	registry[name] = watcher
}

func ListWatchers() map[string]Watcher {
	registryLock.Lock()
	defer registryLock.Unlock()

	return registry
}
