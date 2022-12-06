package monitoring

import (
	"log"
	"time"

	"github.com/getsentry/sentry-go"
)

type LogData struct {
	Method        string
	Endpoint      string
	RequestBody   interface{}
	RequestParams interface{}
	RequestQuery  interface{}
	StackTrace    interface{}
}

type SentryMonitor interface {
	Close()
	LogInfo(message string, data LogData) *sentry.EventID
	LogError(message string, data LogData) *sentry.EventID
	LogPanic(message string, data LogData) *sentry.EventID
}

type sentryMonitor struct {
}

func NewSentryMonitor(dsn string) SentryMonitor {
	err := sentry.Init(sentry.ClientOptions{
		Dsn: dsn,
		// Set TracesSampleRate to 1.0 to capture 100%
		// of transactions for performance monitoring.
		// We recommend adjusting this value in production,
		TracesSampleRate: 1.0,
	})
	if err != nil {
		log.Fatalf("sentry.Init: %s", err)
	}

	return &sentryMonitor{}
}

func (s *sentryMonitor) Close() {
	sentry.Flush(2 * time.Second)
}

func (s *sentryMonitor) dataToMap(data LogData) map[string]interface{} {
	return map[string]interface{}{
		"Method":         data.Method,
		"Endpoint":       data.Endpoint,
		"Request Body":   data.RequestBody,
		"Request Params": data.RequestParams,
		"Request Query":  data.RequestQuery,
		"Stack Trace":    data.StackTrace,
	}
}

func (s *sentryMonitor) LogInfo(message string, data LogData) *sentry.EventID {
	localHub := sentry.CurrentHub().Clone()
	localHub.ConfigureScope(func(scope *sentry.Scope) {
		scope.SetLevel(sentry.LevelInfo)
		scope.SetContext("info", s.dataToMap(data))
	})

	return localHub.CaptureMessage(message)
}

func (s *sentryMonitor) LogError(message string, data LogData) *sentry.EventID {
	localHub := sentry.CurrentHub().Clone()
	localHub.ConfigureScope(func(scope *sentry.Scope) {
		scope.SetLevel(sentry.LevelError)
		scope.SetContext("info", s.dataToMap(data))
	})

	return localHub.CaptureMessage(message)
}

func (s *sentryMonitor) LogPanic(message string, data LogData) *sentry.EventID {
	localHub := sentry.CurrentHub().Clone()
	localHub.ConfigureScope(func(scope *sentry.Scope) {
		scope.SetLevel(sentry.LevelFatal)
		scope.SetContext("info", s.dataToMap(data))
	})

	return localHub.CaptureMessage(message)
}
