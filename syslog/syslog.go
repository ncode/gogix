package syslog

import (
	"fmt"
	"github.com/ncode/gogix/utils"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type Graylog2Parsed struct {
	Host         string `json:"host"`
	Timestamp    int64  `json:"timestamp"`
	Facility     string `json:"facility"`
	Level        int    `json:"level"`
	Version      string `json:"version"`
	ShortMessage string `json:"short_message"`
}

var Severity = []string{"emerg", "alert", "crit", "err", "warn", "notice", "info", "debug"}
var Facility = []string{"kern", "user", "mail", "daemon", "auth", "syslog", "lpr",
	"news", "uucp", "cron", "authpriv", "ftp", "ntp", "audit",
	"alert", "at", "local0", "local1", "local2", "local3",
	"local4", "local5", "local6", "local7"}

var LvlRegex = regexp.MustCompile("^<(.+?)>([A-Za-z]{3} .*)")

func Graylog2ParseLog(line string, remote_addr string) Graylog2Parsed {
	parsed := Graylog2Parsed{}
	now := time.Now()
	parsed.Timestamp = now.Unix()
	parsed.Version = "1.0"
	if strings.Contains(remote_addr, "127.0.0.1") {
		hostname, err := os.Hostname()
		utils.CheckPanic(err, fmt.Sprintf("Unable to get my hostname"))
		parsed.Host = hostname
	} else {
		parsed.Host = remote_addr
	}

	lvl := LvlRegex.FindStringSubmatch(line)
	if len(lvl) >= 2 {
		i, err := strconv.Atoi(lvl[1])
		utils.Check(err, fmt.Sprintf("Unable to convert %s to int", i))
		parsed.Facility = strings.Title(Facility[i/8])
		parsed.Level = i % 8
		parsed.ShortMessage = strings.TrimRight(lvl[2], "\u0000")
	} else {
		parsed.Facility = "Syslog"
		parsed.Level = 6
		parsed.ShortMessage = strings.TrimRight(line, "\u0000")
	}
	return parsed
}
