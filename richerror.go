package richerror

import (
	"errors"
	"fmt"
	"runtime"
	"strings"
)

type richError struct {
	wrappedError error
	message      string
	fields       Metadata

	runtimeInfo []RuntimeInfo

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
		fields:       make(map[string]interface{}),

		runtimeInfo: []RuntimeInfo{
			{
				LineNumber:   lineNumber,
				FileName:     fileName,
				FunctionName: functionName,
			},
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

// WithLevel assigns an error level to the error which can be used for log purposes
func (r *richError) WithLevel(level Level) *richError {
	r.level = level
	return r
}

// WithKind assigns an error kind to the error which can be used to decide the return code of the failed request
func (r *richError) WithKind(kind Kind) *richError {
	r.kind = kind
	return r
}

func (r *richError) WithOperation(operation Operation) *richError {
	r.operation = operation
	return r
}

// WithError wraps the underlying error and copies level, kind, type, and operation of the underlying error if not explicitly specified
func (r *richError) WithError(err error) *richError {
	r.wrappedError = err

	var wrappedRichError RichError
	ok := errors.As(err, &wrappedRichError)
	if !ok {
		return r
	}

	if r.level == UnknownLevel {
		r.level = wrappedRichError.Level()
	}

	if r.kind == UnknownKind {
		r.kind = wrappedRichError.Kind()
	}

	if r.operation == "" {
		r.operation = wrappedRichError.Operation()
	}

	if r._type == nil {
		r._type = wrappedRichError.Type()
	}

	for key, value := range wrappedRichError.Metadata() {
		if _, ok := r.fields[key]; !ok {
			r.fields[key] = value
		}
	}

	r.runtimeInfo = append(r.runtimeInfo, wrappedRichError.RuntimeInfo()...)

	return r
}

// NilIfNoError returns nil if wrapped error is nil, useful for direct return of the error
func (r *richError) NilIfNoError() RichError {
	if r.wrappedError == nil {
		return nil
	}

	return r
}

// String representation of the error
func (r *richError) String() string {
	return r.string(0)
}

func (r *richError) string(step int) string {
	msg := ``

	if r.message != "" {
		msg += fmt.Sprintf("message: %s ", r.message)
	}

	if r.operation != "" {
		msg += fmt.Sprintf("operation: %s ", r.operation)
	}

	if r.level != UnknownLevel {
		msg += fmt.Sprintf("level: %s ", r.level)
	}

	if r.kind != UnknownKind {
		msg += fmt.Sprintf("kind: %s ", r.kind)
	}

	if r._type != nil {
		msg += fmt.Sprintf("type: %s ", r._type)
	}

	if r.fields != nil {
		msg += fmt.Sprintf("fileds: %+v ", r.fields)
	}

	msg += fmt.Sprintf("code_info: %s ", r.runtimeInfo[0].String())

	if r.wrappedError != nil {
		innerError, ok := r.wrappedError.(*richError)
		if ok {
			return fmt.Sprintf("%s%s\n%s", strings.Repeat("\t", step), msg, innerError.string(step+1))
		}

		return fmt.Sprintf("%s%s\n%smessage: %s\n", strings.Repeat("\t", step), msg,
			strings.Repeat("\t", step+1), r.wrappedError.Error())
	}

	return fmt.Sprintf("%s%s\n", strings.Repeat("\t", step), msg)
}

func (r *richError) Error() string {
	if r.wrappedError == nil {
		return r.message
	}

	return fmt.Sprintf("%s -> %s", r.message, r.wrappedError.Error())
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
	return errors.As(r, &target)
}

func (r *richError) Metadata() Metadata {
	return r.fields
}

func (r *richError) RuntimeInfo() []RuntimeInfo {
	return r.runtimeInfo
}

func (r *richError) Operation() Operation {
	return r.operation
}

func (r *richError) Level() Level {
	if r.level == UnknownLevel {
		return Error
	}

	return r.level
}

func (r *richError) Type() Type {
	return r._type
}

func (r *richError) Kind() Kind {
	if r.kind == UnknownKind {
		return Unknown
	}

	return r.kind
}

// Deprecated: CodeInfo has been renamed to RuntimeInfo and will be removed in V2
func (r *richError) CodeInfo() CodeInfo {
	return CodeInfo{
		LineNumber:   r.runtimeInfo[0].LineNumber,
		FileName:     r.runtimeInfo[0].FileName,
		FunctionName: r.runtimeInfo[0].FunctionName,
	}
}
