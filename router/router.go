package router

import (
	"flag"
	"github.com/gorilla/mux"
	"github.com/realjf/goframe/config"
	"github.com/realjf/goframe/middleware"
	"net/http"
)

type Router struct {
	Router *mux.Router
	Config config.IConfig
	Logger *middleware.Log
}

func NewRouter(Config config.IConfig, Logger *middleware.Log) *Router {
	return &Router{
		Router: mux.NewRouter(),
		Config: Config,
		Logger: Logger,
	}
}

// register url
func (r *Router) InitRouter() *Router {
	var dir string
	flag.StringVar(&dir, "dir", "assets", "")
	flag.Parse()

	// static files
	r.Router.PathPrefix("/assets/").Handler(http.StripPrefix("/assets/", http.FileServer(http.Dir(dir))))

	// match url
	RegisterUrl(r)
	RegisterApi(r)

	// logs
	if r.Config.IsLog() {
		r.Router.Use(middleware.AccessLogger)
	}

	// authentication
	if r.Config.IsAuth() {
		amw := middleware.NewAuthentication()
		amw.Populate()
		r.Router.Use(amw.Middleware)
	}

	// safe handler
	r.Router.Use(middleware.SafeHandler)

	// jwt
	//r.Router.Use(middleware.JWT)

	return r
}
