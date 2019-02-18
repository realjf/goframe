package control

import (
	"kboard/config"
	"kboard/template"
	"net/http"
)

type CtlIndex struct {
	Control
}

func NewCtlIndex(config config.IConfig, w http.ResponseWriter, r *http.Request) *CtlIndex {
	return &CtlIndex{
		Control{
			Config:    config,
			TplEngine: template.NewTplEngine(w, r),
			Module:    "index",
			Actions:   map[string]func(){},
			R:         r,
			W:         w,
		},
	}
}

func (this *CtlIndex) Index() {
	this.Display("index")
}
