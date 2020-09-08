package template

import (
	"encoding/json"
	"github.com/realjf/goframe/exception"
	"html/template"
	"net/http"
)

type ITplEngine interface {
	NewTemplate(string) *TplEngine
	Parse(string) *TplEngine
	ParseFiles(...string) *TplEngine
	AssignMap(map[string]interface{}) *TplEngine
	Assign(string, interface{}) *TplEngine
	Display(string)
	Response(int, interface{}, string)
	ResponseWithHeader(int, interface{}, string, map[string]string)
}

type ResponseData struct {
	Code    int         `json:"code"`
	Result  interface{} `json:"result"`
	Message string      `json:"message"`
}

type TplEngine struct {
	tpl     *template.Template
	TplData map[string]interface{}
	W       http.ResponseWriter
	R       *http.Request
	Header  map[string]string
}

func NewTplEngine(w http.ResponseWriter, r *http.Request) *TplEngine {
	return &TplEngine{
		W:       w,
		R:       r,
		TplData: make(map[string]interface{}),
		Header:  map[string]string{},
	}
}

func (t *TplEngine) NewTemplate(name string) *TplEngine {
	t.tpl = template.New(name)
	return t
}

func (t *TplEngine) Parse(content string) *TplEngine {
	var err error
	t.tpl, err = t.tpl.Parse(content)
	exception.CheckError(err, 2003)
	return t
}

func getDefinedTpl(file string) []string {
	file = "web/" + file + ".html"
	tpls := []string{
		file,
	}

	return tpls
}

func (t *TplEngine) ParseFiles(tplName ...string) *TplEngine {
	commFiles := getDefinedTpl(tplName[0])
	for _, f := range tplName {
		ff := "web/" + f + ".html"
		commFiles = append(commFiles, ff)
	}
	var err error
	t.tpl, err = t.tpl.ParseFiles(commFiles...)
	exception.CheckError(err, 2001)
	return t
}

func (t *TplEngine) AssignMap(data map[string]interface{}) *TplEngine {
	for k, v := range data {
		t.TplData[k] = v
	}
	return t
}

func (t *TplEngine) Assign(key string, data interface{}) *TplEngine {
	t.TplData[key] = data
	return t
}

func (t *TplEngine) Display(tpl string) {
	err := t.ParseFiles(tpl).tpl.Execute(t.W, t.TplData)
	exception.CheckError(err, 2004)
}

func (t *TplEngine) DisplayMulti(tpl string, subTpl []string) {
	tpls := []string{tpl}
	tpls = append(tpls, subTpl...)
	err := t.ParseFiles(tpls...).tpl.Execute(t.W, t.TplData)
	exception.CheckError(err, 2004)
}

func (t *TplEngine) Response(code int, result interface{}, message string) {
	data := ResponseData{
		Code:    code,
		Result:  result,
		Message: message,
	}
	jsonData, err := json.Marshal(data)
	exception.CheckError(err, 2005)
	t.W.Write(jsonData)
}

func (t *TplEngine) ResponseWithHeader(code int, result interface{}, message string, headerOptions map[string]string) {
	data := ResponseData{
		Code:    code,
		Result:  result,
		Message: message,
	}
	jsonData, err := json.Marshal(data)
	exception.CheckError(err, 2005)
	if len(headerOptions) > 0 {
		for field, val := range headerOptions {
			t.W.Header().Set(field, val)
		}
	}
	t.W.Write(jsonData)
}
