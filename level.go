package richerror

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
