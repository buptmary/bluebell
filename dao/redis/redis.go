package redis

import (
	"bluebell/settings"
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
)

// 声明一个全局的rdb变量
var (
	rdb *redis.Client
	ctx = context.Background()
)

// Init 初始化连接
func Init(redisConfig *settings.RedisConfig) (err error) {

	// 创建一个Redis客户端实例
	rdb = redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%d",
			redisConfig.Host,
			redisConfig.Port,
		),
		Password: redisConfig.Password,
		DB:       redisConfig.DB,
	})

	_, err = rdb.Ping(ctx).Result()
	if err != nil {
		return err
	}
	return err
}

func Close() {
	_ = rdb.Close()
}
