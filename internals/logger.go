package internals

import (
	"os"

	"github.com/op/go-logging"
)

// Logger is a server logger
type Logger interface {
	Fatal(args ...interface{})
	Fatalf(format string, args ...interface{})
	Panic(args ...interface{})
	Panicf(format string, args ...interface{})
	Critical(args ...interface{})
	Criticalf(format string, args ...interface{})
	Error(args ...interface{})
	Errorf(format string, args ...interface{})
	Warning(args ...interface{})
	Warningf(format string, args ...interface{})
	Notice(args ...interface{})
	Noticef(format string, args ...interface{})
	Info(args ...interface{})
	Infof(format string, args ...interface{})
	Debug(args ...interface{})
	Debugf(format string, args ...interface{})
}

// LoggerID ia a logger ID
const LoggerID = "yc-backend"

var (
	// Logger settings
	logger           = logging.MustGetLogger(LoggerID)
	logConsoleFormat = logging.MustStringFormatter(
		`%{color}%{time:2006/01/02 15:04:05} [YC-Backend] >> %{message} %{color:reset}`,
	)
)

func init() {
	// Prepare logger
	logConsoleBackend := logging.NewLogBackend(os.Stderr, "", 0)
	logConsolePrettyBackend := logging.NewBackendFormatter(logConsoleBackend, logConsoleFormat)
	lvl := logging.AddModuleLevel(logConsolePrettyBackend)
	logging.SetBackend(lvl)
	//logging.SetLevel(logging.INFO, LoggerID)
	// Set log level based on env
	switch os.Getenv("APP_LOG_LEVEL") {
	case "debug":
		logging.SetLevel(logging.DEBUG, LoggerID)
	case "info":
		logging.SetLevel(logging.INFO, LoggerID)
	case "warn":
		logging.SetLevel(logging.WARNING, LoggerID)
	case "err":
		logging.SetLevel(logging.ERROR, LoggerID)
	default:
		logging.SetLevel(logging.INFO, LoggerID) // log everything by default
	}
}

// GetLogger returns the logger
func GetLogger() Logger {
	return logger
}
