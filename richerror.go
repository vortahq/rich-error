package richerror

import (
	"errors"
	"runtime"
)

type richError struct {
	wrappedError error
	message      string
	fields       Metadata

	codeInfo CodeInfo

	_type     Type
	level     Level
	kind      Kind
	operation Operation
}

// New creates a new richError
func New(message string) *richError {
	pc, fileName, lineNumber, _ := runtime.Caller(1)

	funcPt := runtime.FuncForPC(pc)
	functionName := "Unknown"
	if funcPt != nil {
		functionName = funcPt.Name()
	}

	return &richError{
		wrappedError: nil,
		message:      message,
		fields:       make(map[string]interface{}, 0),

		codeInfo: CodeInfo{
			LineNumber:   lineNumber,
			FileName:     fileName,
			FunctionName: functionName,
		},

		_type:     nil,
		level:     UnknownLevel,
		kind:      UnknownKind,
		operation: "",
	}
}

// WithFields appends given fields to already existing ones
func (r *richError) WithFields(fields Metadata) *richError {
	for key, value := range fields {
		r.fields[key] = value
	}
	return r
}

// WithField appends given field to already existing ones
func (r *richError) WithField(key string, value interface{}) *richError {
	r.fields[key] = value
	return r
}

// WithType specifies type of the error
func (r *richError) WithType(_type Type) *richError {
	r._type = _type
	return r
}

func (r *richError) WithLevel(level Level) *richError {
	r.level = level
	return r
}

func (r *richError) WithKind(kind Kind) *richError {
	r.kind = kind
	return r
}

func (r *richError) WithOperation(operation Operation) *richError {
	r.operation = operation
	return r
}

// WithError make the caller wraps the argument
func (r *richError) WithError(err error) *richError {
	r.wrappedError = err

	if re, ok := err.(RichError); ok {
		if r._type == nil {
			r._type = re.Type()
		}

		if r.level == UnknownLevel {
			r.level = re.Level()
		}

		if r.kind == UnknownKind {
			r.kind = re.Kind()
		}

		if r.operation == "" {
			r.operation = re.Operation()
		}
	}

	return r
}

// NilIfNoError returns nil if inner error is nil
func (r *richError) NilIfNoError() RichError {
	if r.Unwrap() == nil {
		return nil
	}

	return r
}

func (r *richError) Error() string {
	return r.message
}

func (r *richError) Unwrap() error {
	if innerError, ok := r.wrappedError.(interface{ Unwrap() error }); ok {
		return innerError.Unwrap()
	}

	return r.wrappedError
}

func (r *richError) Is(target error) bool {
	return errors.Is(r, target)
}

func (r *richError) As(target interface{}) bool {
	return errors.As(r, target)
}

func (r *richError) Metadata() Metadata {
	return r.fields
}

func (r *richError) CodeInfo() CodeInfo {
	return r.codeInfo
}

func (r *richError) Operation() Operation {
	return r.operation
}

func (r *richError) Level() Level {
	return r.level
}

func (r *richError) Type() Type {
	return r._type
}

func (r *richError) Kind() Kind {
	return r.kind
}
