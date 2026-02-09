package mysql

import (
	"bluebell/models"
	"crypto/md5"
	"database/sql"
	"encoding/hex"
	"errors"

	"go.uber.org/zap"
)

const secret = "khy"

func CheckUserExist(username string) error {
	sqlStr := `select count(user_id) from user where username = ?`

	var count int64

	if err := db.Get(&count, sqlStr, username); err != nil {
		// 数据库查询错误, 返回
		return err
	}
	if count > 0 {
		// 返回错误, 用户已经存在,无法注册
		return ErrorUserExist
	}
	return nil
}

// InsertUser把注册的用户信息插入到数据库当中去
func InsertUser(user *models.User) error {
	// 1 首先对用户密码加密
	user.Password = encryptPassword(user.Password)

	// 2 执行sql语句将user插入到数据库
	sqlStr := `insert into user (user_id, username, password) values (?, ?, ?)`
	_, err := db.Exec(sqlStr, user.UserID, user.Username, user.Password)
	return err

}

// Login检测用户输入的用户名和密码是否正确
func Login(user *models.User) error {
	oPassword := user.Password // 记录一下原始密码,与后面的数据库密码进行比较
	Password := encryptPassword(oPassword)

	sqlStr := `select user_id, username, password from user where username = ?`
	if err := db.Get(user, sqlStr, user.Username); err != nil {
		zap.L().Error("mysql.Query fail", zap.Error(err))
		if errors.Is(err, sql.ErrNoRows) {
			return ErrorUserNotExist
		}
		return err
	}

	// 判断密码是否正确
	if Password != user.Password {
		return ErrorInvalidPassword
	}
	return nil

}

// GetUserByID 根据userID查询user
func GetUserByID(id int64) (user *models.User, err error) {
	sqlStr := `select user_id, username, password from user where user_id = ?`
	user = new(models.User)
	err = db.Get(user, sqlStr, id)
	return
}

// encryptPassword 密码加密
func encryptPassword(oPassword string) string {
	h := md5.New()
	h.Write([]byte(secret))
	return hex.EncodeToString(h.Sum([]byte(oPassword)))
}
