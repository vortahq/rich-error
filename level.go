package richerror

import "github.com/getsentry/sentry-go"

// Level identifies severity of the error
type Level level
type level uint8

const (
	UnknownLevel Level = iota
	Fatal
	Error
	Warning
	Info
)

var levelStrings = [...]string{"_", "Fatal", "Error", "Warning", "Info"}

func (l Level) String() string {
	return levelStrings[l]
}

func (l Level) MarshalJSON() ([]byte, error) {
	return []byte(l.String()), nil
}

func (l Level) SentryLevel() sentry.Level {
	switch l {
	case Fatal:
		return sentry.LevelFatal
	case Error:
		return sentry.LevelError
	case Warning:
		return sentry.LevelWarning
	case Info:
		return sentry.LevelInfo
	default:
		return sentry.LevelError
	}
}
