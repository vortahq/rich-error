package helpers

import (
	"context"
	"errors"

	richerror "gitlab.com/orderhq/rich-error"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// UnaryInterceptor returns a gRPC unary interceptor that intercepts every gRPC request and in case of error prints
// (or logs) error and sets the grpc status code according to the error Kind. It also recovers panics.
func (h helper) UnaryInterceptor() grpc.UnaryServerInterceptor {
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
func (h helper) StreamInterceptor() grpc.StreamServerInterceptor {
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

func (h helper) getGPRCError(err error) error {
	var rErr richerror.RichError
	if !errors.As(err, &rErr) {
		message := err.Error()
		if rErr.Type() != nil {
			message = rErr.Type().String()
		}

		return status.Errorf(rErr.Kind().GRPCStatusCode(), "error: %s", message)
	}

	return status.Errorf(codes.Unknown, "error: %s", err.Error())
}
