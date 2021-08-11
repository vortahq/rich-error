package helpers

import (
	"fmt"
	"runtime/debug"

	richerror "gitlab.com/orderhq/rich-error"
)

type helper struct {
	basicLogger     basicLogger
	formattedLogger formattedLogger
	contextLogger   contextLogger

	sentryEnabled bool
	environment   string
	serverName    string
}

func New() *helper {
	return &helper{}
}

func (h *helper) WithLogger(logger basicLogger) {
	h.basicLogger = logger
}

func (h *helper) WithFormattedLogger(logger formattedLogger) {
	h.formattedLogger = logger
}

func (h *helper) WithContextLogger(logger contextLogger) {
	h.contextLogger = logger
}

func (h *helper) WithSentry() {
	h.sentryEnabled = true
}

func (h *helper) WithEnvironment(env string) {
	h.environment = env
}

func (h *helper) WithServerName(name string) {
	h.serverName = name
}

func (h helper) recover(path string) error {
	if r := recover(); r != nil {
		errType := richerror.StringType(fmt.Sprintf("panic: %s", r))
		err := richerror.New("panic detected").WithType(errType).WithFields(richerror.Metadata{
			"path":        path,
			"panic":       r,
			"stack_trace": debug.Stack(),
		})
		return err
	}

	return nil
}
