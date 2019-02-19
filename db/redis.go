package db

import (
	"github.com/gomodule/redigo/redis"
	"goframe/config"
	"goframe/utils"
)

var RedisClient *redis.Conn

type RedisDriver struct {
	Host string
	Port string
}

func NewRedis(config config.IConfig) *RedisDriver {
	configData := config.GetConfigData()
	return &RedisDriver{
		Host: configData.Redis.Host,
		Port: utils.ToString(configData.Redis.Port),
	}
}

func (r *RedisDriver) Init() {

}


