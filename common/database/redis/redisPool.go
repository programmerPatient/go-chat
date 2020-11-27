package redis

import (
	"../../../config"
	"github.com/garyburd/redigo/redis"
	"strconv"
)

var Pool *RedisPool

type RedisPool struct {
	MaxCount int //最大连接数
	MinCount int //最小创建数
	pool *redis.Pool //连接池
}



func New() *RedisPool {
	max,_ :=strconv.Atoi(config.RedisConf["maxCount"])
	min,_ :=strconv.Atoi(config.RedisConf["minCount"])
	addr := config.RedisConf["host"]+":"+config.RedisConf["port"]
	pool := &redis.Pool{
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp",addr)
		},
		MaxIdle:         min,
		MaxActive:       max,
		IdleTimeout:     300,
	}
	return &RedisPool{
		MaxCount: max,
		MinCount: min,
		pool:     pool,
	}
}

/**
获取连接
 */
func (r *RedisPool) Get() *RedisDb{
	return &RedisDb{
		C:r.pool.Get(),
	}
}




