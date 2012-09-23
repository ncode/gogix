package util

import (
    "fmt"
    "log"
    "log/syslog"
    "github.com/kless/goconfig/config"
)

var Cfg *config.Config

func CheckPanic(err error, message string) {
  if err != nil {
    msg := fmt.Sprintf("%s: %s", message, err)
    _log, err := syslog.New(syslog.LOG_ERR, "[go]")
    if err != nil {
        log.Fatalln("Unable to write syslog message")
    }
    defer _log.Close()
    log.Fatalln(msg)
    panic(fmt.Sprintf(msg))
  }
}