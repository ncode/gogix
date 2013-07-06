package main

import (
	"github.com/ncode/gogix/broker"
	"github.com/ncode/gogix/syslog"
	"github.com/ncode/gogix/util"
	"flag"
	"fmt"
	"github.com/msbranco/goconfig"
	"net"
	"os"
	"strings"
)

var uri string
var queue string
var Cfg *goconfig.ConfigFile
var user = flag.String("u", "gogix", "username")
var debug = flag.Bool("d", false, "debug")

func main() {
	flag.Parse()
	var err error
	config_file := os.Getenv("GOGIX_CONF")
	if strings.TrimSpace(config_file) == "" {
		config_file = "/etc/gogix/gogix.conf"
	}

	Cfg, err = goconfig.ReadConfigFile(config_file)
	util.CheckPanic(err, "File not found")

	bind_addr, err := Cfg.GetString("server", "bind_addr")
	util.CheckPanic(err, "Unable to get bind_addr from gogix.conf")
	queue, err = Cfg.GetString("transport", "queue")
	util.CheckPanic(err, "Unable to get queue from gogix.conf")
	uri, err = Cfg.GetString("transport", "uri")
	util.CheckPanic(err, "Unable to get transport from gogix.conf")
	message_ttl, err := Cfg.GetInt64("transport", "message_ttl")
	util.CheckPanic(err, "Unable to get message_ttl from gogix.conf")

	addr, err := net.ResolveUDPAddr("udp", bind_addr)
	util.CheckPanic(err, "Unable to resolve bind address")

	l, err := net.ListenUDP("udp", addr)
	util.CheckPanic(err, fmt.Sprintf("Unable to bind %s", addr))

	for {
		recv := make([]byte, 1024)
		_, _, err := l.ReadFromUDP(recv)
		util.CheckPanic(err, "Problem receiving data")
		go handle_data(string(recv), message_ttl)
	}
}

func handle_data(data string, message_ttl int64) {
	parsed := syslog.ParseLog(data)
	var conn broker.Connection
	if *debug == true {
		fmt.Printf("Received log %s\n", data)
		fmt.Printf("Connecting to Broker %s\n", uri)
	}
	conn = conn.Dial(uri)
	if *debug == true {
		fmt.Printf("Setup queue %s\n", queue)
	}
	conn = conn.SetupBroker(queue, message_ttl)
	if *debug == true {
		fmt.Printf("Sending data %s\n", parsed)
	}
	conn.Send(parsed)
}
