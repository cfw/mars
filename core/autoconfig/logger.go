package autoconfig

import (
	"bytes"
	"fmt"
	log "github.com/sirupsen/logrus"
	"runtime"
	"strconv"
	"strings"
)

type logFormatter struct {
	timestampFormat string
}

func (s *logFormatter) Format(entry *log.Entry) ([]byte, error) {
	entry.Time.Format(s.timestampFormat)
	var fileLine string
	if entry.Caller != nil {
		file := trimmedPath(entry)
		line := entry.Caller.Line
		fileLine = fmt.Sprintf("%s:%d", file, line)
	}
	msg := fmt.Sprintf("%s [%s][%10d] %-30s: %s\n",
		entry.Time.Format(s.timestampFormat),
		strings.ToUpper(entry.Level.String()),
		getGID(),
		fileLine,
		entry.Message)
	return []byte(msg), nil
}

func getGID() uint64 {
	b := make([]byte, 64)
	b = b[:runtime.Stack(b, false)]
	b = bytes.TrimPrefix(b, []byte("goroutine "))
	b = b[:bytes.IndexByte(b, ' ')]
	n, _ := strconv.ParseUint(string(b), 10, 64)
	return n
}
func trimmedPath(ec *log.Entry) string {
	idx := strings.LastIndexByte(ec.Caller.File, '/')
	if idx == -1 {
		return ec.Caller.File
	}
	// Find the penultimate separator.
	idx = strings.LastIndexByte(ec.Caller.File[:idx], '/')
	if idx == -1 {
		return ec.Caller.File
	}
	buf := new(strings.Builder)
	// Keep everything after the penultimate separator.
	buf.WriteString(ec.Caller.File[idx+1:])
	caller := buf.String()
	return caller
}

func init() {
	log.SetFormatter(&logFormatter{timestampFormat: "2006-01-02 15:04:05.000"})
	log.SetReportCaller(true)
}
