package richerror

// RichError is a richer type of error that holds runtime information with itself
type RichError interface {
	String() string
	Error() string
	Unwrap() error
	Is(target error) bool
	As(target interface{}) bool

	Metadata() Metadata
	RuntimeInfo() RuntimeInfo

	Operation() Operation
	Level() Level
	Type() Type
	Kind() Kind

	// Deprecated: CodeInfo has been renamed to RuntimeInfo and will be removed in V2
	CodeInfo() CodeInfo
}
