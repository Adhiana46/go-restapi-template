package logger

type Logger interface {
	Info(v ...interface{})
	Warn(v ...interface{})
	Error(v ...interface{})

	Infof(format string, v ...interface{})
	Warnf(format string, v ...interface{})
	Errorf(format string, v ...interface{})
}
