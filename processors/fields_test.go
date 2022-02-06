package processors_test

import (
	"testing"

	"google.golang.org/protobuf/encoding/protojson"

	"google.golang.org/protobuf/types/dynamicpb"

	"github.com/faunists/deal-go/mocks"

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
				protoreflect.EnumNumber(1), //nolint:revive // random number
			),
			field: protoFields.getField(
				t, "MessageWithComplexFields", "enumField",
			),
			expectedFormat: "EnumNumbers_TWO",
		},
		{
			name: "should format correctly when value is a message",
			value: protoreflect.ValueOfMessage(
				getMessage(
					t,
					"SimpleMessage",
					readFixture(t, "simple_message_value.json"),
				),
			),
			field: &protogen.Field{
				Message: protoFields.getMessage(t, "SimpleMessage"),
			},
			expectedFormat: `&SimpleMessage{IntField: 42}`,
		},
		{
			name: "should format correctly when value is List with a single element",
			value: protoreflect.ValueOfList(&mocks.ProtoList{
				Values: []protoreflect.Value{
					protoreflect.ValueOfString("awesome string"),
				},
			}),
			field: protoFields.getField(
				t, "MessageWithComplexFields", "stringListField",
			),
			expectedFormat: `[]string{"awesome string"}`,
		},
		{
			name: "should format correctly when value is List with more than one element",
			value: protoreflect.ValueOfList(&mocks.ProtoList{
				Values: []protoreflect.Value{
					protoreflect.ValueOfString("first string"),
					protoreflect.ValueOfString("second string"),
				},
			}),
			field: protoFields.getField(
				t, "MessageWithComplexFields", "stringListField",
			),
			expectedFormat: `[]string{"first string", "second string"}`,
		},
		{
			name: "should format correctly when value is List with a message",
			value: protoreflect.ValueOfList(&mocks.ProtoList{
				Values: []protoreflect.Value{
					protoreflect.ValueOfMessage(
						getMessage(
							t,
							"SimpleMessage",
							readFixture(t, "simple_message_value.json"),
						),
					),
				},
			}),
			field: protoFields.getField(
				t, "MessageWithComplexFields", "listSimpleMessageField",
			),
			expectedFormat: `[]*SimpleMessage{&SimpleMessage{IntField: 42}}`,
		},
		{
			name: "should format correctly when value is Map",
			value: protoreflect.ValueOfMap(&mocks.ProtoMap{
				Map: map[interface{}]protoreflect.Value{
					int64(42): protoreflect.ValueOfString("test"), //nolint:revive
				},
			}),
			field: protoFields.getField(
				t, "MessageWithComplexFields", "mapField",
			),
			expectedFormat: `map[int64]string{42: "test"}`,
		},
		{
			name: "should format correctly when value is Map and the value is a message",
			value: protoreflect.ValueOfMap(&mocks.ProtoMap{
				Map: map[interface{}]protoreflect.Value{
					"my value": protoreflect.ValueOfMessage(
						getMessage(
							t,
							"SimpleMessage",
							readFixture(t, "simple_message_value.json"),
						),
					),
				},
			}),
			field: protoFields.getField(
				t, "MessageWithComplexFields", "mapSimpleMessageField",
			),
			expectedFormat: `map[string]*SimpleMessage{"my value": &SimpleMessage{IntField: 42}}`,
		},
	}

	identFunc := func(ident protogen.GoIdent) string {
		return ident.GoName
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actualFormat, err := processors.FormatFieldValue(
				identFunc,
				test.field,
				test.value,
			)

			if err != nil {
				t.Fatalf("Unespected error: %v", err)
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

func TestFormatFieldValue_Errors(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		value         protoreflect.Value
		field         *protogen.Field
		expectedError string
	}{
		{
			name: "should return an error when enum value is out of the range - positive",
			value: protoreflect.ValueOfEnum(
				protoreflect.EnumNumber(999), //nolint:revive // random number
			),
			field: protoFields.getField(
				t, "MessageWithComplexFields", "enumField",
			),
			expectedError: "enum option out of range for 'EnumNumbers'",
		},
		{
			name: "should return an error when enum value is out of the range - negative",
			value: protoreflect.ValueOfEnum(
				protoreflect.EnumNumber(-1), //nolint:revive // random number
			),
			field: protoFields.getField(
				t, "MessageWithComplexFields", "enumField",
			),
			expectedError: "enum option out of range for 'EnumNumbers'",
		},
	}

	identFunc := func(ident protogen.GoIdent) string {
		return ident.GoName
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, err := processors.FormatFieldValue(
				identFunc,
				test.field,
				test.value,
			)

			if err == nil {
				t.Fatal("expected an error but none was returned")
			}

			if err.Error() != test.expectedError {
				t.Fatalf(`wrong error, given: "%s" expected: "%s"`, err, test.expectedError)
			}
		})
	}
}

func getMessage(t *testing.T, messageName string, messageValue []byte) *dynamicpb.Message {
	t.Helper()

	protoMessage := protoFields.getMessage(t, messageName)

	message := dynamicpb.NewMessage(protoMessage.Desc)
	if err := protojson.Unmarshal(messageValue, message); err != nil {
		t.Fatalf("error populating message: %v", err)
	}

	return message
}
