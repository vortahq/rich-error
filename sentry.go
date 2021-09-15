package richerror

import (
	"errors"
	"time"

	"github.com/getsentry/sentry-go"
)

type SentryLogger func(err error)

// SentryLoggerFactory is a Factory that creates a SentryLogger based on the options given to it.
// The returned SentryLogger will report details of your errors (if they're RichError) to the sentry using `sentry-go`
// module.
func SentryLoggerFactory(options ...SentryLoggerOptions) SentryLogger {
	option := aggregateOptions(options...)

	return func(err error) {
		sentryHub := sentry.CurrentHub().Clone()

		var rErr RichError
		ok := errors.As(err, &rErr)
		if !ok {
			sentryHub.CaptureException(err)
			return
		}

		event := sentry.NewEvent()

		event.Contexts = rErr.Metadata()
		event.Environment = option.environment
		event.Level = rErr.Level().SentryLevel()
		event.Message = rErr.Error()
		event.ServerName = option.serverName
		event.Timestamp = time.Now()

		event.Tags["kind"] = rErr.Kind().String()
		if rErr.Operation() != "" {
			event.Tags["operation"] = string(rErr.Operation())
		}

		sentryHub.CaptureEvent(event)
	}
}

type SentryLoggerOptions struct {
	environment string
	serverName  string
}

func EnvironmentOption(environment string) SentryLoggerOptions {
	return SentryLoggerOptions{environment: environment}
}

func ServerNameOption(serverName string) SentryLoggerOptions {
	return SentryLoggerOptions{serverName: serverName}
}

func aggregateOptions(options ...SentryLoggerOptions) SentryLoggerOptions {
	aggregate := SentryLoggerOptions{}

	for _, option := range options {
		if option.serverName != "" {
			aggregate.serverName = option.serverName
		}

		if option.environment != "" {
			aggregate.environment = option.environment
		}
	}

	return aggregate
}
