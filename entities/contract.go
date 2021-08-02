package entities

type Contract struct {
	Name     string             `json:"name"`
	Services map[string]Service `json:"services"`
}

type Service map[string]Method

type Method struct {
	SuccessCases []SuccessCase `json:"successCases"`
	FailureCases []FailureCase `json:"failureCases"`
}

type SuccessCase struct {
	Description string      `json:"description"`
	Request     interface{} `json:"request"`
	Response    interface{} `json:"response"`
}

type FailureCase struct {
	Description string      `json:"description"`
	Request     interface{} `json:"request"`
	Error       GRPCError   `json:"error"`
}

type GRPCError struct {
	ErrorCode string `json:"errorCode"`
	Message   string `json:"message"`
}
