package db

import (
	"fmt"
	"github.com/pangudashu/memcache"
	"goframe/config"
	"goframe/exception"
	"goframe/middleware"
	"goframe/utils"
	"time"
)

var McClient *memcache.Memcache

type McDriver struct {
	Host string
	Port string
}

func NewMemcache(config config.IConfig) *McDriver {
	configData := config.GetConfigData()
	return &McDriver{
		Host: configData.Memcache.Host,
		Port: utils.ToString(configData.Memcache.Port),
	}
}

func (mc *McDriver) Init() {
	if mc.Host == "" || mc.Port == "" {
		exception.CheckError(exception.NewError("memcache config is error"), 3001)
	}
	serv1 := &memcache.Server{Address: fmt.Sprintf("%s:%s", mc.Host, mc.Port), Weight: 50}

	mcClient, err := memcache.NewMemcache([]*memcache.Server{serv1})
	exception.CheckError(err, 3000)
	// 设置是否自动剔除无法连接的server，默认不开启(建议开启)
	// 如果开启此选项被踢除的server如果恢复正常将会再次被加入server列表
	McClient = mcClient
	McClient.SetRemoveBadServer(true)

	McClient.SetTimeout(time.Second*2, time.Second, time.Second)

	middleware.Logger.Logger.Info("init memcache...")
}
