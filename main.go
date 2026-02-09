package main

import (
	"bluebell/dao/mysql"
	"bluebell/dao/redis"
	"bluebell/logger"
	"bluebell/pkg/snowflake"
	"bluebell/router"
	"bluebell/setting"
	"fmt"

	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func main() {

	// 第一步 一定是加载配置文件yaml 才能够去做后续的操作
	if err := setting.Init(); err != nil {
		fmt.Println("load config failed, err:", zap.Error(err))
		return
	}

	// 注意只有当配置文件生效之后才调用其他函数的初始化函数
	// 初始化logger
	logger.Init(viper.GetString("app.mode"))

	// 连接mysql, 最后记得关闭数据库
	if err := mysql.Init(); err != nil {
		zap.L().Error("init mysql failed, err:%v\n", zap.Error(err))
		return
	}
	defer mysql.Close() // 程序退出关闭数据库连接

	// 连接redis, 并记得关闭数据库
	if err := redis.Init(); err != nil {
		zap.L().Error("init redis failed, err:%v\n", zap.Error(err))
		return
	}
	defer redis.Close()

	// 初始化雪花算法
	if err := snowflake.Init(viper.GetString("app.start_time"), viper.GetInt64("app.machine_id")); err != nil {
		zap.L().Error("init snowflake failed", zap.Error(err))
		return
	}

	// 注册路由
	r := router.SetupRouter(viper.GetString("app.mode"))
	err := r.Run(fmt.Sprintf(":%d", viper.GetInt("app.port")))
	if err != nil {
		zap.L().Error("start server failed", zap.Error(err))
		return
	}
}
