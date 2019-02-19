package db

import (
	"testing"
	config2 "goframe/config"
)

func TestNewRedis(t *testing.T) {
	config := config2.NewConfigYaml().LoadConfigFile("../config/config.yaml")


}

