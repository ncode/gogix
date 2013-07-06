package util

import (
	"fmt"
	"github.com/msbranco/goconfig"
	"log"
	"log/syslog"
)

var Cfg *config.Config

func CheckPanic(err error, message string) {
	if err != nil {
		msg := fmt.Sprintf("%s: %s", message, err)
		_log, err := syslog.New(syslog.LOG_ERR, "gogix")
		if err != nil {
			log.Fatalln("Unable to write syslog message")
		}
		_log.Warning(msg)
		defer _log.Close()
		log.Fatalln(msg)
		// panic(msg)
	}
}
