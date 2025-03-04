package store

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"github.com/jianghushinian/blog-go-example/nightwatch/pkg/meta"
	"github.com/jianghushinian/blog-go-example/nightwatch/pkg/model"
)

type TaskStore interface {
	Create(ctx context.Context, task *model.Task) error
	Get(ctx context.Context, taskID string) (*model.Task, error)
	List(ctx context.Context, opts ...meta.ListOption) (int64, []*model.Task, error)
	Update(ctx context.Context, task *model.Task) error
	Delete(ctx context.Context, taskID string) error
}

type taskStore struct {
	ds *datastore
}

func newTaskStore(ds *datastore) *taskStore {
	return &taskStore{ds}
}

func (d *taskStore) db(ctx context.Context) *gorm.DB {
	return d.ds.Core(ctx)
}

func (d *taskStore) Create(ctx context.Context, task *model.Task) error {
	return d.db(ctx).Create(&task).Error
}

func (d *taskStore) Get(ctx context.Context, taskID string) (*model.Task, error) {
	task := &model.Task{}
	if err := d.db(ctx).Where("id = ?", taskID).First(&task).Error; err != nil {
		return nil, err
	}

	return task, nil
}

func (d *taskStore) List(ctx context.Context, opts ...meta.ListOption) (count int64, ret []*model.Task, err error) {
	o := meta.NewListOptions(opts...)

	ans := d.db(ctx).
		Where(o.Filters).
		Not(o.Not).
		Offset(o.Offset).
		Limit(defaultLimit(o.Limit)).
		Order(defaultOrder(o.Order)).
		Find(&ret).
		Offset(-1).
		Limit(-1).
		Count(&count)

	return count, ret, ans.Error
}

func (d *taskStore) Update(ctx context.Context, task *model.Task) error {
	return d.db(ctx).Save(task).Error
}

func (d *taskStore) Delete(ctx context.Context, taskID string) error {
	err := d.db(ctx).Where("id = ?", taskID).Delete(&model.Task{}).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	return nil
}
