package processors

import (
	"errors"
	"fmt"
	"strings"

	"google.golang.org/protobuf/compiler/protogen"

	"google.golang.org/protobuf/reflect/protoreflect"
)

// IdentFunc is used when we need to use the QualifiedGoIdent method
// from protogen.GeneratedFile without passing the entire struct.
type IdentFunc func(ident protogen.GoIdent) string

// FieldsByNumber represents the relation between the proto field number
// with a protogen.Field. It's very useful to correlate with fields from
// protoreflect.Message.
type FieldsByNumber map[protoreflect.FieldNumber]*protogen.Field

// FormatFieldValue receives a value (protoreflect.Value) and converts it
// to a string format as an instantiate representation. See the example below:
//
// If value is a string the format will be myString -> "myString" because using
// double quote is the way we use to instantiate/create a string.
// Some other examples:
//   - int64: 64 -> 64 (to instantiate/create a int we just need the number)
//   - []bytes: []byte("abcd") -> []byte{0x61, 0x62, 0x63, 0x64}
func FormatFieldValue(
	identFunc IdentFunc,
	field *protogen.Field,
	value protoreflect.Value,
) (string, error) {
	// TODO: Make sure this works for proto enum, list, and map
	// References:
	//   - protoreflect.List
	//   - protoreflect.Map
	switch v := value.Interface(); v.(type) {
	case float32, float64:
		return fmt.Sprintf("%f", v), nil
	case string:
		return fmt.Sprintf("%q", v), nil
	case []byte:
		return fmt.Sprintf("%#v", v), nil
	case protoreflect.EnumNumber:
		return fmt.Sprintf("%d", v), nil
	case protoreflect.Message:
		fieldsByNumber := CreateFieldsByNumber(field.Message.Fields)
		value.Message()

		return FormatMessageFieldNew(
			identFunc,
			field.Message.GoIdent,
			fieldsByNumber,
			value.Message(),
		)
	case protoreflect.List, protoreflect.Map:
		return "", errors.New(fmt.Sprintf(`"Unsupported type: %T"`, v))
	default:
		return fmt.Sprintf("%v", v), nil
	}
}

func FormatMessageFieldNew(
	identFunc IdentFunc,
	ident protogen.GoIdent,
	fieldsByNumber FieldsByNumber,
	message protoreflect.Message,
) (string, error) {
	var err error

	messageArguments := make([]string, 0)
	message.Range(
		func(descriptor protoreflect.FieldDescriptor, value protoreflect.Value) bool {
			field, exists := fieldsByNumber[descriptor.Number()]
			if !exists {
				err = fmt.Errorf(
					"field not found %s while inspecting message %s",
					descriptor.Name(), message.Type().Descriptor().FullName(),
				)
				return false
			}

			// We need to declare formattedField before assign the return of
			// FormatFieldValue to it because using the shorthand assign `:=`
			// we will lose the closure of `err`.
			var formattedField string
			formattedField, err = FormatFieldValue(identFunc, field, value)
			if err != nil {
				return false
			}

			messageArguments = append(
				messageArguments,
				fmt.Sprintf("%s: %s", field.GoName, formattedField),
			)

			return true
		},
	)

	return fmt.Sprintf(
		"&%s{%s}",
		identFunc(ident),
		strings.Join(messageArguments, ", "),
	), err
}

// CreateFieldsByNumber transform a slice of protogen.Field into a FieldsByNumber,
// so we can access the field by its number. This is very handy when we need to correlate
// fields from a protogen.Message with fields from a protoreflect.Message.
func CreateFieldsByNumber(fields []*protogen.Field) FieldsByNumber {
	fieldsByNumber := make(FieldsByNumber)
	for _, field := range fields {
		fieldsByNumber[field.Desc.Number()] = field
	}

	return fieldsByNumber
}
