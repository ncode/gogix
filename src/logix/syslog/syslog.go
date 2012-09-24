package syslog

import (
    "time"
    "fmt"
    "regexp"
    "../../logix/util"
    "strconv"
)

type Parser struct {
    host string
    timestamp int64
    facility string
    level string
    version float32
    short_message string
}

var Severity = []string{"emerg", "alert", "crit", "err", "warn", "notice", "info", "debug"}
var Facility = []string{"kern", "user", "mail", "daemon", "auth", "syslog", "lpr",
                        "news", "uucp", "cron", "authpriv", "ftp", "ntp", "audit",
                        "alert", "at", "local0", "local1", "local2", "local3",
                        "local4", "local5", "local6", "local7"}

var LvlRegex = regexp.MustCompile("^<(.+?)>([A-Za-z]{3} .*)")

func ParseLog(line string) Parser {
    parsed := Parser{}
    now := time.Now()
    parsed.timestamp = now.Unix()
    parsed.version = 1.0
    lvl := LvlRegex.FindStringSubmatch(line)
    if len(lvl) >= 2 {
        i, err := strconv.Atoi(lvl[1])
        util.CheckPanic(err, fmt.Sprintf("Unable to convert %s to int", i))
        parsed.facility = Facility[i/8]
        parsed.level = Severity[i%8]
        parsed.short_message = lvl[2]
    }
    parsed.facility = "syslog"
    parsed.level = "info"
    parsed.short_message = line
    return parsed
}