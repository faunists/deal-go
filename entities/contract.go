package entities

import "fmt"

// Contract represents the root of everything that will be generated
type Contract struct {
	Name     string             `json:"name" yaml:"name"`
	Services map[string]Service `json:"services" yaml:"services"`
}

// Service is a named type of a map[string]Method
// The key represents the service name
type Service map[string]Method

// Method handles the two possibles cases to be generated:
//   - Success
//   - Failure
type Method struct {
	SuccessCases []SuccessCase `json:"successCases" yaml:"successCases"`
	FailureCases []FailureCase `json:"failureCases" yaml:"failureCases"`
}

// SuccessCase handles the information about the request and response of a method
type SuccessCase struct {
	Description string      `json:"description" yaml:"description"`
	Request     interface{} `json:"request" yaml:"request"`
	Response    interface{} `json:"response" yaml:"response"`
}

// FailureCase handles the information about the request and the error that should be returned
// for a given request
type FailureCase struct {
	Description string      `json:"description" yaml:"description"`
	Request     interface{} `json:"request" yaml:"request"`
	Error       GRPCError   `json:"error" yaml:"error"`
}

// GRPCError handles the information about the error code and the string message of a GRPC error
type GRPCError struct {
	ErrorCode string `json:"errorCode" yaml:"errorCode"`
	Message   string `json:"message" yaml:"message"`
}

func (e GRPCError) String() string {
	return fmt.Sprintf("rpc error: code = %s desc = %s", e.ErrorCode, e.Message)
}
