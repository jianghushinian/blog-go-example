package biz

import (
	"context"

	"github.com/google/wire"
	"github.com/jinzhu/copier"

	"github.com/jianghushinian/blog-go-example/wire/user/internal/model"
	"github.com/jianghushinian/blog-go-example/wire/user/internal/store"
	"github.com/jianghushinian/blog-go-example/wire/user/pkg/api"
)

// ProviderSet 一个 Wire provider sets，用来初始化 biz 实例对象，并将 UserBiz 接口绑定到 *userBiz 类型实现上
var ProviderSet = wire.NewSet(New, wire.Bind(new(UserBiz), new(*userBiz)))

// UserBiz 定义 user 业务逻辑操作方法
type UserBiz interface {
	Create(ctx context.Context, r *api.CreateUserRequest) error
}

// UserBiz 接口的实现
type userBiz struct {
	s store.UserStore
}

// 确保 userBiz 实现了 UserBiz 接口
var _ UserBiz = (*userBiz)(nil)

// New userBiz 构造函数
func New(s store.UserStore) *userBiz {
	return &userBiz{s: s}
}

// Create 创建用户
func (b *userBiz) Create(ctx context.Context, r *api.CreateUserRequest) error {
	var userM model.UserM
	_ = copier.Copy(&userM, r)

	return b.s.Create(ctx, &userM)
}
