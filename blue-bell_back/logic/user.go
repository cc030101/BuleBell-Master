package logic

import (
	"blue-bell_back/dao/mysql"
	"blue-bell_back/models"
	"blue-bell_back/pkg/snowflake"
)

// SignUp 用户注册函数
// 参数 p 包含用户输入的用户名和密码
// 返回值 error 用于返回注册过程中可能发生的错误
func SignUp(p *models.ParamSignUp) (err error) {
	//处理注册逻辑
	//1.判断用户是否存在 数据库中检查
	if err := mysql.CheckUserExist(p.Username); err != nil {
		return err
	}

	//2.生成UID
	userID := snowflake.GenID()
	//构造一个User实例
	user := &models.User{
		UserID:   userID,
		Username: p.Username,
		Password: p.RePassword,
	}

	//3.密码加密

	//4.保存进入数据库
	return mysql.InsertUser(user)

}

// Login 用户登录函数
// 参数 p 包含用户输入的用户名和密码
// 返回值 error 用于返回登录过程中可能发生的错误

func Login(p *models.ParamLogin) error {
	user := &models.User{
		Username: p.Username,
		Password: p.Password,
	}

	return mysql.Login(user)
}
