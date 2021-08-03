package processors_test

import (
	"testing"

	"github.com/faunists/deal-go/processors"
)

func TestIsErrorCodeValid(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		errorCode     string
		expectedValue bool
	}{
		{
			name:          `"OK" should be a valid option`,
			errorCode:     "OK",
			expectedValue: true,
		},
		{
			name:          `"Canceled" should be a valid option`,
			errorCode:     "Canceled",
			expectedValue: true,
		},
		{
			name:          `"Unknown" should be a valid option`,
			errorCode:     "Unknown",
			expectedValue: true,
		},
		{
			name:          `"InvalidArgument" should be a valid option`,
			errorCode:     "InvalidArgument",
			expectedValue: true,
		},
		{
			name:          `"DeadlineExceeded" should be a valid option`,
			errorCode:     "DeadlineExceeded",
			expectedValue: true,
		},
		{
			name:          `"NotFound" should be a valid option`,
			errorCode:     "NotFound",
			expectedValue: true,
		},
		{
			name:          `"AlreadyExists" should be a valid option`,
			errorCode:     "AlreadyExists",
			expectedValue: true,
		},
		{
			name:          `"PermissionDenied" should be a valid option`,
			errorCode:     "PermissionDenied",
			expectedValue: true,
		},
		{
			name:          `"ResourceExhausted" should be a valid option`,
			errorCode:     "ResourceExhausted",
			expectedValue: true,
		},
		{
			name:          `"FailedPrecondition" should be a valid option`,
			errorCode:     "FailedPrecondition",
			expectedValue: true,
		},
		{
			name:          `"Aborted" should be a valid option`,
			errorCode:     "Aborted",
			expectedValue: true,
		},
		{
			name:          `"OutOfRange" should be a valid option`,
			errorCode:     "OutOfRange",
			expectedValue: true,
		},
		{
			name:          `"Unimplemented" should be a valid option`,
			errorCode:     "Unimplemented",
			expectedValue: true,
		},
		{
			name:          `"Internal" should be a valid option`,
			errorCode:     "Internal",
			expectedValue: true,
		},
		{
			name:          `"Unavailable" should be a valid option`,
			errorCode:     "Unavailable",
			expectedValue: true,
		},
		{
			name:          `"DataLoss" should be a valid option`,
			errorCode:     "DataLoss",
			expectedValue: true,
		},
		{
			name:          `"Unauthenticated" should be a valid option`,
			errorCode:     "Unauthenticated",
			expectedValue: true,
		},
		{
			name:          `"MyTest" should not be a valid option`,
			errorCode:     "MyTest",
			expectedValue: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actualValue := processors.IsErrorCodeValid(test.errorCode)
			if actualValue != test.expectedValue {
				t.Errorf(
					"Given: %v, expected: %v",
					actualValue, test.expectedValue,
				)
			}
		})
	}
}
