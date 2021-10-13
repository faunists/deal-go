package processors_test

import (
	"testing"

	"google.golang.org/protobuf/reflect/protoreflect"

	"github.com/faunists/deal-go/processors"
)

func TestFormatFieldValue(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		value          protoreflect.Value
		expectedFormat string
	}{
		{
			name:           "should format correctly when value is boolean",
			value:          protoreflect.ValueOfBool(true),
			expectedFormat: "true",
		},
		{
			name:           "should format correctly when value is int32",
			value:          protoreflect.ValueOfInt32(32), //nolint:revive // random number
			expectedFormat: "32",
		},
		{
			name:           "should format correctly when value is int64",
			value:          protoreflect.ValueOfInt64(64), //nolint:revive // random number
			expectedFormat: "64",
		},
		{
			name:           "should format correctly when value is uint32",
			value:          protoreflect.ValueOfUint32(32), //nolint:revive // random number
			expectedFormat: "32",
		},
		{
			name:           "should format correctly when value is uint64",
			value:          protoreflect.ValueOfUint64(64), //nolint:revive // random number
			expectedFormat: "64",
		},
		{
			name:           "should format correctly when value is float32",
			value:          protoreflect.ValueOfFloat32(32.0), //nolint:revive // random number
			expectedFormat: "32.000000",
		},
		{
			name:           "should format correctly when value is float64",
			value:          protoreflect.ValueOfFloat32(64.0), //nolint:revive // random number
			expectedFormat: "64.000000",
		},
		{
			name:           "should format correctly when value is string",
			value:          protoreflect.ValueOfString("some-string"),
			expectedFormat: `"some-string"`,
		},
		{
			name:           "should format correctly when value is bytes",
			value:          protoreflect.ValueOfBytes([]byte("abcd")),
			expectedFormat: "[]byte{0x61, 0x62, 0x63, 0x64}",
		},
		{
			name: "should format correctly when value is EnumNumber",
			value: protoreflect.ValueOfEnum(
				protoreflect.EnumNumber(128), //nolint:revive // random number
			),
			expectedFormat: "128",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actualFormat := processors.FormatFieldValue(test.value)
			if actualFormat != test.expectedFormat {
				t.Errorf(
					"Wrong format, given: %s expected %s",
					actualFormat, test.expectedFormat,
				)
			}
		})
	}
}
