package richerror

import "encoding/json"

// JsonMode is a flag controlling the format of generated output, the default format is string
var JsonMode = false

// IsJson returns the JsonMode flag value.
func IsJson() bool {
	return JsonMode
}

func (r *richError) json() []byte {
	b, err := json.Marshal(r)
	if err != nil {
		return []byte(err.Error())
	}

	return b
}

type richErrorJson struct {
	Message      string   `json:"message,omitempty"`
	Operation    string   `json:"operation,omitempty"`
	Level        uint8    `json:"level,omitempty"`
	Kind         uint8    `json:"kind,omitempty"`
	Type         string   `json:"type,omitempty"`
	Fields       Metadata `json:"fields,omitempty"`
	CodeInfo     CodeInfo `json:"code_info,omitempty"`
	WrappedError error    `json:"wrapped_error,omitempty"`
}

func (r *richError) MarshalJSON() ([]byte, error) {
	j := &richErrorJson{
		Message:   r.message,
		Operation: string(r.operation),
		Level:     uint8(r.level),
		Kind:      uint8(r.kind),
		Fields:    r.fields,
		CodeInfo:  r.codeInfo,
	}

	if r.Type() != nil {
		j.Type = r.Type().String()
	}

	if r.wrappedError != nil {
		if _, ok := r.wrappedError.(json.Marshaler); ok {
			j.WrappedError = r.wrappedError
		} else {
			j.WrappedError = &simpleError{Message: r.wrappedError.Error()}
		}
	}

	return json.Marshal(j)
}

type simpleError struct {
	Message string `json:"message"`
}

func (s *simpleError) Error() string {
	return s.Message
}
