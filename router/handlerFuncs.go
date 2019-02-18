package router

import (
	"kboard/api"
	"kboard/config"
	"kboard/control"
	"net/http"

	"github.com/gorilla/mux"
)

// 注册路由
func UrlRegister(r *Router) {
	r.Router.HandleFunc("/", C_DefaultHandler(r.Config))

	// api的路由特殊处理
	r.Router.HandleFunc("/api/user/{action:[a-z]+}", I_UserHandler(r.Config))
	r.Router.HandleFunc("/api/node/{action:[a-z]+}", I_NodeHandler(r.Config))
	r.Router.HandleFunc("/api/order/{action:[a-z]+}", I_OrderHandler(r.Config))
	r.Router.HandleFunc("/api/team/{action:[a-z]+}", I_TeamHandler(r.Config))
	r.Router.HandleFunc("/api/login/{action:[a-z]+}", I_LoginHandler(r.Config))

	// control下的路由
	r.Router.HandleFunc("/login/{action:[a-z]+}", C_LoginHandler(r.Config))
	r.Router.HandleFunc("/index/{action:[a-z]+}", C_IndexHandler(r.Config))

}

// api下的路由处理handler在此处理
func I_UserHandler(c config.IConfig) (f func(http.ResponseWriter, *http.Request)) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		action := mux.Vars(r)["action"]
		i := api.NewIUser(c, w, r)
		i.Register("index", i.Index).Run(action)
	}

	return handler
}

func I_NodeHandler(c config.IConfig) (f func(http.ResponseWriter, *http.Request)) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		action := mux.Vars(r)["action"]
		i := api.NewINode(c, w, r)
		i.Register("scale", i.Scale)
		i.Register("index", i.Index).Run(action)
	}

	return handler
}

func I_OrderHandler(c config.IConfig) (f func(http.ResponseWriter, *http.Request)) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		action := mux.Vars(r)["action"]
		i := api.NewIOrder(c, w, r)
		i.Register("list", i.List)
		i.Register("index", i.Index).Run(action)
	}

	return handler
}

func I_TeamHandler(c config.IConfig) (f func(http.ResponseWriter, *http.Request)) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		action := mux.Vars(r)["action"]
		i := api.NewITeam(c, w, r)

		i.Register("index", i.Index).Run(action)
	}

	return handler
}

func I_LoginHandler(c config.IConfig) (f func(http.ResponseWriter, *http.Request)) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		action := mux.Vars(r)["action"]
		i := api.NewILogin(c, w, r)

		i.Register("index", i.Index).Run(action)
	}

	return handler
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
