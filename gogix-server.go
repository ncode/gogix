package main

import (
	"flag"
	"fmt"
	"github.com/msbranco/goconfig"
	"github.com/ncode/gogix/broker"
	"github.com/ncode/gogix/syslog"
	"github.com/ncode/gogix/utils"
	"net"
	"os"
	"strings"
)

var (
	Cfg   *goconfig.ConfigFile
	user  = flag.String("u", "gogix", "username")
	debug = flag.Bool("d", false, "debug")
)

func main() {
	flag.Parse()
	var err error

	config_file := os.Getenv("GOGIX_CONF")
	if strings.TrimSpace(config_file) == "" {
		config_file = "/etc/gogix/gogix.conf"
	}

	Cfg, err = goconfig.ReadConfigFile(config_file)
	utils.CheckPanic(err, "File not found")

	bind_addr, err := Cfg.GetString("server", "bind_addr")
	utils.CheckPanic(err, "Unable to get bind_addr from gogix.conf")
	queue, err := Cfg.GetString("transport", "queue")
	utils.CheckPanic(err, "Unable to get queue from gogix.conf")
	uri, err := Cfg.GetString("transport", "uri")
	utils.CheckPanic(err, "Unable to get transport from gogix.conf")
	message_ttl, err := Cfg.GetString("transport", "message_ttl")
	utils.CheckPanic(err, "Unable to get message_ttl from gogix.conf")
	addr, err := net.ResolveUDPAddr("udp", bind_addr)
	utils.CheckPanic(err, "Unable to resolve bind address")
	l, err := net.ListenUDP("udp", addr)
	utils.CheckPanic(err, fmt.Sprintf("Unable to bind %s", addr))

	var conn broker.Connection
	if *debug == true {
		fmt.Printf("Connecting to Broker %s\n", uri)
	}
	conn = conn.Dial(uri)

	if *debug == true {
		fmt.Printf("Setting-up queue %s\n", queue)
	}
	conn = conn.SetupBroker(queue, message_ttl)

	for {
		recv := make([]byte, 1024)
		_, remote_addr, err := l.ReadFromUDP(recv)
		utils.Check(err, "Problem receiving data")
		go handle_data(string(recv), message_ttl, conn, remote_addr.IP)
	}

	defer conn.Close()
}

func handle_data(data string, message_ttl string, conn broker.Connection, remote_addr string) {
	if *debug == true {
		fmt.Printf("Received log %s\n", data)
	}

	parsed := syslog.ParseLog(data, remote_addr)

	if *debug == true {
		fmt.Printf("Sending data %s\n", parsed)
	}
	conn.Send(parsed)
}
