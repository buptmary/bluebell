package mysql

import (
	"bluebell/models"
	"errors"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// Login 登录业务
func Login(user *models.User) (err error) {
	originPassword := user.Password // 记录用户的登录密码
	err = db.Where("username = ?", user.UserName).First(&user).Error
	if err != nil {
		// 用户不存在
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrorUserNotExist
		}
		// 查询数据库出错
		return err
	}

	// 比较密码
	if ok := ComparePassword(user.Password, originPassword); !ok {
		return ErrorPasswordWrong
	}
	return nil
}

// CheckUserExist 检查指定用户名用户是否存在
func CheckUserExist(username string) (err error) {
	var count int64
	err = db.Model(models.User{}).Where("username = ?", username).Count(&count).Error
	if count > 0 {
		return ErrorUserExist
	}
	return err
}

// InsertUser 向数据库插入一条新的用户记录
func InsertUser(user *models.User) error {
	// 对密码进行加密
	hashPass, err := hashEncode(user.Password)
	if err != nil {
		return err
	}
	user.Password = hashPass
	// 执行插入SQL语句
	err = db.Create(user).Error
	if err != nil {
		return ErrorInsertFailed
	}
	return nil
}

// GetUserByID 根据ID查询用户信息
func GetUserByID(uid int64) (user *models.User, err error) {
	user = new(models.User)
	err = db.Select("user_id", "username").Where("user_id = ?", uid).First(user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			err = ErrorUserNotExist
			zap.L().Error("mysql GetUserByID failed", zap.Error(err))
			return nil, err
		}
		return nil, ErrorInvalidID
	}
	return user, nil
}

/*------------------------ 以下为工具函数 ------------------------*/

// 对用户密码进行加密 []byte 字节切片
func hashEncode(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

// ComparePassword 验证密码，password1为加密的密码，password2为待验证的密码
func ComparePassword(hashPass, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashPass), []byte(password))
	if err != nil {
		return false
	}
	return true
}
