package controller

import "github.com/gin-gonic/gin"

type ResponseData struct {
	Code ResCode     `json:"code"`
	Msg  interface{} `json:"msg"`
	Data interface{} `json:"data"`
}

func ResponseSuccess(c *gin.Context, data interface{}) {
	c.JSON(200, &ResponseData{
		Code: CodeSuccess,
		Msg:  CodeSuccess.Msg(),
		Data: data,
	})
}

func ResponseError(c *gin.Context, code ResCode) {
	c.JSON(200, &ResponseData{
		Code: code,
		Msg:  code.Msg(),
		Data: nil,
	})
}

func ResponseErrorWithMsg(c *gin.Context, code ResCode, msg string) {
	c.JSON(200, &ResponseData{
		Code: code,
		Msg:  msg,
		Data: nil,
	})
}
