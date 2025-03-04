package task

import (
	"context"
	"log/slog"
	"sync"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"

	"github.com/jianghushinian/blog-go-example/nightwatch/internal/watcher"
	"github.com/jianghushinian/blog-go-example/nightwatch/pkg/meta"
	"github.com/jianghushinian/blog-go-example/nightwatch/pkg/model"
	"github.com/jianghushinian/blog-go-example/nightwatch/pkg/store"
)

var _ watcher.Watcher = (*taskWatcher)(nil)

type taskWatcher struct {
	store     store.IStore
	clientset kubernetes.Interface

	wg sync.WaitGroup
}

func (w *taskWatcher) Init(ctx context.Context, config *watcher.Config) error {
	w.store = config.Store
	w.clientset = config.Clientset
	return nil
}

func (w *taskWatcher) Spec() string {
	return "@every 30s"
}

// Run 运行 task watcher 任务
func (w *taskWatcher) Run() {
	w.wg.Add(2)

	slog.Debug("Sync period is start")

	// NOTE: 将 Normal 状态任务在 K8s 中启动
	go func() {
		defer w.wg.Done()
		ctx := context.Background()

		_, tasks, err := w.store.Tasks().List(ctx, meta.WithFilter(map[string]any{
			"status": model.TaskStatusNormal,
		}))
		if err != nil {
			slog.Error(err.Error(), "Failed to list tasks")
			return
		}

		var wg sync.WaitGroup
		wg.Add(len(tasks))
		for _, task := range tasks {
			go func(task *model.Task) {
				defer wg.Done()
				job, err := w.clientset.BatchV1().Jobs(task.Namespace).Create(ctx, toJob(task), metav1.CreateOptions{})
				if err != nil {
					slog.Error(err.Error(), "Failed to create job")
					return
				}

				task.Status = model.TaskStatusPending
				if err := w.store.Tasks().Update(ctx, task); err != nil {
					slog.Error(err.Error(), "Failed to update task status")
					return
				}
				slog.Info("Successfully created job", "namespace", job.Namespace, "name", job.Name)
			}(task)
		}
		wg.Wait()
	}()

	// NOTE: 同步中间状态的任务在 K8s 中的状态到表中
	go func() {
		defer w.wg.Done()
		ctx := context.Background()

		_, tasks, err := w.store.Tasks().List(ctx, meta.WithFilterNot(map[string]any{
			// 排除这几个状态
			"status": []model.TaskStatus{model.TaskStatusNormal, model.TaskStatusSucceeded, model.TaskStatusFailed},
		}))
		if err != nil {
			slog.Error(err.Error(), "Failed to list tasks")
			return
		}

		var wg sync.WaitGroup
		wg.Add(len(tasks))
		for _, task := range tasks {
			go func(task *model.Task) {
				defer wg.Done()
				job, err := w.clientset.BatchV1().Jobs(task.Namespace).Get(ctx, task.Name, metav1.GetOptions{})
				if err != nil {
					slog.Error(err.Error(), "Failed to get task")
					return
				}

				task.Status = toTaskStatus(job)
				if err := w.store.Tasks().Update(ctx, task); err != nil {
					slog.Error(err.Error(), "Failed to update task status")
					return
				}
				slog.Info("Successfully sync job status to task", "namespace", job.Namespace, "name", job.Name, "status", task.Status)
			}(task)
		}
		wg.Wait()
	}()

	w.wg.Wait()
	slog.Debug("Sync period is complete")
}

func init() {
	watcher.Register(&taskWatcher{})
}
