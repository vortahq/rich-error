package richerror

// RichError is a richer type of error that holds runtime information with itself
type RichError interface {
	Error() string
	Unwrap() error
	Is(target error) bool
	As(target interface{}) bool

	Metadata() Metadata
	CodeInfo() CodeInfo

	Operation() Operation
	Level() Level
	Type() Type
	Kind() Kind
}
