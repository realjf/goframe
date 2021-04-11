package goframe

import (
	"flag"

	"github.com/realjf/goframe/server"
)

var (
	// flag启动参数
	ConfigPath string
	CaCertPath string
	CaKeyPath  string

	DefaultServer *server.Server
)

func init() {
	// 启动参数处理
	// 配置文件路径
	flag.StringVar(&ConfigPath, "config-path", "config/config.yaml", "--config-path, specify config file path;default path is config/conf.toml")
	flag.StringVar(&CaCertPath, "ca-cert", "config/ca.cer", "--ca-cert, specify ca-cert file path;default path is config/ca.cer")
	flag.StringVar(&CaKeyPath, "ca-key", "config/ca.key", "--ca-key, specify ca-key file path;default path is config/ca.key")

	flag.Parse()
}

func NewServer(configPath string) *server.Server {
	return &server.Server{
		ConfigPath: configPath,
	}
}

func NewDefaultServer() *server.Server {
	if DefaultServer != nil {
		return DefaultServer
	}
	DefaultServer = NewServer("./config/config.yaml")
	return DefaultServer
}
