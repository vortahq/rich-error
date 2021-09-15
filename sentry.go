package richerror

import (
	"errors"
	"time"

	"github.com/getsentry/sentry-go"
)

// SentryLogger is a ErrorLogger that logs to sentry. The returned ErrorLogger will report details of your errors (if
// they're RichError) to the sentry using `sentry-go` module.
type SentryLogger struct {
	environment string
	serverName  string
}

func (s SentryLogger) Log(err error) {
	sentryHub := sentry.CurrentHub().Clone()

	var rErr RichError
	ok := errors.As(err, &rErr)
	if !ok {
		sentryHub.CaptureException(err)
		return
	}

	event := sentry.NewEvent()

	event.Contexts = rErr.Metadata()
	event.Environment = s.environment
	event.Level = rErr.Level().SentryLevel()
	event.Message = rErr.Error()
	event.ServerName = s.serverName
	event.Timestamp = time.Now()

	event.Tags["kind"] = rErr.Kind().String()
	if rErr.Operation() != "" {
		event.Tags["operation"] = string(rErr.Operation())
	}

	sentryHub.CaptureEvent(event)
}
