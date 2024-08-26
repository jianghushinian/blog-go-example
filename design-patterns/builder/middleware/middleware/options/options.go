package options

// Option 用于配置 RBACMiddleware 的函数类型
type Option func(*RBACMiddleware)

// NewRBACMiddleware 创建并配置一个新的 RBACMiddleware 实例
func NewRBACMiddleware(opts ...Option) *RBACMiddleware {
	r := &RBACMiddleware{}
	for _, opt := range opts {
		opt(r)
	}
	return r
}

// WithRole 添加一个允许访问的角色
func WithRole(role string) Option {
	return func(r *RBACMiddleware) {
		r.allowedRoles = append(r.allowedRoles, role)
	}
}
