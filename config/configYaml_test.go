package config

import "testing"

func TestNewConfigYaml(t *testing.T) {
	conf := NewConfigYaml().LoadConfigFile("./config.yml")
	t.Errorf("%+v", conf.GetConfigData())
}
