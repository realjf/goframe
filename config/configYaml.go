package config

import (
	"fmt"
	"os"
	"sync"

	"gopkg.in/yaml.v2"

	"github.com/realjf/goframe/pkg/exception"
	"github.com/realjf/goframe/pkg/utils"
)

type ConfigYaml struct {
	Path string
	Data ConfigData
	Once sync.Once // 实现单例模式
	Lock sync.RWMutex
}

func NewConfigYaml() IConfig {
	return &ConfigYaml{
		Path: "",
		Data: ConfigData{},
		Once: sync.Once{},
		Lock: sync.RWMutex{},
	}
}

func (c *ConfigYaml) LoadConfigFile(path string) IConfig {
	fmt.Println("loading config file [path: ", path, "]")
	c.Once.Do(func() {
		if path == "" {
			exception.CheckError(exception.NewError("config-path is empty"), 1000)
		}
		c.Path = path
		// 检查文件是否存在
		fileData, err := os.ReadFile(path)
		if err != nil || len(fileData) <= 0 {
			exception.CheckError(exception.NewError("read yaml config file error"), 0)
		}
		if err := yaml.Unmarshal(fileData, &c.Data); err != nil {
			exception.CheckError(err, 1001)
		}
	})

	return c
}

func (c *ConfigYaml) GetConfigData() ConfigData {
	return c.Data
}

// 重新加载配置文件
func (c *ConfigYaml) ReloadConfigFile() {
	fmt.Println("reloading config file...")
	c.Lock.RLock()
	defer c.Lock.RUnlock()
	c.Once.Do(func() {
		c.Lock.Lock()
		defer c.Lock.Unlock()
		fileData, err := os.ReadFile(c.Path)
		if err != nil || len(fileData) <= 0 {
			exception.CheckError(exception.NewError("read yaml config file error"), 0)
		}
		if err := yaml.Unmarshal(fileData, &c.Data); err != nil {
			exception.CheckError(err, 1001)
		}
	})
}

func (c *ConfigYaml) SetTSL(tsl *ClusterTSL) error {
	if tsl.Cert == "" || tsl.Key == "" {
		return exception.NewError("server tsl contain invalid value")
	}
	c.Data.Cluster.TLS.Key = tsl.Key
	c.Data.Cluster.TLS.Cert = tsl.Cert
	return nil
}

func (c *ConfigYaml) GetAddress() string {
	if c.Data.Cluster.Host == "" {
		exception.CheckError(exception.NewError("server host is empty"), 1004)
	}
	port := c.Data.Cluster.Port
	if port <= 0 || port > 65535 {
		exception.CheckError(exception.NewError("server port is invalid"), 1004)
	}
	return c.Data.Cluster.Host + ":" + utils.ToString(port)
}

func (c *ConfigYaml) GetTSL() *ClusterTSL {
	cert := c.Data.Cluster.TLS.Cert
	key := c.Data.Cluster.TLS.Key
	if cert == "" || key == "" {
		exception.CheckError(exception.NewError("cert or key is empty"), 1005)
	}
	return &ClusterTSL{
		Cert: cert,
		Key:  key,
	}
}

func (c *ConfigYaml) IsHttps() bool {
	return c.Data.Cluster.Openssl
}

func (c *ConfigYaml) IsLog() bool {
	return c.Data.Cluster.Log
}

func (c *ConfigYaml) IsAuth() bool {
	return c.Data.Cluster.Auth
}

func (c *ConfigYaml) IsHttp2() bool {
	if c.Data.Cluster.HttpVersion == "2.0" {
		return true
	}
	return false
}

func (c *ConfigYaml) GetHttpVersion() string {
	var httpVersion string
	switch c.Data.Cluster.HttpVersion {
	case "1.0":
		fallthrough
	case "1.1":
		fallthrough
	case "2.0":
		httpVersion = c.Data.Cluster.HttpVersion
	default:
		httpVersion = "1.1"
	}
	return "HTTP/" + httpVersion
}
