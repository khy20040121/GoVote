package controller

import (
	"bluebell/dao/redis"
	"bluebell/logic"
	"bluebell/models"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// 帖子投票的控制函数
func PostVoteHandler(c *gin.Context) {
	// 1 参数绑定
	p := new(models.ParamVoteData)
	if err := c.ShouldBind(p); err != nil {
		zap.L().Error("PostVoteHandler ShouldBind error", zap.Error(err))
		ResponseError(c, CodeInvalidParam)
		return
	}

	// 2 业务处理

	// 获取当前投票用户的id
	userID, err := getCurrentUser(c)
	if err != nil {
		zap.L().Error("without login", zap.Error(err))
		ResponseError(c, CodeNeedLogin)
		return
	}

	// 具体投票的业务逻辑
	if err := logic.VoteForPost(userID, p); err != nil {
		zap.L().Error("logic.VoteForPost error", zap.Error(err))
		if err == redis.ErrVoteRepeated {
			ResponseError(c, CodeVoteRepeated)
			return
		}
		ResponseError(c, CodeServerBusy)
		return
	}

	// 3 返回响应
	ResponseSuccess(c, nil)
}
