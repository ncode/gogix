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

type Parser struct {
	Host         string  `json:"host"`
	Timestamp    int64   `json:"timestamp"`
	Facility     string  `json:"facility"`
	Level        string  `json:"level"`
	Version      float64 `json:"version"`
	ShortMessage string  `json:"short_message"`
}

var Severity = []string{"emerg", "alert", "crit", "err", "warn", "notice", "info", "debug"}
var Facility = []string{"kern", "user", "mail", "daemon", "auth", "syslog", "lpr",
	"news", "uucp", "cron", "authpriv", "ftp", "ntp", "audit",
	"alert", "at", "local0", "local1", "local2", "local3",
	"local4", "local5", "local6", "local7"}

var LvlRegex = regexp.MustCompile("^<(.+?)>([A-Za-z]{3} .*)")

func ParseLog(line string, remote_addr string) Parser {
	parsed := Parser{}
	now := time.Now()
	parsed.Timestamp = now.Unix()
	parsed.Version = 1.0
	hostname, err := os.Hostname()
	parsed.Host = hostname
	utils.Check(err, fmt.Sprintf("Unable to get my hostname"))
	lvl := LvlRegex.FindStringSubmatch(line)
	if len(lvl) >= 2 {
		i, err := strconv.Atoi(lvl[1])
		utils.Check(err, fmt.Sprintf("Unable to convert %s to int", i))
		parsed.Facility = Facility[i/8]
		parsed.Level = Severity[i%8]
		parsed.ShortMessage = strings.TrimRight(lvl[2], "\u0000")
	} else {
		parsed.Facility = "syslog"
		parsed.Level = "info"
		parsed.ShortMessage = strings.TrimRight(line, "\u0000")
	}
	return parsed
}
