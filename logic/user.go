package logic

import (
	"bluebell/dao/mysql"
	"bluebell/models"
	"bluebell/pkg/jwt"
	"bluebell/pkg/snowflake"
)

// SignUp 用户注册信息的logic
func SignUp(p *models.ParamSignUp) error {
	// 1 判断用户存在不存在
	if err := mysql.CheckUserExist(p.Username); err != nil {
		return err
	}

	// 2生成UID
	userID := snowflake.GenID()

	// 构造一个user实例
	user := &models.User{
		UserID:   userID,
		Username: p.Username,
		Password: p.Password,
	}

	// 3 将user保存在数据库当中去
	return mysql.InsertUser(user)
}

// Login 用户登录的logic
func Login(p *models.ParamLogin) (string, int64, string, error) {
	user := &models.User{
		Username: p.Username,
		Password: p.Password,
	}

	// 进行数据库层面的处理
	if err := mysql.Login(user); err != nil {
		return "", 0, "", err
	}
	token, err := jwt.GenToken(user.UserID, user.Username)
	if err != nil {
		return "", 0, "", err
	}
	return token, user.UserID, user.Username, nil
}
