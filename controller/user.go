package controller

import (
	"bluebell/dao/mysql"
	"bluebell/logic"
	"bluebell/models"
	"errors"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// SignUpHandler处理注册路由
func SignUpHandler(c *gin.Context) {
	// 1 参数校验
	p := new(models.ParamSignUp)
	if err := c.ShouldBind(&p); err != nil {
		//请求参数有误
		zap.L().Error("Signup with invalid param", zap.Error(err))
		ResponseError(c, CodeInvalidParam)
		return
	}

	// 2业务处理
	if err := logic.SignUp(p); err != nil {
		zap.L().Error("logic.Signup failed", zap.Error(err))
		if errors.Is(err, mysql.ErrorUserExist) {
			// 用户已经存在,提示前端
			ResponseError(c, CodeUserExist)
			return
		}
		ResponseError(c, CodeServerBusy)
		return
	}

	// 3返回响应
	ResponseSuccess(c, nil)
}

// LoginHandler登录请求路由
func LoginHandler(c *gin.Context) {
	// 1 参数校验
	p := new(models.ParamLogin)
	if err := c.ShouldBind(&p); err != nil {
		//请求参数有误
		zap.L().Error("Signup with invalid param", zap.Error(err))
		ResponseError(c, CodeInvalidParam)
		return
	}

	// 2 业务逻辑处理
	token, userID, username, err := logic.Login(p)

	if err != nil {
		zap.L().Error("Login error", zap.Error(err))
		if errors.Is(err, mysql.ErrorUserNotExist) {
			ResponseError(c, CodeUserNotExist)
			return
		} else if errors.Is(err, mysql.ErrorInvalidPassword) {
			ResponseError(c, CodeInvalidPassword)
			return
		} else {
			ResponseError(c, CodeServerBusy)
			return
		}
	}

	// 3 返回响应
	ResponseSuccess(c, gin.H{
		"user_id":  strconv.FormatInt(userID, 10),
		"username": username,
		"token":    token,
	})
}
