package redis

import (
	"github.com/gogf/gf/container/gmap"
	"github.com/gogf/gf/errors/gerror"
	"github.com/gogf/gf/text/gregex"
	"github.com/gogf/gf/text/gstr"
	"github.com/gogf/gf/util/gconv"
	"time"
)

const (
	DEFAULT_GROUP_NAME = "default"
	DEFAULT_REDIS_PORT = 6379
)

var (
	// Configuration groups.
	configs = gmap.NewStrAnyMap(true)
)

func SetConfigByStr(str string, name ...string) error {
	group := DEFAULT_GROUP_NAME
	if len(name) > 0 {
		group = name[0]
	}
	config, err := ConfigFromStr(str)
	if err != nil {
		return err
	}
	configs.Set(group, config)
	instances.Remove(group)
	return nil
}

func GetConfig(name ...string) (config Config, ok bool) {
	group := DEFAULT_GROUP_NAME
	if len(name) > 0 {
		group = name[0]
	}
	if v := configs.Get(group); v != nil {
		return v.(Config), true
	}
	return Config{}, false
}

func RemoveConfig(name ...string) {
	group := DEFAULT_GROUP_NAME
	if len(name) > 0 {
		group = name[0]
	}
	configs.Remove(group)
	instances.Remove(group)
}

func ConfigFromStr(str string) (config Config, err error) {
	array, _ := gregex.MatchString(`([^:]+):*(\d*),{0,1}(\d*),{0,1}(.*)\?(.+)`, str)
	if len(array) == 6 {
		parse, _ := gstr.Parse(array[5])
		config = Config{
			Host: array[1],
			Port: gconv.Int(array[2]),
			Db:   gconv.Int(array[3]),
			Passwd: array[4],
		}
		if config.Port == 0 {
			config.Port = DEFAULT_REDIS_PORT
		}
		if v, ok := parse["maxIdle"]; ok {
			config.MaxIdle = gconv.Int(v)
		}
		if v, ok := parse["maxActive"]; ok {
			config.MaxActive = gconv.Int(v)
		}
		if v, ok := parse["idleTimeout"]; ok {
			config.IdleTimeout = gconv.Duration(v) * time.Second
		}
		if v, ok := parse["maxConnLifetime"]; ok {
			config.MaxConnLifetime = gconv.Duration(v) * time.Second
		}
		return
	}
	array, _ = gregex.MatchString(`([^:]+):*(\d*),{0,1}(\d*),{0,1}(.*)`, str)
	if len(array) == 5 {
		config = Config{
			Host: array[1],
			Port: gconv.Int(array[2]),
			Db:   gconv.Int(array[3]),
			Passwd: array[4],
		}
		if config.Port == 0 {
			config.Port = DEFAULT_REDIS_PORT
		}
	} else {
		err = gerror.Newf(`invalid redis configuration: "%s"`, str)
	}
	return
}

func ClearConfig() {
	configs.Clear()
	instances.Clear()
}


