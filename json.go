package richerror

import "encoding/json"

type richErrorJson struct {
	Message      string      `json:"message,omitempty"`
	Operation    Operation   `json:"operation,omitempty"`
	Level        Level       `json:"level,omitempty"`
	Kind         Kind        `json:"kind,omitempty"`
	Type         string      `json:"type,omitempty"`
	Fields       Metadata    `json:"fields,omitempty"`
	RuntimeInfo  RuntimeInfo `json:"runtime_info,omitempty"`
	WrappedError error       `json:"wrapped_error,omitempty"`
}

func (r *richError) MarshalJSON() ([]byte, error) {
	jsonStruct := &richErrorJson{
		Message:     r.message,
		Operation:   r.operation,
		Level:       r.level,
		Kind:        r.kind,
		Fields:      r.fields,
		RuntimeInfo: r.runtimeInfo,
	}

	if r.Type() != nil {
		jsonStruct.Type = r.Type().String()
	}

	if r.wrappedError != nil {
		if _, ok := r.wrappedError.(json.Marshaler); ok {
			jsonStruct.WrappedError = r.wrappedError
		} else {
			jsonStruct.WrappedError = &simpleError{Message: r.wrappedError.Error()}
		}
	}

	return json.Marshal(jsonStruct)
}

type simpleError struct {
	Message string `json:"message"`
}

func (s *simpleError) Error() string {
	return s.Message
}
