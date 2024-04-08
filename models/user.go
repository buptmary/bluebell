package models

// User 定义请求参数结构体
type User struct {
	UserID   int64  `json:"user_id,string"`                  // 序列化为string类型，防止数字失真
	UserName string `json:"username" gorm:"column:username"` // 为字段指定列标签
	Password string `json:"password"`
}

// RegisterForm 注册请求结构体
type RegisterForm struct {
	UserName        string `json:"username" binding:"required"`
	Password        string `json:"password" binding:"required"`
	ConfirmPassword string `json:"confirm_password" binding:"required,eqfield=Password"`
}

type LoginForm struct {
	UserName string `json:"username" bind:"required"`
	Password string `json:"password" binding:"required"`
}
