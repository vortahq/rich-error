package richerror

import "fmt"

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
