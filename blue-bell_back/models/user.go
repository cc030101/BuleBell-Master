package models

// 存放数据相关的结构体
type User struct {
	UserID   int64  `db:"user_id"`
	UserName string `db:"username"`
	Password string `db:"password"`
	Token    string `json:"token"`
}
