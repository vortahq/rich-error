package richerror

import (
	"context"
	"errors"
	"fmt"
	"runtime/debug"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// GRPCInterceptors is a helper that provides unary and stream grpc interceptors that will catch and log errors
// of your grpc server. If your grpc services return RichError it will set the grpc status code based on their Kind.
// Keep in mind that these interceptors will not log errors regarding the reflection API.
type GRPCInterceptors struct {
	Logger ErrorLogger
}

// UnaryInterceptor returns a gRPC unary interceptor that intercepts every gRPC request and in case of error prints
// (or logs) error and sets the grpc status code according to the error Kind. It also recovers panics.
func (h GRPCInterceptors) UnaryInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		defer func() {
			if e := h.recover(info.FullMethod); e != nil {
				h.log(info.FullMethod, e)
				err = h.getGPRCError(e)
			}
		}()

		if resp, err = handler(ctx, req); err != nil {
			h.log(info.FullMethod, err)
			err = h.getGPRCError(err)
			resp = nil
		}

		return
	}
}

// StreamInterceptor returns a gRPC stream interceptor that intercepts every gRPC request and in case of error prints
// (or logs) error and sets the grpc status code according to the error Kind. It also recovers panics.
func (h GRPCInterceptors) StreamInterceptor() grpc.StreamServerInterceptor {
	return func(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) (err error) {
		defer func() {
			if e := h.recover(info.FullMethod); e != nil {
				h.log(info.FullMethod, e)
				err = h.getGPRCError(e)
			}
		}()

		if err = handler(srv, stream); err != nil {
			h.log(info.FullMethod, err)
			err = h.getGPRCError(err)
		}

		return
	}
}

func (GRPCInterceptors) recover(path string) error {
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

func (GRPCInterceptors) getGPRCError(err error) error {
	var rErr RichError
	if !errors.As(err, &rErr) {
		message := err.Error()
		if rErr.Type() != nil {
			message = rErr.Type().String()
		}

		return status.Errorf(rErr.Kind().GRPCStatusCode(), "error: %s", message)
	}

	return status.Errorf(codes.Unknown, "error: %s", err.Error())
}

func (h GRPCInterceptors) log(path string, err error) {
	if strings.HasPrefix(path, "/grpc.reflection.v1alpha.ServerReflection/") {
		return
	}

	h.Logger.Log(err)
}
