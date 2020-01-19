package server

import (
	"errors"
	"goframe/config"
	"goframe/db"
	"goframe/db/memcache"
	"goframe/db/mysql"
	"goframe/middleware"
	"goframe/router"
	"golang.org/x/net/http2"
	"log"
	"net/http"
	"strings"
	"time"
)

type Server struct {
	ConfigPath string
	server http.Server
	Config config.IConfig
	Router *router.Router
}

func (s *Server) SetAddress(addr string) {
	s.server.Addr = addr
}

func (s *Server) SetTimeout(timeout int) {
	s.server.ReadTimeout = time.Duration(timeout) * time.Second
	s.server.WriteTimeout = time.Duration(timeout) * time.Second
	s.server.IdleTimeout = time.Duration(timeout) * time.Second
}

func (s *Server) Run() {
	s.InitConfig(s.ConfigPath)
	s.InitLogger()

	s.InitDb()
	s.InitRedis()
	s.InitMemcached()

	s.InitRouter()
	s.InitServer()

	s.Start()
}

func (s *Server) InitServer() {
	s.server = http.Server{
		Addr:         s.Config.GetAddress(),
		Handler:      s.Router.Router,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	if s.Config.IsHttp2() {
		err := http2.ConfigureServer(&s.server, &http2.Server{})
		if err != nil {
			log.Panic(err)
		}
	}
}

func (s *Server) Start() {
	if s.Config.IsHttps() {
		ca := s.Config.GetTSL()
		s.server.ListenAndServeTLS(ca.Cert, ca.Key)
	} else {
		s.server.ListenAndServe()
	}
}

// 初始化日志
func (s *Server) InitLogger() {
	middleware.Logger = middleware.NewLogger().Init()
}

// 初始化路由
func (s *Server) InitRouter() {
	s.Router = router.NewRouter(s.Config, middleware.Logger).InitRouter()
}

// 初始化配置文件
func (s *Server) InitConfig(confPath string) {
	if confPath == "" {
		log.Panic(errors.New("配置文件为空"))
	}

	// 确认文件后缀名，根据文件后缀名确定加载的文件
	if strings.HasSuffix(confPath, "yaml") || strings.HasSuffix(confPath, "yml") {
		s.Config = config.NewConfigYaml().LoadConfigFile(confPath)
	}else if strings.HasSuffix(confPath, "toml") {
		s.Config = config.NewConfig().LoadConfigFile(confPath)
	}
}

// 初始化db
func (s *Server) InitDb() {
	db.NewRedis(s.Config).Init()
}

// 初始化redis
func (s *Server) InitRedis() {
	mysql.NewMysql(s.Config).Init()
}

// 初始化memcache
func (s *Server) InitMemcached() {
	memcache.NewMemcache(s.Config).Init()
}



