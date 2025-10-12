package mysql

import (
	"blue-bell_back/models"
	"crypto/md5"
	"database/sql"
	"encoding/hex"
	"errors"
)

// CheckUserExist 检查用户名是否已存在
// 参数:
//
//	username - 待检查的用户名
//
// 返回值:
//
//	如果用户名已存在或查询过程中出错，返回相应的错误

func CheckUserExist(username string) (err error) {
	sqlStr := `select count(user_id) from user where username=?`
	var count int
	if err := db.Get(&count, sqlStr, username); err != nil {
		return err
	}
	if count > 0 {
		return errors.New("用户已存在")
	}
	return
}

// sercret是用于密码加密的密钥
const secret = "liwenzhou.com"

// InsertUser 向数据库中插入一条新的用户记录
// 参数:
//
//	user - 包含用户信息的结构体指针
//
// 返回值:
//
//	如果插入过程中出错，返回相应的错误
func InsertUser(user *models.User) (err error) {
	user.Password = encryptPassword(user.Password)
	sqlStr := `insert into user(user_id, username, password) values(?,?,?)`
	_, err = db.Exec(sqlStr, user.UserID, user.UserName, user.Password)
	return
}

// encryptPassword 加密用户密码
// 参数:
//
//	oPassword - 原始密码
//
// 返回值:
//
//	加密后的密码

func encryptPassword(oPassword string) string {
	h := md5.New()
	h.Write([]byte(secret))
	return hex.EncodeToString(h.Sum([]byte(oPassword)))
}

// Login 验证用户登录
// 参数:
//
//	user - 包含用户登录信息的结构体指针
//
// 返回值:
//
//	如果用户不存在或密码错误，返回相应的错误

func Login(user *models.User) (err error) {
	oPassword := user.Password
	sqlStr := `select user_id, username, password from user where username=?`
	err = db.Get(user, sqlStr, user.UserName)
	if err == sql.ErrNoRows {
		return errors.New("用户不存在")
	}

	if err != nil {
		//查询数据库失败
		return err
	}

	//判断密码是否正确
	password := encryptPassword(oPassword)
	if password != user.Password {
		return errors.New("密码错误")
	}
	return
}
