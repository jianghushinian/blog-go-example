package store

import (
	"context"

	"github.com/google/wire"
	"gorm.io/gorm"

	"github.com/jianghushinian/blog-go-example/wire/user/internal/model"
)

// ProviderSet 一个 Wire provider sets，用来初始化 store 实例对象，并将 UserStore 接口绑定到 *userStore 类型实现上
var ProviderSet = wire.NewSet(New, wire.Bind(new(UserStore), new(*userStore)))

// UserStore 定义 user 暴露的 CRUD 方法
type UserStore interface {
	Create(ctx context.Context, user *model.UserM) error
}

// UserStore 接口实现
type userStore struct {
	db *gorm.DB
}

// 确保 userStore 实现了 UserStore 接口
var _ UserStore = (*userStore)(nil)

// New userStore 构造函数
func New(db *gorm.DB) *userStore {
	return &userStore{db}
}

// Create 插入一条 user 记录
func (u *userStore) Create(ctx context.Context, user *model.UserM) error {
	return u.db.Create(&user).Error
}
