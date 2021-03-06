package logger

import (
	"fmt"
	"github.com/getsentry/raven-go"
	"os"

	logtry "github.com/evalphobia/logrus_sentry"
	"github.com/sirupsen/logrus"
)

// Fields is for passing semi-structured data that doesn't already have a type to the logger
type Fields map[string]interface{}

// Fielder is a interface that allows any type to be converted into log fields
type Fielder interface {
	ToFields() Fields
}

var (
	log          = logrus.New()
	globalFields = logrus.Fields{}

	sentry *logtry.SentryHook
)

func init() {
	log.Formatter = new(logrus.JSONFormatter)
	log.Out = os.Stderr

	// don't configure sentry in development
	if os.Getenv("APP_ENV") == "development" {
		return
	}

	sentryDSN := os.Getenv("SENTRY_DSN")

	sentry, err := logtry.NewAsyncSentryHook(sentryDSN, []logrus.Level{
		logrus.PanicLevel,
		logrus.FatalLevel,
		logrus.ErrorLevel,
		logrus.WarnLevel,
	})
	if err != nil {
		Fatal(err)
	}

	log.Hooks.Add(sentry)
}

func SetGlobalFields(data Fields) {
	var tags raven.Tags

	for key, value := range data {
		var stringValue string
		if str, ok := value.(string); ok {
			stringValue = str
		} else {
			stringValue = fmt.Sprintf("%#v", value)
		}

		tags = append(tags, raven.Tag{
			Key:   key,
			Value: stringValue,
		})
	}

	globalFields = logrus.Fields{"tags":tags}
}

// Info logs with the given message & addition fields at the INFO Level
// This doesn't log to sentry
func Info(message string, data Fields) {
	log.WithFields(globalFields).WithFields(logrus.Fields(data)).Info(message)
}

// Warn logs with the given message & addition fields at the Warn Level
func Warn(message string, data Fields) {
	log.WithFields(globalFields).WithFields(logrus.Fields(data)).Warn(message)
}

// WarnError logs an error & fields at the Warn Level
func WarnError(err error) {
	log.WithFields(globalFields).WithFields(errorToFields(err)).Warn(err.Error())
}

// Error logs an error with the error message & converts the message to fields if Fielder is implemented
func Error(err error) {
	log.WithFields(globalFields).WithFields(errorToFields(err)).Error(err.Error())
}

// Fatal logs errors in the same way as Error then flushes the errors and calls os.Exit(1)
func Fatal(err error) {
	// this doesn't use log.Fatal because we need to flush the hook before exiting
	log.WithFields(globalFields).WithFields(errorToFields(err)).Error(err.Error())
	Flush()
	os.Exit(1)
}

// Flush blocks until all pending messages have been processed
func Flush() {
	if sentry != nil {
		sentry.Flush()
	}
}

func errorToFields(err error) logrus.Fields {
	if err == nil {
		return nil
	}

	if fielder, ok := err.(Fielder); ok {
		return logrus.Fields(fielder.ToFields())
	}

	return logrus.Fields{
		"error": err.Error(),
		"data":  err,
	}
}
