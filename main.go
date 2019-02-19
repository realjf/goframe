package main

import (
	"flag"
	"goframe/config"
	"goframe/exception"
	"goframe/middleware"
	"goframe/router"
	"golang.org/x/net/http2"
	"log"
	"net/http"
	"time"
	"goframe/db"
)

var (
	Config             config.IConfig
	NotifyReloadConfig chan int

	// flag启动参数
	ConfigPath string
	CaCertPath string
	CaKeyPath  string
)

func init() {
	// 启动参数处理
	// 配置文件路径
	flag.StringVar(&ConfigPath, "config-path", "config/config.yaml", "--config-path, specify config file path;default path is config/conf.toml")
	flag.StringVar(&CaCertPath, "ca-cert", "config/ca.cer", "--ca-cert, specify ca-cert file path;default path is config/ca.cer")
	flag.StringVar(&CaKeyPath, "ca-key", "config/ca.key", "--ca-key, specify ca-key file path;default path is config/ca.key")

	flag.Parse()

	// init config
	Config = config.NewConfig().LoadConfigFile(ConfigPath)

	// watch config file to reload
	NotifyReloadConfig = make(chan int, 1)
	go func() {
		for {
			<-NotifyReloadConfig
			Config.ReloadConfigFile()
		}
	}()

	// init log
	middleware.Logger = middleware.NewLogger().Init()

	// init db、cache、control and so on
	db.NewMysql(Config).Init()
}

func main() {
	r := router.NewRouter(Config, middleware.Logger).InitRouter()
	log.Println("Listen On", Config.GetAddress())
	server := http.Server{
		Addr:         Config.GetAddress(),
		Handler:      r.Router,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	// turn http/2.0 on
	if Config.IsHttp2() {
		err := http2.ConfigureServer(&server, &http2.Server{})
		exception.CheckError(err, 11)
	}
	middleware.Logger.Logger.Info(Config.GetHttpVersion())

	if Config.IsHttps() {
		ca := Config.GetTSL()
		if CaKeyPath != "" && CaCertPath != "" {
			middleware.Logger.Logger.Fatal(server.ListenAndServeTLS(CaCertPath, CaKeyPath))
		} else {
			middleware.Logger.Logger.Fatal(server.ListenAndServeTLS(ca.Cert, ca.Key))
		}
	} else {
		middleware.Logger.Logger.Fatal(server.ListenAndServe())
	}
}
