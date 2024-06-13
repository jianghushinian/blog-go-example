package api

// CreateUserRequest `POST /users` 创建用户接口请求参数
type CreateUserRequest struct {
	Email    string `json:"email"`
	Nickname string `json:"nickname"`
	Username string `json:"username"`
	Password string `json:"password"`
}
