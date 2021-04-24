package router

import (
	"net/http"

	"github.com/realjf/goframe/config"
	"github.com/realjf/goframe/internal/control"

	"github.com/gorilla/mux"
)

// 注册路由
func RegisterUrl(r *Router) {
	r.Router.HandleFunc("/", C_DefaultHandler(r.Config))

	// control下的路由
	r.Router.HandleFunc("/login/{action:[a-z]+}", C_LoginHandler(r.Config))
	r.Router.HandleFunc("/index/{action:[a-z]+}", C_IndexHandler(r.Config))

}

// control下的路由处理handler在此处理
func C_LoginHandler(c config.IConfig) (f func(http.ResponseWriter, *http.Request)) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		action := mux.Vars(r)["action"]
		c := control.NewCtlLogin(c, w, r)
		c.Register("index", c.Index).Run(action)
	}

	return handler
}

func C_IndexHandler(c config.IConfig) (f func(http.ResponseWriter, *http.Request)) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		action := mux.Vars(r)["action"]
		c := control.NewCtlIndex(c, w, r)
		c.Register("index", c.Index).Run(action)
	}

	return handler
}

func C_DefaultHandler(c config.IConfig) (f func(http.ResponseWriter, *http.Request)) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		action := mux.Vars(r)["action"]
		c := control.NewCtlDefault(c, w, r)
		c.Register("index", c.Index).Run(action)
	}

	return handler
}
