//go:build codegen

package protonats

//go:generate go run xiam.li/errnumer/cmd/errnumer -type=errnumerErr -output=errors_enum.go
type errnumerErr struct {
	Code        string
	Description string
}

var (
	errnumerErrInternal            = errnumerErr{"500", "Internal server error."}
	errnumerErrUnmarshallingFailed = errnumerErr{"501", "Failed to unmarshal proto message."}
	errnumerErrMarshallingFailed   = errnumerErr{"502", "Failed to marshal proto message."}
	errnumerErrPanicRecovered      = errnumerErr{"503", "A panic occurred and was recovered."}
	errnumerErrValidationFailed    = errnumerErr{"400", "Validation Failed."}
	errnumerErrUnauthorized        = errnumerErr{"401", "Unauthorized."}
	errnumerErrForbidden           = errnumerErr{"403", "Forbidden."}
)
