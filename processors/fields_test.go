package processors_test

import (
	"testing"

	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/reflect/protoreflect"

	"github.com/faunists/deal-go/processors"
)

func TestFormatFieldValue(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		value          protoreflect.Value
		field          *protogen.Field
		expectedFormat string
	}{
		{
			name:           "should format correctly when value is boolean",
			value:          protoreflect.ValueOfBool(true),
			field:          nil,
			expectedFormat: "true",
		},
		{
			name:           "should format correctly when value is int32",
			value:          protoreflect.ValueOfInt32(32), //nolint:revive // random number
			field:          nil,
			expectedFormat: "32",
		},
		{
			name:           "should format correctly when value is int64",
			value:          protoreflect.ValueOfInt64(64), //nolint:revive // random number
			field:          nil,
			expectedFormat: "64",
		},
		{
			name:           "should format correctly when value is uint32",
			value:          protoreflect.ValueOfUint32(32), //nolint:revive // random number
			field:          nil,
			expectedFormat: "32",
		},
		{
			name:           "should format correctly when value is uint64",
			value:          protoreflect.ValueOfUint64(64), //nolint:revive // random number
			field:          nil,
			expectedFormat: "64",
		},
		{
			name:           "should format correctly when value is float32",
			value:          protoreflect.ValueOfFloat32(32.0), //nolint:revive // random number
			field:          nil,
			expectedFormat: "32.000000",
		},
		{
			name:           "should format correctly when value is float64",
			value:          protoreflect.ValueOfFloat32(64.0), //nolint:revive // random number
			field:          nil,
			expectedFormat: "64.000000",
		},
		{
			name:           "should format correctly when value is string",
			value:          protoreflect.ValueOfString("some-string"),
			field:          nil,
			expectedFormat: `"some-string"`,
		},
		{
			name:           "should format correctly when value is bytes",
			value:          protoreflect.ValueOfBytes([]byte("abcd")),
			field:          nil,
			expectedFormat: "[]byte{0x61, 0x62, 0x63, 0x64}",
		},
		{
			name: "should format correctly when value is EnumNumber",
			value: protoreflect.ValueOfEnum(
				protoreflect.EnumNumber(128), //nolint:revive // random number
			),
			field:          nil,
			expectedFormat: "128",
		},
		// TODO: Find a way to create a message and generate the field from them
		//{
		//	name:           "should format correctly when value is a message",
		//	value:          protoreflect.ValueOfMessage(simpleMessage),
		//	field:          nil,
		//	expectedFormat: `{Name: "Name", Age: 10}`,
		//},
		// TODO: Find a way to create a list and generate the field from them
		//{
		//	name: "should format correctly when value is List",
		//	value: protoreflect.ValueOfList(&testList{
		//		values: []protoreflect.Value{
		//			protoreflect.ValueOfInt64(42),
		//		},
		//	}),
		//	field:          nil,
		//	expectedFormat: "[]int64{42}",
		//},
	}

	identFunc := func(ident protogen.GoIdent) string {
		return ident.String()
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actualFormat, err := processors.FormatFieldValue(
				identFunc,
				test.field,
				test.value,
			)

			if err != nil {
				t.Errorf("Unespected error: %v", err)
			}

			if actualFormat != test.expectedFormat {
				t.Errorf(
					"Wrong format, given: %s expected %s",
					actualFormat, test.expectedFormat,
				)
			}
		})
	}
}
