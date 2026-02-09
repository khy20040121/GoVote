package setting

import (
	"fmt"
	"strings"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

func Init() (err error) {

	viper.SetConfigFile("./config/config.yaml") // 设置文件的路径
	err = viper.ReadInConfig()                  // 读取配置信息

	if err != nil { // 读取配置信息失败
		fmt.Printf("viper.ReadInConfig failed, err:%v\n", err)
		return
	}

	// 支持环境变量 可以在服务器配置上面自动配置环境
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	// 监视, 如果配置文件发生了改变,启动回调函数,重新绑定配置到结构体上
	viper.WatchConfig()
	viper.OnConfigChange(func(in fsnotify.Event) {
		fmt.Println("配置文件修改了...")
	})
	return
}
