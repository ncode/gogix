package main

import (
    "net"
    "./logix/util"
    "./logix/syslog"
    "./logix/broker"
)

func main(){
    addr, err := net.ResolveUDPAddr("up4", ":6660")
    util.Checkp(err)

    l, err := net.ListenUDP("udp", addr)
    util.Checkp(err)

    for {
        recv := make([]byte, 1024)
         _, _, err := l.ReadFromUDP(recv)
        util.Checkp(err)
        go handle_data(string(recv))
    }
}

func handle_data(data string){
    parsed := syslog.ParseLog(data)
    var conn broker.Connection
    conn = conn.Dial("amqp://guest:guest@karoly:5672/")
    conn = conn.SetupBroker("logix")
    conn.Send(parsed)
}