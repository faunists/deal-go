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
	switch v := value.Interface(); v.(type) {
	case float32, float64:
		return fmt.Sprintf("%f", v), nil
	case string:
		return fmt.Sprintf("%q", v), nil
	case []byte:
		return fmt.Sprintf("%#v", v), nil
	case protoreflect.EnumNumber:
		enum := field.Enum

		if value.Enum() < 0 || int(value.Enum()) >= len(enum.Values) {
			return "", fmt.Errorf("enum option out of range for '%s'", enum.Desc.Name())
		}

		return fmt.Sprintf("%s", identFunc(enum.Values[value.Enum()].GoIdent)), nil
	case protoreflect.Message:
		fieldsByNumber := CreateFieldsByNumber(field.Message.Fields)
		value.Message()

		return FormatMessageField(
			identFunc,
			field.Message.GoIdent,
			fieldsByNumber,
			value.Message(),
		)
	case protoreflect.List:
		return formatList(identFunc, field, value)
	case protoreflect.Map:
		return formatMap(identFunc, field, value)
	default:
		return fmt.Sprintf("%v", v), nil
	}
}

// FormatMessageField takes care of formatting a message to
// a properly string format.
func FormatMessageField(
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

func formatList(
	identFunc IdentFunc,
	field *protogen.Field,
	value protoreflect.Value,
) (string, error) {
	list := value.List()

	formattedValues := make([]string, 0, list.Len())
	for i := 0; i < list.Len(); i++ {
		item := list.Get(i)
		formattedValue, err := FormatFieldValue(identFunc, field, item)
		if err != nil {
			return "", err
		}
		formattedValues = append(formattedValues, formattedValue)
	}

	switch field.Desc.Kind() {
	case protoreflect.GroupKind:
		return "", errors.New("we don't support groups yet")
	case protoreflect.MessageKind:
		return fmt.Sprintf(
			"[]*%s{%s}",
			identFunc(field.Message.GoIdent),
			strings.Join(formattedValues, ", "),
		), nil
	default:
		return fmt.Sprintf(
			"[]%s{%s}",
			field.Desc.Kind().String(),
			strings.Join(formattedValues, ", "),
		), nil
	}
}

func formatMap(
	identFunc IdentFunc,
	field *protogen.Field,
	value protoreflect.Value,
) (string, error) {
	var err error

	m := value.Map()

	// A `Map` is represented by a message which contains two fields,
	// you can use the `field.Message.Desc.IsMapEntry` method to verify whether
	// a message represents a map. Explanation about the fields extracted from godoc:
	//
	// Map entry messages have only two fields:
	//	• a "key" field with a field number of 1
	//	• a "value" field with a field number of 2
	// The key and value types are determined by these two fields.
	//
	// nolint:lll // link
	// Source: https://pkg.go.dev/google.golang.org/protobuf/reflect/protoreflect#MessageDescriptor
	keyField := field.Message.Fields[0]
	valueField := field.Message.Fields[1]

	formattedValues := make([]string, 0, m.Len())
	m.Range(func(key protoreflect.MapKey, insideValue protoreflect.Value) bool {
		var formattedKey string
		formattedKey, err = FormatFieldValue(identFunc, valueField, key.Value())
		if err != nil {
			return false
		}

		var formattedValue string
		formattedValue, err = FormatFieldValue(identFunc, valueField, insideValue)
		if err != nil {
			return false
		}

		formattedValues = append(
			formattedValues,
			fmt.Sprintf("%s: %s", formattedKey, formattedValue),
		)

		return true
	})

	if err != nil {
		return "", err
	}

	switch valueField.Desc.Kind() {
	case protoreflect.GroupKind:
		return "", errors.New("we don't support groups yet")
	case protoreflect.MessageKind:
		return fmt.Sprintf(
			"map[%s]*%s{%s}",
			keyField.Desc.Kind(),
			identFunc(valueField.Message.GoIdent),
			strings.Join(formattedValues, ", "), //nolint:revive // don't need a const for sep
		), nil
	default:
		return fmt.Sprintf(
			"map[%s]%s{%s}",
			keyField.Desc.Kind(),
			valueField.Desc.Kind(),
			strings.Join(formattedValues, ", "), //nolint:revive // don't need a const for sep
		), nil
	}
}
