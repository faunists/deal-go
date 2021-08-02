package processors

var allowedErrorCodeNames = []string{
	"OK",
	"Canceled", // It's not a typo here, this is the actual identifier in grpc codes
	"Unknown",
	"InvalidArgument",
	"DeadlineExceeded",
	"NotFound",
	"AlreadyExists",
	"PermissionDenied",
	"ResourceExhausted",
	"FailedPrecondition",
	"Aborted",
	"OutOfRange",
	"Unimplemented",
	"Internal",
	"Unavailable",
	"DataLoss",
	"Unauthenticated",
}

// IsErrorCodeValid returns true when a error code exists in the GRPC Codes package
func IsErrorCodeValid(errorCode string) bool {
	for _, allowedCode := range allowedErrorCodeNames {
		if errorCode == allowedCode {
			return true
		}
	}
	return false
}
