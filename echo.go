package richerror

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4/middleware"

	"github.com/labstack/echo/v4"
)

func GetEchoLoggerMiddleware(logger ErrorLogger) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		middleware.Recover()
		return func(c echo.Context) (err error) {
			defer func() {
				if e := recoverAndReturnError(c.Path()); e != nil {
					logger.Log(e)
					c.Error(e)
				}
			}()

			if err := next(c); err != nil {
				logger.Log(err)
				code, msg := getErrorStatusCodeAndMessage(err)
				return echo.NewHTTPError(code, msg)
			}

			return nil
		}
	}
}

func getErrorStatusCodeAndMessage(err error) (int, string) {
	var rErr RichError
	if errors.As(err, &rErr) {
		message := err.Error()
		if rErr.Type() != nil {
			message = rErr.Type().String()
		}

		return rErr.Kind().HttpStatusCode(), message
	}

	return http.StatusInternalServerError, err.Error()
}
