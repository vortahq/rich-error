package richerror

import "fmt"

// Type stores information that we want to show to user
type Type interface {
	String() string
}

// Kind hints about underlying cause of error
type Kind kind
type kind uint8

const (
	UnknownKind Kind = iota
	NotFound
	Unauthorized
	Unexpected
	Invalid
	Unavailable
)

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

// Operation represents the operation that caused the error
type Operation string

// Metadata stores metadata of error
type Metadata map[string]interface{}

// CodeInfo stores runtime information about the code
type CodeInfo struct {
	LineNumber   int    `json:"line_number,omitempty"`
	FileName     string `json:"file_name,omitempty"`
	FunctionName string `json:"function_name,omitempty"`
}

func (s *CodeInfo) String() string {
	return fmt.Sprintf("In %s:%s line %d",
		s.FileName,
		s.FunctionName,
		s.LineNumber)
}
