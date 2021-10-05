package richerror

import (
	"errors"
	"fmt"
)

// Assert Logger implements ErrorLogger
var _ ErrorLogger = ChainLogger{}

// ChainLogger is an ErrorLogger that upon being called will call a set of loggers one after the other
type ChainLogger struct {
	Loggers []ErrorLogger
}

func (c ChainLogger) Log(err error) {
	for _, logger := range c.Loggers {
		logger.Log(err)
	}
}

func (c ChainLogger) LogInfo(msg string) {
	for _, logger := range c.Loggers {
		logger.LogInfo(msg)
	}
}

func (c ChainLogger) LogInfoWithMetadata(msg string, metadata ...interface{}) {
	for _, logger := range c.Loggers {
		logger.LogInfo(msg)
	}
}

func ChainLoggers(loggers ...ErrorLogger) ErrorLogger {
	return ChainLogger{Loggers: loggers}
}

// Assert Logger implements ErrorLogger
var _ ErrorLogger = Logger{}

// Logger is a struct that provides an ErrorLogger. The resulting ErrorLogger will log RichErrors given to it as
// descriptive as it can (based on the loggers abilities). Keep in mind that it's the module users' responsibility to
// give the struct their desired loggers. It will try to use ContextLogger which logs the error along with all of its
// Metadata. If not exists it will try FormattedLogger, BasicLogger, and GoLogger in that order. Finally, if no logger
// has been defined it will use fmt.Println to log the error.
type Logger struct {
	GoLogger        GoLogger
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

func (l Logger) LogInfo(msg string) {
	l.logRichError(New(msg).WithLevel(Info))
}

func (l Logger) LogInfoWithMetadata(msg string, metadata ...interface{}) {
	err := New(msg).WithLevel(Info)
	var key string
	for _, datum := range metadata {
		if key == "" {
			key = datum.(string)
			continue
		}

		err = err.WithField(key, datum)
		key = ""
	}
	l.logRichError(err)
}

type GoLogger interface {
	Printf(format string, v ...interface{})
	Println(v ...interface{})
	Fatalf(format string, v ...interface{})
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

	if l.GoLogger != nil {
		l.GoLogger.Println(err.Error())
		return
	}

	fmt.Printf("error: %s\n", err.Error())
}

func (l Logger) logRichError(err RichError) {
	if l.ContextLogger != nil {
		contexts := []interface{}{
			"metadata", err.Metadata(),
		}

		if err.Level() != Info {
			contexts = append(contexts, "runtime_info", err.RuntimeInfo())
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

	if l.GoLogger != nil {
		if err.Level() == Fatal {
			l.GoLogger.Fatalf("%s: %s\n", err.Kind(), err.Error())
			return
		}

		l.GoLogger.Printf("%s: %s\n", err.Kind(), err.Error())
		return
	}

	fmt.Printf("%s: %s\n", err.Kind(), err.Error())
}
