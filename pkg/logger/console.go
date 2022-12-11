package logger

import (
	"fmt"
	"log"
	"os"
)

type logger struct {
	infoLogger *log.Logger
	warnLogger *log.Logger
	errLogger  *log.Logger
}

var (
	colorReset = "\033[0m"
	colorInfo  = "\033[36m"
	colorWarn  = "\033[33m"
	colorErr   = "\033[31m"
)

func NewConsoleLogger() Logger {
	flags := log.LstdFlags | log.Lshortfile
	infoLogger := log.New(os.Stdout, fmt.Sprintf("%s[INFO]%s: ", colorInfo, colorReset), flags)
	warnLogger := log.New(os.Stdout, fmt.Sprintf("%s[WARN]%s: ", colorWarn, colorReset), flags)
	errLogger := log.New(os.Stdout, fmt.Sprintf("%s[ERROR]%s: ", colorErr, colorReset), flags)

	return &logger{
		infoLogger: infoLogger,
		warnLogger: warnLogger,
		errLogger:  errLogger,
	}
}

func (l *logger) Info(v ...interface{}) {
	l.infoLogger.Println(v...)
}

func (l *logger) Warn(v ...interface{}) {
	l.warnLogger.Println(v...)
}

func (l *logger) Error(v ...interface{}) {
	l.errLogger.Println(v...)
}

func (l *logger) Infof(format string, v ...interface{}) {
	l.infoLogger.Printf(format, v...)
}

func (l *logger) Warnf(format string, v ...interface{}) {
	l.warnLogger.Printf(format, v...)
}

func (l *logger) Errorf(format string, v ...interface{}) {
	l.errLogger.Printf(format, v...)
}
