package logic

import (
	"blue-bell_back/dao/mysql"
	"blue-bell_back/models"
	"blue-bell_back/pkg/jwt"
	"blue-bell_back/pkg/snowflake"

	"go.uber.org/zap"
)

// SignUp 用户注册函数
// 参数 p 包含用户输入的用户名和密码
// 返回值 error 用于返回注册过程中可能发生的错误
func SignUp(p *models.ParamSignUp) (err error) {
	//处理注册逻辑
	//1.判断用户是否存在 数据库中检查
	if err := mysql.CheckUserExist(p.UserName); err != nil {
		return err
	}

	//2.生成UID
	userID := snowflake.GenID()
	if err != nil {
		zap.L().Error("user snowflake.GenId failed.", zap.Error(err))
		return
	}
	//构造一个User实例
	user := &models.User{
		UserID:   userID,
		UserName: p.UserName,
		Password: p.RePassword,
	}
	//3.密码加密并保存进数据库
	return mysql.InsertUser(user)

}

// Login 用户登录函数
// 参数 p 包含用户输入的用户名和密码
// 返回值 user 是登录成功的用户信息，包括用户ID、用户名和令牌(Token)
func Login(p *models.ParamLogin) (user *models.User, err error) {
	//初始化用户信息
	user = &models.User{
		UserName: p.UserName,
		Password: p.Password,
	}

	//return mysql.Login(user)
	// 调用mysql.Login函数执行登录操作，如果登录失败，返回错误信息

	if err := mysql.Login(user); err != nil {
		return nil, err
	}

	//生成用户登录令牌，生成失败，返回错误信息
	user.Token, err = jwt.GenToken(user.UserID, user.UserName)
	if err != nil {
		return
	}

	// //生成JWT
	// return jwt.GenToken(user.UserID, p.UserName)
	return
}
