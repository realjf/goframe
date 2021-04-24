package control

import (
	"net/http"

	"github.com/realjf/goframe/config"
	"github.com/realjf/goframe/internal/template"
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
