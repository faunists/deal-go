package entities

// Contract represents the root of everything that will be generated
type Contract struct {
	Name     string             `json:"name"`
	Services map[string]Service `json:"services"`
}

// Service is a named type of a map[string]Method
// The key represents the service name
type Service map[string]Method

// Method handles the two possibles cases to be generated:
//   - Success
//   - Failure
type Method struct {
	SuccessCases []SuccessCase `json:"successCases"`
	FailureCases []FailureCase `json:"failureCases"`
}

// SuccessCase handles the information about the request and response of a method
type SuccessCase struct {
	Description string      `json:"description"`
	Request     interface{} `json:"request"`
	Response    interface{} `json:"response"`
}

// FailureCase handles the information about the request and the error that should be returned
// for a given request
type FailureCase struct {
	Description string      `json:"description"`
	Request     interface{} `json:"request"`
	Error       GRPCError   `json:"error"`
}

// GRPCError handles the information about the error code and the string message of a GRPC error
type GRPCError struct {
	ErrorCode string `json:"errorCode"`
	Message   string `json:"message"`
}
