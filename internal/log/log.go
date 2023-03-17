package log

import (
	"fmt"
	"log"
	"os"
)

type logLevel string

const (
	INFO logLevel = "[I]"
	WARN logLevel = "[W]"
	ERRO logLevel = "[E]"
)

func Setup() {
	setupSys()
	log.SetOutput(os.Stdout)
}

func logToStd(level logLevel, data string) {
	log.Println(level, data)
}

func logf(level logLevel, format string, args []any) {
	data := fmt.Sprintf(format, args...)
	logToSys(level, data)
	logToStd(level, data)
}

func Infof(format string, args ...any) {
	logf(INFO, format, args)
}

func Warnf(format string, args ...any) {
	logf(WARN, format, args)
}

func Errorf(format string, args ...any) {
	logf(ERRO, format, args)
}

func Fatalf(format string, args ...any) {
	logf(ERRO, format, args)
	os.Exit(1)
}
