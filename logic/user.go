package logic

import (
	"bluebell/dao/mysql"
	"bluebell/models"
	"bluebell/pkg/snowflake"
)

// SignUp 注册业务逻辑
func SignUp(fo *models.RegisterForm) (err error) {
	// 1. 判断用户是否存在
	err = mysql.CheckUserExist(fo.UserName)
	if err != nil {
		return err
	}

	// 2. 为注册用户生成UserID 雪花ID
	userID := snowflake.GenID()

	// 3. 构造User实例
	u := &models.User{
		UserID:   userID,
		UserName: fo.UserName,
		Password: fo.Password,
	}

	// 4. 存入数据库
	err = mysql.InsertUser(u)
	return err
}

// Login 登录业务逻辑
func Login(fo *models.LoginForm) (user *models.User, err error) {
	// 1. 新建一个user实例
	user = &models.User{
		UserName: fo.UserName,
		Password: fo.Password,
	}

	// 2. 用户登录验证
	if err := mysql.Login(user); err != nil {
		return nil, err
	}
	return user, nil
}
