package control

import (
	"fmt"
	"log"
	"net/http"

	"github.com/realjf/goframe/config"
	"github.com/realjf/goframe/pkg/exception"
	"github.com/realjf/goframe/pkg/utils"
	"github.com/realjf/goframe/template"
)

type IControl interface {
	Register(string, func()) *Control
	Run(string)
	ResponseWithHeader(int, interface{}, string)
	Response(int, interface{}, string)
	Display(string)
}

type Control struct {
	Config    config.IConfig
	TplEngine *template.TplEngine
	Module    string
	Namespace string
	Actions   map[string]func()
	W         http.ResponseWriter
	R         *http.Request
	Header    map[string]string
}

func NewControl(config config.IConfig, w http.ResponseWriter, r *http.Request) *Control {
	return &Control{
		Config:    config,
		TplEngine: template.NewTplEngine(w, r),
		Module:    "base",
		Namespace: "control",
		Actions:   map[string]func(){},
		R:         r,
		W:         w,
		Header: map[string]string{
			"Access-Control-Allow-Origin":   "*",
			"Access-Control-Allow-Methods":  "*",
			"Access-Control-Allow-Headers":  "Content-Type,Access-Token,X-Access-Token,X-Session-Token",
			"Access-Control-Expose-Headers": "*",
		},
	}
}

func (c *Control) GetString(name string) string {
	return utils.ToString(c.R.URL.Query().Get(name))
}

func (c *Control) GetInt(name string) int {
	return utils.ToInt(c.R.URL.Query().Get(name))
}

func (c *Control) PostString(name string) string {
	return utils.ToString(c.R.FormValue(name))
}

func (c *Control) PostInt(name string) int {
	return utils.ToInt(c.R.FormValue(name))
}

func (c *Control) Register(action string, f func()) *Control {
	if c.Actions == nil {
		c.Actions = map[string]func(){}
	}
	if c.Module == "" {
		exception.CheckError(exception.NewError("error: control is empty!"), 999)
	}
	c.Actions[action] = f
	return c
}

func (c *Control) Run(action string) {
	// 注册全局变量
	if c.TplEngine.TplData["GModule"] == nil || c.TplEngine.TplData["GModule"] == "" {
		c.TplEngine.TplData["GModule"] = c.Module
	}
	if c.TplEngine.TplData["GAction"] == nil || c.TplEngine.TplData["GAction"] == "" {
		c.TplEngine.TplData["GAction"] = action
	}
	// 检查action方法是否存在
	if f, ok := c.Actions[action]; !ok {
		if defaultFunc, ok1 := c.Actions["index"]; !ok1 {
			fmt.Fprintln(c.TplEngine.W, "404 page not found!")
			log.Println("404")
		} else {
			c.TplEngine.TplData["GAction"] = "index"
			defaultFunc()
		}
	} else {
		// run action
		f()
	}
}

func (c *Control) Index() {
	fmt.Fprintln(c.TplEngine.W, "hello world, this is default index")
}

func (c *Control) ResponseWithHeader(code int, result interface{}, message string) {
	c.TplEngine.ResponseWithHeader(code, result, message, c.Header)
}

func (c *Control) Response(code int, result interface{}, message string) {
	c.TplEngine.Response(code, result, message)
}

func (c *Control) Display(tpl string) {
	c.TplEngine.Display(tpl)
}
