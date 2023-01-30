// SDR - Labo 4
// Nicolas Crausaz & Maxime Scharwath
// Defines various log functions

package log

import (
	"fmt"
	"time"
)

type LogLevel struct {
	color string
	name  string
	level int
}

// create enum for log levels
var (
	Trace = LogLevel{"\033[0;36m", "TRACE", 0}
	Debug = LogLevel{"\033[0;34m", "DEBUG", 1}
	Info  = LogLevel{"\033[0;32m", "INFO", 2}
	Warn  = LogLevel{"\033[0;33m", "WARN", 3}
	Error = LogLevel{"\033[0;31m", "ERROR", 4}
	Fatal = LogLevel{"\033[0;35m", "FATAL", 5}
)

var logLevel = 0
var logEnabled = true

func SetLogLevel(level LogLevel) {
	logLevel = level.level
}

func SetLogLevelByValue(level int) {
	logLevel = level
}

func SetLogEnabled(enabled bool) {
	logEnabled = enabled
}

func Log(level LogLevel, message string) {
	if logEnabled && level.level >= logLevel {
		timestamp := time.Now().Format("15:04:05")
		colorReset := "\033[0m"
		message = fmt.Sprintf("%s%s [%s]%s: %s", level.color, timestamp, level.name, colorReset, message)
		if level.level == Fatal.level {
			panic(message)
		} else {
			fmt.Println(message)
		}
	}
}

func Logf(level LogLevel, message string, args ...interface{}) {
	Log(level, fmt.Sprintf(message, args...))
}
