package middleware

import (
	"net/http"
	"runtime/debug"
)

// log
func SafeHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if e, ok := recover().(error); ok {
				http.Error(w, e.Error(), http.StatusInternalServerError)
				// 输出自定义页面
				Logger.Logger.Warningf("[Warning]: panic in %+v - %+v", r.RequestURI, e)
				Logger.Logger.Info(string(debug.Stack()))
			}
		}()
		next.ServeHTTP(w, r)
	})
}
