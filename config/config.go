package config

import (
	"fmt"
	"io/ioutil"
	"sync"

	"github.com/BurntSushi/toml"
	"goframe/exception"
	"goframe/utils"
)

type IConfig interface {
	LoadConfigFile(path string) IConfig
	ReloadConfigFile()
	SetTSL(tsl *ClusterTSL) error
	GetTSL() *ClusterTSL
	GetAddress() string
	IsHttps() bool // 是否开启https
	IsLog() bool
	IsAuth() bool // 是否开启鉴权
	IsHttp2() bool
	GetHttpVersion() string
	GetConfigData() ConfigData
}

type ClusterTSL struct {
	Cert string
	Key  string
}

type Config struct {
	Path string
	Data ConfigData
	Once sync.Once // 实现单例模式
	Lock sync.RWMutex
}

type ConfigData struct {
	Cluster struct {
		Host        string
		Port        int
		Log         bool   // 日志记录
		Auth        bool   // 鉴权
		HttpVersion string `toml:"httpVersion",yaml:"httpVersion"` // 1.0, 1.1, 2.0
		Openssl     bool   // https
		TLS         struct {
			Cert string
			Key  string
		}
	}
	Mysql struct {
		Host         string
		Port         int
		Username     string
		Password     string
		Dbname       string
		Charset      string
		MaxOpenConns int `toml:"maxOpenConns",yaml:"maxOpenConns"`
	}
	Memcache struct {
		Host string
		Port int
	}
	Redis struct {
		Host    string
		Port    int
		Timeout int
	}
}

// load config file
// singleton
func NewConfig() IConfig {
	return &Config{
		Path: "",
		Data: ConfigData{},
		Once: sync.Once{},
		Lock: sync.RWMutex{},
	}
}

// error code 1000 ~ 1200
func (c *Config) LoadConfigFile(path string) IConfig {
	fmt.Println("loading config file [path: ", path, "]")
	c.Once.Do(func() {
		if path == "" {
			exception.CheckError(exception.NewError("config-path is empty"), 1000)
		}
		c.Path = path
		// 检查文件是否存在
		fileData, err := ioutil.ReadFile(path)
		if err != nil || len(fileData) <= 0 {
			exception.CheckError(exception.NewError("read toml config file error"), 0)
		}
		if _, err := toml.DecodeFile(path, &c.Data); err != nil {
			exception.CheckError(err, 1001)
		}
	})

	return c
}

func (c *Config) GetConfigData() ConfigData {
	return c.Data
}

// 重新加载配置文件
func (c *Config) ReloadConfigFile() {
	fmt.Println("reloading config file...")
	c.Lock.RLock()
	defer c.Lock.RUnlock()
	c.Once.Do(func() {
		c.Lock.Lock()
		defer c.Lock.Unlock()
		if _, err := toml.DecodeFile(c.Path, &c.Data); err != nil {
			exception.CheckError(err, 1001)
		}
	})
}

func (c *Config) SetTSL(tsl *ClusterTSL) error {
	if tsl.Cert == "" || tsl.Key == "" {
		return exception.NewError("server tsl contain invalid value")
	}
	c.Data.Cluster.TLS.Key = tsl.Key
	c.Data.Cluster.TLS.Cert = tsl.Cert
	return nil
}

func (c *Config) GetAddress() string {
	if c.Data.Cluster.Host == "" {
		exception.CheckError(exception.NewError("server host is empty"), 1004)
	}
	port := c.Data.Cluster.Port
	if port <= 0 || port > 65535 {
		exception.CheckError(exception.NewError("server port is invalid"), 1004)
	}
	return c.Data.Cluster.Host + ":" + utils.ToString(port)
}

func (c *Config) GetTSL() *ClusterTSL {
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

func (c *Config) IsHttps() bool {
	return c.Data.Cluster.Openssl
}

func (c *Config) IsLog() bool {
	return c.Data.Cluster.Log
}

func (c *Config) IsAuth() bool {
	return c.Data.Cluster.Auth
}

func (c *Config) IsHttp2() bool {
	if c.Data.Cluster.HttpVersion == "2.0" {
		return true
	}
	return false
}

func (c *Config) GetHttpVersion() string {
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
