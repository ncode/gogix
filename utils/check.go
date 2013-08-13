package utils

import (
	"fmt"
	"log"
	"log/syslog"
)

var (
	_log, s_err = syslog.New(syslog.LOG_ERR, "gogix")
)

func Check(err error, message string) {
	check(err, message, false)
}

func CheckPanic(err error, message string) {
	check(err, message, true)
}

func Log(message string) {
	CheckPanic(s_err, "Unable to write syslog message")
	_log.Info(message)
	defer _log.Close()
}

func check(err error, message string, _panic bool) {
	if err != nil {
		msg := fmt.Sprintf("%s: %s", message, err)
		if s_err != nil {
			log.Fatalln("Unable to write syslog message")
		}
		_log.Warning(msg)
		defer _log.Close()
		log.Fatalln(msg)
		if _panic {
			panic(msg)
		}
	}
}
