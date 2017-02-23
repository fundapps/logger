package logger

import (
	"os"

	"github.com/Sirupsen/logrus"
	logtry "github.com/evalphobia/logrus_sentry"
)

// Fields is for passing semi-structured data that doesn't already have a type to the logger
type Fields map[string]interface{}

type Fielder interface {
	ToFields() Fields
}

var (
	log    = logrus.New()
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
	})
	if err != nil {
		Fatal(err)
	}

	log.Hooks.Add(sentry)
}

func Info(message string, data Fields) {
	log.WithFields(logrus.Fields(data)).Info(message)
}

func Error(err error) {
	log.WithField("error", errorToField(err)).Error(err.Error())
}

func Fatal(err error) {
	// this doesn't use log.Fatal because we need to flush the hook before exiting
	log.WithField("error", errorToField(err)).Error(err.Error())
	Flush()
	os.Exit(1)
}

func Flush() {
	if sentry != nil {
		sentry.Flush()
	}
}

func errorToField(err error) interface{} {
	if err == nil {
		return nil
	}

	if fielder, ok := err.(Fielder); ok {
		return fielder.ToFields()
	}

	return err.Error()
}
