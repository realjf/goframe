package middleware

import (
	"net/http"
	"regexp"
	"github.com/op/go-logging"
	"kboard/utils"
	"os"
)

// Example format string. Everything except the message has a custom color
// which is dependent on the log level. Many fields have a custom output
// formatting too, eg. the time returns the hour down to the milli second.
var format = logging.MustStringFormatter(
	`%{color}%{time:15:04:05.000} %{shortfunc} ▶ %{level:.4s} %{id:04x}%{color:reset} %{message}`,
)

type Log struct {
	Logger *logging.Logger
}

func NewLogger() *Log {
	return &Log{
		Logger: logging.MustGetLogger("mylogger"),
	}
}

func (l *Log) Init() *Log {
	// For demo purposes, create two backend for os.Stderr.
	backend1 := logging.NewLogBackend(os.Stderr, "", 0)
	backend2 := logging.NewLogBackend(os.Stderr, "", 0)

	// For messages written to backend2 we want to add some additional
	// information to the output, including the used log level and the name of
	// the function.
	backend2Formatter := logging.NewBackendFormatter(backend2, format)

	// Only errors and more severe messages should be sent to backend1
	backend1Leveled := logging.AddModuleLevel(backend1)
	backend1Leveled.SetLevel(logging.ERROR, "")

	// Set the backends to be used.
	logging.SetBackend(backend1Leveled, backend2Formatter)

	return l
}

// log
func Logger(next http.Handler) http.Handler {
	Logger := NewLogger()
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if matched, _ := regexp.Match("^/assets/.*", []byte(r.RequestURI)); !matched {
			//log.Println(fmt.Sprintf("%s %s %s | %s", r.Method, r.RequestURI, r.Proto, utils.GetIPAdress(r)))
			Logger.Logger.Infof("%s %s %s | %s", r.Method, r.RequestURI, r.Proto, utils.GetIPAdress(r))
		}
		next.ServeHTTP(w, r)
	})
}
