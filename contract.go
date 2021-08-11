package richerror

import "google.golang.org/grpc"

// RichError is a richer type of error that holds runtime information with itself
type RichError interface {
	String() string
	Error() string
	Unwrap() error
	Is(target error) bool
	As(target interface{}) bool

	Metadata() Metadata
	RuntimeInfo() []RuntimeInfo

	Operation() Operation
	Level() Level
	Type() Type
	Kind() Kind

	// Deprecated: CodeInfo has been renamed to RuntimeInfo and will be removed in V2
	CodeInfo() CodeInfo
}

// Operation can be used to group or organize error
type Operation string

// Metadata stores metadata of error
type Metadata map[string]interface{}

type GRPCHelper interface {
	UnaryInterceptor() grpc.UnaryServerInterceptor
	StreamInterceptor() grpc.StreamServerInterceptor
}
