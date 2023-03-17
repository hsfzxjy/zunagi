//go:build unix

package log

import "log/syslog"

var writer *syslog.Writer

func setupSys() {
	var err error
	writer, err = syslog.New(syslog.LOG_USER, "vbic")
	if err != nil {
		panic(err)
	}
}

func logToSys(level logLevel, data string) {
	switch level {
	case INFO:
		writer.Info(data)
	case WARN:
		writer.Warning(data)
	case ERRO:
		writer.Err(data)
	}
}
