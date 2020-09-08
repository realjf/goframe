package control

import (
	"net/http"

	"github.com/realjf/goframe/config"
	"github.com/realjf/goframe/template"
)

type CtlDefault struct {
	Control
}

func NewCtlDefault(config config.IConfig, w http.ResponseWriter, r *http.Request) *CtlDefault {
	return &CtlDefault{
		Control{
			Config:    config,
			TplEngine: template.NewTplEngine(w, r),
			Module:    "default",
			Actions:   map[string]func(){},
			R:         r,
			W:         w,
		},
	}
}

func (this *CtlDefault) Index() {
	this.Response(100, "", "ok")
}
