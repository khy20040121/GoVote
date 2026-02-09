package controller

import (
	"bluebell/logic"
	"bluebell/models"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// CreatePostHandler 创建帖子的处理函数
func CreatePostHandler(c *gin.Context) {
	// 1 获取参数以及参数校验
	p := new(models.Post)
	if err := c.ShouldBind(p); err != nil {
		zap.L().Error("create post with invalid param", zap.Error(err))
		ResponseError(c, CodeInvalidParam)
		return
	}

	// 从c中获取用户id
	userID, err := getCurrentUser(c)
	if err != nil {
		zap.L().Error("user need to login again", zap.Error(err))
		ResponseError(c, CodeNeedLogin)
		return
	}
	p.AuthorID = userID

	// 2 logic处理
	if err = logic.CreatePost(p); err != nil {
		zap.L().Error("logic.createpost failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}

	// 3 返回响应
	ResponseSuccess(c, nil)
}

// GetPostDetailHandler 获取帖子详情的处理函数
func GetPostDetailHandler(c *gin.Context) {
	// 1 参数以及校验
	pidStr := c.Param("id")
	pid, err := strconv.Atoi(pidStr)
	if err != nil {
		zap.L().Error("invalid param", zap.Error(err))
		ResponseError(c, CodeInvalidParam)
		return
	}

	// 尝试获取当前用户ID（如果已登录）
	var userID int64
	// 这里不强制要求登录，所以忽略错误
	if uid, err := getCurrentUser(c); err == nil {
		userID = uid
	}

	// logic处理,根据帖子的id来查询帖子的具体数据
	data, err := logic.GetPostByID(int64(pid), userID)
	if err != nil {
		zap.L().Error("logic.get post by id failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}

	// 3返回响应
	ResponseSuccess(c, data)
}

// GetPostListHandler 获取所有帖子列表的处理函数
func GetPostListHandler(c *gin.Context) {
	// 获取分页参数
	page, size := GetPageInfo(c)

	// 获取数据
	data, err := logic.GetPostList(page, size)
	if err != nil {
		zap.L().Error("logic.GetPostList failed", zap.Error(err))
		ResponseError(c, CodePostNotExist)
		return
	}

	// 返回响应
	ResponseSuccess(c, data)
}

// GetPostListHandler2 根据顺序来查找所有帖子
func GetPostListHandler2(c *gin.Context) {
	// 处理请求参数, 默认值如下
	p := &models.ParamPostList{
		Page:  1,
		Size:  10,
		Order: models.OrderTime,
	}

	if err := c.ShouldBind(p); err != nil {
		zap.L().Error("GetPostListHandler2 param failed", zap.Error(err))
		ResponseError(c, CodeInvalidParam)
		return
	}

	// 2 获取帖子数据
	data, err := logic.GetPostListNew(p)
	if err != nil {
		zap.L().Error("logic.GetPostList failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}

	// 3 返回响应
	ResponseSuccess(c, data)
}
