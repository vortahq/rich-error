package richerror

import (
	"fmt"
)

// Type stores information that we want to show to user
type Type interface {
	String() string
}

// Kind hints about underlying cause of error
type Kind kind
type kind uint8

const (
	UnknownKind      Kind = iota
	Canceled              // grpc: CANCELLED 			- http 500 Internal Server Error
	Unknown               // grpc: UNKNOWN 				- http 500 Internal Server Error
	InvalidArgument       // grpc: INVALID_ARGUMENT 	- http 400 Bad Request
	Timeout               // grpc: DEADLINE_EXCEEDED 	- http 500 Internal Server Error
	NotFound              // grpc: NOT_FOUND 			- http 404 Not Found
	AlreadyExists         // grpc: ALREADY_EXISTS		- http 409 Conflict
	PermissionDenied      // grpc: PERMISSION_DENIED 	- http 403 Forbidden
	TooManyRequests       // grpc: RESOURCE_EXHAUSTED 	- http 429 Too Many Requests
	Unimplemented         // grpc: UNIMPLEMENTED		- http 501 Not Implemented
	Internal              // grpc: INTERNAL 			- http 500 Internal Server Error
	Unavailable           // grpc: UNAVAILABLE 			- http 503 Service Unavailable
	Unauthenticated       // grpc: UNAUTHENTICATED 		- http 401 Unauthorized

	Unauthorized = Unauthenticated
	Invalid      = InvalidArgument
	Unexpected   = Unknown
)

var kindStrings = [...]string{"_", "Canceled", "Unknown", "Invalid Argument", "Timeout", "NotFound", "Already Exists",
	"Permission Denied", "Too Many Requests", "Unimplemented", "Internal", "Unavailable", "Unauthenticated"}

func (k Kind) String() string {
	return kindStrings[k]
}

func (k Kind) MarshalJSON() ([]byte, error) {
	return []byte(k.String()), nil
}

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

// Operation can be used to group or organize error
type Operation string

// Metadata stores metadata of error
type Metadata map[string]interface{}

// RuntimeInfo stores runtime information about the code
type RuntimeInfo struct {
	LineNumber   int    `json:"line_number,omitempty"`
	FileName     string `json:"file_name,omitempty"`
	FunctionName string `json:"function_name,omitempty"`
}

func (s *RuntimeInfo) String() string {
	return fmt.Sprintf("In %s:%s line %d",
		s.FileName,
		s.FunctionName,
		s.LineNumber)
}

// Deprecated: CodeInfo has been renamed to RuntimeInfo and will be removed in V2
type CodeInfo RuntimeInfo
