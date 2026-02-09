package controller

import (
	"bluebell/logic"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// CommunityHandler查询所有的社区的列表
func CommunityHandler(c *gin.Context) {
	// 查询到所有的社区,以community_id, community_name的形式返回
	data, err := logic.GetCommunityList(c)
	if err != nil {
		ResponseError(c, CodeServerBusy)
		return
	}

	ResponseSuccess(c, data)
}

// CommunityDetailHandler根据社区的id查询社区的详情
func CommunityDetailHandler(c *gin.Context) {
	// 1 获取社区的id
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		zap.L().Error("wrong CommunityID param ", zap.Error(err))
		ResponseError(c, CodeInvalidParam)
		return
	}

	// 2 根据社区id查询社区详情
	data, err := logic.GetCommunityDetail(c, int64(id))
	if err != nil {
		zap.L().Error("logic.GetCommunityDetail failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}

	ResponseSuccess(c, data)
}
