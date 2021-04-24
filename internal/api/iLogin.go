package api

import (
	"github.com/realjf/goframe/config"
	"net/http"
)

type ILogin struct {
	IApi
}

func NewILogin(config config.IConfig, w http.ResponseWriter, r *http.Request) *ILogin {
	order := &ILogin{
		IApi: *NewIApi(config, w, r),
	}
	order.Module = "login"
	return order
}

func (this *ILogin) Index() {
	email := this.PostString("email")
	password := this.PostString("password")
	if email == "" || password == "" {
		this.ResponseWithHeader(101, "", "缺少数据")
	}
	result := map[string]string{
		"email":    email,
		"password": password,
		"username": email,
	}
	this.ResponseWithHeader(100, result, "登录成功")
}
