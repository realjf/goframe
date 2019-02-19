package api

import (
	"fmt"
	"goframe/config"
	"goframe/exception"
	"goframe/template"
	"log"
	"net/http"
	"goframe/utils"
)

type iApi interface {
	Register(string, func()) *IApi
	Run(string)
	ResponseWithHeader(int, interface{}, string)
	Response(int, interface{}, string)
}

type IApi struct {
	Config    config.IConfig
	TplEngine *template.TplEngine
	Module    string
	Actions   map[string]func()
	W         http.ResponseWriter
	R         *http.Request
	Header    map[string]string
	Namespace string
}

func NewIApi(config config.IConfig, w http.ResponseWriter, r *http.Request) *IApi {
	return &IApi{
		Config:    config,
		TplEngine: template.NewTplEngine(w, r),
		Module:    "",
		Namespace: "Api",
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

func (i *IApi) GetString(name string) string {
	return utils.ToString(i.R.URL.Query().Get(name))
}

func (i *IApi) GetInt(name string) int {
	return utils.ToInt(i.R.URL.Query().Get(name))
}

func (i *IApi) PostString(name string) string {
	return utils.ToString(i.R.FormValue(name))
}

func (i *IApi) PostInt(name string) int {
	return utils.ToInt(i.R.FormValue(name))
}

func (i *IApi) Register(action string, f func()) *IApi {
	if i.Actions == nil {
		i.Actions = map[string]func(){}
	}
	if i.Module == "" {
		exception.CheckError(exception.NewError("error: api is empty!"), 999)
	}
	i.Actions[action] = f
	return i
}

func (i *IApi) Run(action string) {
	// 注册全局变量
	if i.TplEngine.TplData["GModule"] == nil || i.TplEngine.TplData["GModule"] == "" {
		i.TplEngine.TplData["GModule"] = i.Module
	}
	if i.TplEngine.TplData["GAction"] == nil || i.TplEngine.TplData["GAction"] == "" {
		i.TplEngine.TplData["GAction"] = action
	}
	// 检查action方法是否存在
	if f, ok := i.Actions[action]; !ok {
		if defaultFunc, ok1 := i.Actions["index"]; !ok1 {
			fmt.Fprintln(i.TplEngine.W, "404 page not found!")
			log.Println("404")
		} else {
			i.TplEngine.TplData["GAction"] = "index"
			defaultFunc()
		}
	} else {
		// run action
		f()
	}
}

func (c *IApi) Index() {
	fmt.Fprintln(c.TplEngine.W, "hello world, this is default index")
}

// base interface
type IBase struct {
	IApi
}

func NewIBase(config config.IConfig, w http.ResponseWriter, r *http.Request) *IBase {
	base := &IBase{
		IApi: *NewIApi(config, w, r),
	}
	base.Module = "base"
	return base
}

func (c *IApi) ResponseWithHeader(code int, result interface{}, message string) {
	c.TplEngine.ResponseWithHeader(code, result, message, c.Header)
}

func (c *IApi) Response(code int, result interface{}, message string) {
	c.TplEngine.Response(code, result, message)
}
