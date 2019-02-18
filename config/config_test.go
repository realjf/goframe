package config

import (
	"testing"
)

func TestConfig_LoadConfigFile(t *testing.T) {
	conf := NewConfig().LoadConfigFile("./conf.toml")
	t.Errorf("%+v", conf.GetConfigData())
}
