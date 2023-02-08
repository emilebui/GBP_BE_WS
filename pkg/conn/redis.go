package conn

import (
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
)

func GetRedisConn(conf *viper.Viper) *redis.Client {

	addr := conf.GetString("redis_addr")
	pwd := conf.GetString("redis_pwd")
	db := conf.GetInt("redis_db")

	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: pwd, // no password set
		DB:       db,  // use default DB
	})

	return rdb

}
