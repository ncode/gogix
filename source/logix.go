package main

import (
    "os"
    "net"
    "fmt"
    "strings"
    "./logix/util"
    "./logix/syslog"
    "./logix/broker"
    "github.com/kless/goconfig/config"
)

/*
[transport]
url = amqp://127.0.0.1:5672
queue = logix

[server]
bind_addr = 127.0.0.1:6660
 */

var Cfg *config.Config
var queue string
var uri string

func main(){
    var err error
    config_file := os.Getenv("LOGIX_CONF")
    if (strings.TrimSpace(config_file) == ""){
        config_file = "/etc/logix/logix.conf"
    }

    Cfg, err = config.ReadDefault(config_file)
    util.CheckPanic(err, "File not found")

    bind_addr, err := Cfg.String("server", "bind_addr")
    util.CheckPanic(err, "Unable to get bind_addr from logix.conf")
    queue, err = Cfg.String("transport", "queue")
    util.CheckPanic(err, "Unable to get queue from logix.conf")
    uri, err = Cfg.String("transport", "uri")
    util.CheckPanic(err, "Unable to get transport from logix.conf")

    addr, err := net.ResolveUDPAddr("up4", bind_addr)
    util.CheckPanic(err, "Unable to resolve bind address")

    l, err := net.ListenUDP("udp", addr)
    util.CheckPanic(err, fmt.Sprintf("Unable to bind %s", addr))

    for {
        recv := make([]byte, 1024)
         _, _, err := l.ReadFromUDP(recv)
        util.CheckPanic(err, "Problem receiving data")
        go handle_data(string(recv))
    }
}

func handle_data(data string){
    parsed := syslog.ParseLog(data)
    var conn broker.Connection
    conn = conn.Dial(uri)
    conn = conn.SetupBroker(queue)
    conn.Send(parsed)
}