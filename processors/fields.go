package processors

import (
	"fmt"

	"google.golang.org/protobuf/reflect/protoreflect"
)

// FormatFieldValue receives a value (protoreflect.Value) and converts it
// to a string format as an instantiate representation. See the example below:
//
// If value is a string the format will be myString -> "myString" because using
// double quote is the way we use to instantiate/create a string.
// Some other examples:
//   - int64: 64 -> 64 (to instantiate/create a int we just need the number)
//   - []bytes: []byte("abcd") -> []byte{0x61, 0x62, 0x63, 0x64}
func FormatFieldValue(value protoreflect.Value) string {
	// TODO: Make sure this works for proto enum, message, list and map
	// References:
	//   - protoreflect.Message
	//   - protoreflect.List
	//   - protoreflect.Map
	//   - protoreflect.EnumNumber
	switch v := value.Interface(); v.(type) {
	case float32, float64:
		return fmt.Sprintf("%f", v)
	case string:
		return fmt.Sprintf("%q", v)
	case []byte:
		return fmt.Sprintf("%#v", v)
	case protoreflect.Message, protoreflect.List, protoreflect.Map, protoreflect.EnumNumber:
		return fmt.Sprintf(`"Unsupported type: %T"`, v)
	default:
		return fmt.Sprintf("%v", v)
	}
}
