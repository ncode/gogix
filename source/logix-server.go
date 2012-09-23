package main

import (
    "os"
    "net"
    "fmt"
    "flag"
    "strings"
    "./logix/util"
    "./logix/syslog"
    "./logix/broker"
    "github.com/kless/goconfig/config"
)


var uri string
var queue string
var Cfg *config.Config
var user = flag.String("u", "logix", "username")
var debug = flag.Bool("d", false, "debug")

func main(){
    flag.Parse()
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
    if (*debug == true) {
        fmt.Printf("Received log %s\n", data)
        fmt.Printf("Connecting to Broker %s\n", uri)
    }
    conn = conn.Dial(uri)
    if (*debug == true) {
        fmt.Printf("Setup queue %s\n", queue)
    }
    conn = conn.SetupBroker(queue)
    if (*debug == true) {
        fmt.Printf("Sending data %s\n", parsed)
    }
    conn.Send(parsed)
}