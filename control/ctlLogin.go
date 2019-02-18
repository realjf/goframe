package control

import (
	"kboard/config"
	"kboard/template"
	"net/http"
)

type CtlLogin struct {
	Control
}

func NewCtlLogin(config config.IConfig, w http.ResponseWriter, r *http.Request) *CtlLogin {
	return &CtlLogin{
		Control{
			Config:    config,
			TplEngine: template.NewTplEngine(w, r),
			Module:    "login",
			Actions:   map[string]func(){},
			R:         r,
			W:         w,
		},
	}
}

func (this *CtlLogin) Index() {
	this.Display("login")
}
