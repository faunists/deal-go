package processors_test

import (
	"testing"

	"github.com/faunists/deal-go/processors"
)

func TestMakeExportedName(t *testing.T) {
	t.Parallel()

	tests := []struct {
		testName     string
		name         string
		expectedName string
	}{
		{
			testName:     "should capitalize the first letter to export name",
			name:         "someName",
			expectedName: "SomeName",
		},
		{
			testName:     "should do nothing when name is already exported",
			name:         "SomeName",
			expectedName: "SomeName",
		},
		{
			testName:     "should work when the string has length equal to one",
			name:         "s",
			expectedName: "S",
		},
		{
			testName:     "should work when the string has length equal to zero",
			name:         "",
			expectedName: "",
		},
	}

	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			actualName := processors.MakeExportedName(test.name)
			if actualName != test.expectedName {
				t.Errorf(
					"Wrong exported name formatting, given: %s, expected: %s",
					actualName, test.expectedName,
				)
			}
		})
	}
}
