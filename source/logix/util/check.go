package util

import (
    "fmt"
    "log"
    "github.com/kless/goconfig/config"
)

var Cfg *config.Config

func CheckPanic(err error, message string) {
  if err != nil {
    msg := fmt.Sprintf("%s: %s", message, err)
    log.Fatalln(msg)
    panic(fmt.Sprintf(msg))
  }
}