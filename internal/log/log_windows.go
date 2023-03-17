//go:build windows

package log

import (
	"golang.org/x/sys/windows/svc/eventlog"
)

var eventLogger *eventlog.Log

func setupSys() {
	var err error
	eventLogger, err = eventlog.Open("zunagi")
	if err != nil {
		panic(err)
	}
}

func logToSys(level logLevel, data string) {
	switch level {
	case INFO:
		eventLogger.Info(0, data)
	case WARN:
		eventLogger.Warning(0, data)
	case ERRO:
		eventLogger.Error(0, data)
	}
}
