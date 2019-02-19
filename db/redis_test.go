package db

import (
	"goframe/config"
	"testing"
	"time"
)

func TestNewRedis(t *testing.T) {
	config := config.NewConfigYaml().LoadConfigFile("../config/config.yaml")
	NewRedis(config).Init()

	err := RedisClient.Set("hello", "world", time.Second * time.Duration(3600)).Err()
	if err != nil {
		t.Fatal(err.Error())
	}
	res, err := RedisClient.Get("hello").Result()
	if err != nil {
		t.Fatal(err.Error())
	}
	t.Fatal(res)
}
