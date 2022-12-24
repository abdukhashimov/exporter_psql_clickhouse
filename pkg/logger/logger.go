// Package logger represents a generic logging interface
package logger

// Logger represent common interface for logging function
type Logger interface {
	Debugf(format string, args ...interface{})
	Infof(format string, args ...interface{})
	Warnf(format string, args ...interface{})
	Errorf(format string, args ...interface{})
	Fatalf(format string, args ...interface{})

	Debug(args ...interface{})
	Info(args ...interface{})
	Warn(args ...interface{})
	Error(args ...interface{})
	Fatal(args ...interface{})

	Sync() error
}

// Log is a package level variable, every program should access logging function through "Log"
var (
	Log Logger
)

// SetLogger is the setter for log variable, it should be the only way to assign value to log
func SetLogger(newLogger Logger) {
	Log = newLogger
}
