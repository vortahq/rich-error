package richerror

import "google.golang.org/grpc/codes"

// Kind hints about underlying cause of error
type Kind kind
type kind uint8

const (
	UnknownKind      Kind = iota
	Canceled              // grpc: CANCELLED 			- http 500 Internal Server Error
	Unknown               // grpc: UNKNOWN 				- http 500 Internal Server Error
	InvalidArgument       // grpc: INVALID_ARGUMENT 	- http 400 Bad Request
	Timeout               // grpc: DEADLINE_EXCEEDED 	- http 500 Internal Server Error
	NotFound              // grpc: NOT_FOUND 			- http 404 Not Found
	AlreadyExists         // grpc: ALREADY_EXISTS		- http 409 Conflict
	PermissionDenied      // grpc: PERMISSION_DENIED 	- http 403 Forbidden
	TooManyRequests       // grpc: RESOURCE_EXHAUSTED 	- http 429 Too Many Requests
	Unimplemented         // grpc: UNIMPLEMENTED		- http 501 Not Implemented
	Internal              // grpc: INTERNAL 			- http 500 Internal Server Error
	Unavailable           // grpc: UNAVAILABLE 			- http 503 Service Unavailable
	Unauthenticated       // grpc: UNAUTHENTICATED 		- http 401 Unauthorized

	Unauthorized = Unauthenticated
	Invalid      = InvalidArgument
	Unexpected   = Unknown
)

var kindStrings = [...]string{"_", "Canceled", "Unknown", "Invalid Argument", "Timeout", "NotFound", "Already Exists",
	"Permission Denied", "Too Many Requests", "Unimplemented", "Internal", "Unavailable", "Unauthenticated"}

func (k Kind) String() string {
	return kindStrings[k]
}

func (k Kind) MarshalJSON() ([]byte, error) {
	return []byte(k.String()), nil
}

func (k Kind) GRPCStatusCode() codes.Code {
	switch k {
	case Canceled:
		return codes.Canceled
	case Unknown:
		return codes.Unknown
	case InvalidArgument:
		return codes.InvalidArgument
	case Timeout:
		return codes.DeadlineExceeded
	case NotFound:
		return codes.NotFound
	case AlreadyExists:
		return codes.AlreadyExists
	case PermissionDenied:
		return codes.PermissionDenied
	case TooManyRequests:
		return codes.ResourceExhausted
	case Unimplemented:
		return codes.Unimplemented
	case Internal:
		return codes.Internal
	case Unavailable:
		return codes.Unavailable
	case Unauthenticated:
		return codes.Unauthenticated
	default:
		return codes.Unknown
	}
}
