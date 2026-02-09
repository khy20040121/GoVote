package redis

import (
	"fmt"

	"github.com/go-redis/redis"
	"github.com/spf13/viper"
)

var (
	client *redis.Client
	Nil    = redis.Nil
)

// Init 初始化连接
func Init() (err error) {
	client = redis.NewClient(&redis.Options{
		Addr:         fmt.Sprintf("%s:%d", viper.GetString("redis.host"), viper.GetInt("redis.port")),
		Password:     viper.GetString("redis.password"), // no password set
		DB:           viper.GetInt("redis.db_name"),     // use default DB
		PoolSize:     viper.GetInt("redis.pool_size"),
		MinIdleConns: viper.GetInt("redis.min_idle_conns"),
	})

	_, err = client.Ping().Result()
	if err != nil {
		return err
	}
	return nil
}

func Close() {
	_ = client.Close()
}
