package helpers

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/getsentry/sentry-go"
	richerror "gitlab.com/orderhq/rich-error"
)

type basicLogger interface {
	Debug(...interface{})
	Info(...interface{})
	Warn(...interface{})
	Error(...interface{})
	Panic(...interface{})
	Fatal(...interface{})
}

type formattedLogger interface {
	Debugf(string, ...interface{})
	Infof(string, ...interface{})
	Warnf(string, ...interface{})
	Errorf(string, ...interface{})
	Panicf(string, ...interface{})
	Fatalf(string, ...interface{})
}

type contextLogger interface {
	Debugw(string, ...interface{})
	Infow(string, ...interface{})
	Warnw(string, ...interface{})
	Errorw(string, ...interface{})
	Panicw(string, ...interface{})
	Fatalw(string, ...interface{})
}

func (h helper) log(path string, err error) {
	if strings.HasPrefix(path, "/grpc.reflection.v1alpha.ServerReflection/") {
		return
	}

	h.logToSentry(path, err)

	var rErr richerror.RichError
	ok := errors.As(err, &rErr)
	if !ok {
		h.logNormalError(path, err)
		return
	}

	h.logRichError(path, rErr)
}

func (h helper) logToSentry(path string, err error) {
	if !h.sentryEnabled {
		return
	}

	sentryHub := sentry.CurrentHub().Clone()

	var rErr richerror.RichError
	ok := errors.As(err, &rErr)
	if !ok {
		sentryHub.CaptureException(err)
		return
	}

	event := sentry.NewEvent()

	event.Contexts = rErr.Metadata()
	event.Environment = h.environment
	event.Level = rErr.Level().SentryLevel()
	event.Message = rErr.Error()
	event.ServerName = h.serverName
	event.Timestamp = time.Now()

	event.Tags["path"] = path
	event.Tags["kind"] = rErr.Kind().String()
	if rErr.Operation() != "" {
		event.Tags["operation"] = string(rErr.Operation())
	}

	sentryHub.CaptureEvent(event)
}

func (h helper) logNormalError(path string, err error) {
	if h.contextLogger != nil {
		h.contextLogger.Errorw(err.Error(), "path", path)
		return
	}

	if h.formattedLogger != nil {
		h.formattedLogger.Errorf("in %s: %s", path, err.Error())
		return
	}

	if h.basicLogger != nil {
		message := fmt.Sprintf("in %s: %s", path, err.Error())
		h.basicLogger.Error(message)
		return
	}

	fmt.Printf("error in %s: %s", path, err.Error())
}

func (h helper) logRichError(path string, err richerror.RichError) {
	if h.contextLogger != nil {
		contexts := []interface{}{
			"path", path,
			"metadata", err.Metadata(),
			"runtime_info", err.RuntimeInfo(),
		}

		switch err.Level() {
		case richerror.Fatal:
			h.contextLogger.Fatalw(err.Error(), contexts...)
		case richerror.Error:
			h.contextLogger.Errorw(err.Error(), contexts...)
		case richerror.Warning:
			h.contextLogger.Warnw(err.Error(), contexts...)
		case richerror.Info:
			h.contextLogger.Infow(err.Error(), contexts...)
		default:
			h.contextLogger.Errorw(err.Error(), contexts...)
		}

		return
	}

	if h.formattedLogger != nil {
		switch err.Level() {
		case richerror.Fatal:
			h.formattedLogger.Fatalf("in %s: %s", path, err.Error())
		case richerror.Error:
			h.formattedLogger.Errorf("in %s: %s", path, err.Error())
		case richerror.Warning:
			h.formattedLogger.Warnf("in %s: %s", path, err.Error())
		case richerror.Info:
			h.formattedLogger.Infof("in %s: %s", path, err.Error())
		default:
			h.formattedLogger.Errorf("in %s: %s", path, err.Error())
		}

		return
	}

	if h.basicLogger != nil {
		message := fmt.Sprintf("in %s: %s", path, err.Error())

		switch err.Level() {
		case richerror.Fatal:
			h.basicLogger.Fatal(message)
		case richerror.Error:
			h.basicLogger.Error(message)
		case richerror.Warning:
			h.basicLogger.Warn(message)
		case richerror.Info:
			h.basicLogger.Info(message)
		default:
			h.basicLogger.Error(message)
		}

		return
	}

	fmt.Printf("%s in %s: %s", err.Kind(), path, err.Error())
}
