package richerror

import (
	"errors"
	"fmt"
)

type ErrorLogger interface {
	Log(err error)
}

type Logger struct {
	BasicLogger     BasicLogger
	FormattedLogger FormattedLogger
	ContextLogger   ContextLogger
}

func (l Logger) Log(err error) {
	var rErr RichError
	ok := errors.As(err, &rErr)
	if !ok {
		l.logNormalError(err)
		return
	}

	l.logRichError(rErr)
}

type BasicLogger interface {
	Debug(...interface{})
	Info(...interface{})
	Warn(...interface{})
	Error(...interface{})
	Panic(...interface{})
	Fatal(...interface{})
}

type FormattedLogger interface {
	Debugf(string, ...interface{})
	Infof(string, ...interface{})
	Warnf(string, ...interface{})
	Errorf(string, ...interface{})
	Panicf(string, ...interface{})
	Fatalf(string, ...interface{})
}

type ContextLogger interface {
	Debugw(string, ...interface{})
	Infow(string, ...interface{})
	Warnw(string, ...interface{})
	Errorw(string, ...interface{})
	Panicw(string, ...interface{})
	Fatalw(string, ...interface{})
}

func (l Logger) logNormalError(err error) {
	if l.ContextLogger != nil {
		l.ContextLogger.Errorw(err.Error())
		return
	}

	if l.FormattedLogger != nil {
		l.FormattedLogger.Errorf(err.Error())
		return
	}

	if l.BasicLogger != nil {
		l.BasicLogger.Error(err.Error())
		return
	}

	fmt.Printf("error: %s\n", err.Error())
}

func (l Logger) logRichError(err RichError) {
	if l.ContextLogger != nil {
		contexts := []interface{}{
			"metadata", err.Metadata(),
			"runtime_info", err.RuntimeInfo(),
		}

		switch err.Level() {
		case Fatal:
			l.ContextLogger.Fatalw(err.Error(), contexts...)
		case Error:
			l.ContextLogger.Errorw(err.Error(), contexts...)
		case Warning:
			l.ContextLogger.Warnw(err.Error(), contexts...)
		case Info:
			l.ContextLogger.Infow(err.Error(), contexts...)
		default:
			l.ContextLogger.Errorw(err.Error(), contexts...)
		}

		return
	}

	if l.FormattedLogger != nil {
		switch err.Level() {
		case Fatal:
			l.FormattedLogger.Fatalf(err.Error())
		case Error:
			l.FormattedLogger.Errorf(err.Error())
		case Warning:
			l.FormattedLogger.Warnf(err.Error())
		case Info:
			l.FormattedLogger.Infof(err.Error())
		default:
			l.FormattedLogger.Errorf(err.Error())
		}

		return
	}

	if l.BasicLogger != nil {
		switch err.Level() {
		case Fatal:
			l.BasicLogger.Fatal(err.Error())
		case Error:
			l.BasicLogger.Error(err.Error())
		case Warning:
			l.BasicLogger.Warn(err.Error())
		case Info:
			l.BasicLogger.Info(err.Error())
		default:
			l.BasicLogger.Error(err.Error())
		}

		return
	}

	fmt.Printf("%s: %s", err.Kind(), err.Error())
}