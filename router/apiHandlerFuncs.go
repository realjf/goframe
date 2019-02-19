package router

import (
	"goframe/api"
	"goframe/config"
	"net/http"

	"github.com/gorilla/mux"
)


func RegisterApi(r *Router) {
	// api的路由特殊处理
	r.Router.HandleFunc("/api/user/{action:[a-z]+}", I_UserHandler(r.Config))
	r.Router.HandleFunc("/api/login/{action:[a-z]+}", I_LoginHandler(r.Config))


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

func I_LoginHandler(c config.IConfig) (f func(http.ResponseWriter, *http.Request)) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		action := mux.Vars(r)["action"]
		i := api.NewILogin(c, w, r)

		i.Register("index", i.Index).Run(action)
	}

	return handler
}

