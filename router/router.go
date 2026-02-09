package router

import (
	"bluebell/controller"
	"bluebell/logger"
	"bluebell/middlewares"

	"github.com/gin-gonic/gin"
)

func SetupRouter(mode string) *gin.Engine {
	if mode == gin.ReleaseMode {
		gin.SetMode(gin.ReleaseMode)
	}
	// 默认为debug模式
	r := gin.New()
	r.Use(logger.GinLogger(), logger.GinRecovery(true))

	v1 := r.Group("/api/v1")

	// 注册
	v1.POST("/signup", controller.SignUpHandler)
	// 登录
	v1.POST("/login", controller.LoginHandler)

	// 根据时间或分数获取帖子列表
	v1.GET("/posts2", controller.GetPostListHandler2)
	v1.GET("/posts", controller.GetPostListHandler)
	v1.GET("/community", controller.CommunityHandler)
	v1.GET("/community/:id", controller.CommunityDetailHandler)

	// 使用 OptionalJWTAuthMiddleware，让 GetPostDetailHandler 可以获取到 userID
	v1.GET("/post/:id", middlewares.OptionalJWTAuthMiddleware(), controller.GetPostDetailHandler)

	v1.Use(middlewares.JWTAuthMiddleware()) // 应用JWT认证中间件

	{
		// 发表帖子
		v1.POST("/post", controller.CreatePostHandler)

		// 为帖子投票
		v1.POST("/vote", controller.PostVoteHandler)
	}

	return r
}
