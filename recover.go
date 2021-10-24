package richerror

import (
	"fmt"
	"runtime/debug"
)

func recoverAndReturnError(path string) error {
	if r := recover(); r != nil {
		errType := StringType(fmt.Sprintf("panic: %s", r))
		err := New("panic detected").WithType(errType).WithFields(Metadata{
			"path":        path,
			"panic":       r,
			"stack_trace": debug.Stack(),
		})
		return err
	}

	return nil
}
