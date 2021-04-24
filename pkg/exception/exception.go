package exception

import (
	"errors"
	"fmt"
	"log"
	"os"
)

func NewError(format string, args ...string) error {
	return errors.New(fmt.Sprintf(format, args))
}

func CheckError(err error, code int) {
	if err != nil {
		log.Println(err.Error())
		if code == 0 || code == -1 {
			os.Exit(code)
		}
	}
}
