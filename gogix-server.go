/* Copyright 2013 Juliano Martinez
   All Rights Reserved.

     Licensed under the Apache License, Version 2.0 (the "License");
     you may not use this file except in compliance with the License.
     You may obtain a copy of the License at

         http://www.apache.org/licenses/LICENSE-2.0

     Unless required by applicable law or agreed to in writing, software
     distributed under the License is distributed on an "AS IS" BASIS,
     WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
     See the License for the specific language governing permissions and
     limitations under the License.

   @author: Juliano Martinez */

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
	var conn broker.Connection

	config_file := os.Getenv("GOGIX_CONF")
	if strings.TrimSpace(config_file) == "" {
		config_file = "/etc/gogix/gogix.conf"
	}

	Cfg, err = goconfig.ReadConfigFile(config_file)
	utils.CheckPanic(err, "File not found")

	bind_addr, err := Cfg.GetString("server", "bind_addr")
	utils.CheckPanic(err, "Unable to get bind_addr from gogix.conf")
	conn.Queue, err = Cfg.GetString("transport", "queue")
	utils.CheckPanic(err, "Unable to get queue from gogix.conf")
	conn.Uri, err = Cfg.GetString("transport", "uri")
	utils.CheckPanic(err, "Unable to get transport from gogix.conf")
	conn.Expiration, err = Cfg.GetString("transport", "message_ttl")
	utils.CheckPanic(err, "Unable to get message_ttl from gogix.conf")
	addr, err := net.ResolveUDPAddr("udp", bind_addr)
	utils.CheckPanic(err, "Unable to resolve bind address")
	l, err := net.ListenUDP("udp", addr)
	utils.CheckPanic(err, fmt.Sprintf("Unable to bind %s", addr))

	if *debug == true {
		fmt.Printf("Setting-Up Broker %s\n", conn.Uri)
	}
	conn = conn.SetupBroker()
	go conn.NotifyClose()

	for {
		recv := make([]byte, 1024)
		_, remote_addr, err := l.ReadFromUDP(recv)
		utils.CheckPanic(err, "Problem receiving data")
		ip := fmt.Sprintf("%s", remote_addr.IP)

		go handle_data(string(recv), conn, ip)
	}

	defer conn.Close()
}

func handle_data(data string, conn broker.Connection, remote_addr string) {
	if *debug == true {
		fmt.Printf("Received log %s\n", data)
	}

	parsed := syslog.Graylog2ParseLog(data, remote_addr)

	if *debug == true {
		fmt.Printf("Sending data %s\n", parsed)
	}
	conn.Send(parsed)
}
