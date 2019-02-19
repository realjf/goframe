package db

import (
	"goframe/config"
	"testing"
)

func TestNewMemcache(t *testing.T) {
	config := config.NewConfigYaml().LoadConfigFile("../config/config.yaml")
	NewMemcache(config).Init()

	_, err := McClient.Set("hello", "world", 3600)
	if err != nil {
		t.Fatal(err)
	}
	res1,_,_ := McClient.Get("hello")
	t.Fatal(res1)
}
