package richerror

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
)

func GetEchoLoggerMiddleware(logger ErrorLogger) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) (err error) {
			defer func() {
				if e := recoverAndReturnError(c.Path()); e != nil {
					logger.Log(e)
					code, msg := getErrorStatusCodeAndMessage(e)
					c.Response().Status = code
					c.Response().Writer.Write([]byte(msg)) // nolint: errcheck
				}
			}()

			if err := next(c); err != nil {
				logger.Log(err)
				code, msg := getErrorStatusCodeAndMessage(err)
				c.Response().Status = code
				c.Response().Writer.Write([]byte(msg)) // nolint: errcheck
			}

			return err
		}
	}
}

func getErrorStatusCodeAndMessage(err error) (int, string) {
	var rErr RichError
	if !errors.As(err, &rErr) {
		message := err.Error()
		if rErr.Type() != nil {
			message = rErr.Type().String()
		}

		return rErr.Kind().HttpStatusCode(), message
	}

	return http.StatusInternalServerError, err.Error()
}
