package processors

import (
	"fmt"
	"strings"
)

// MakeExportedName transforms any string in a Go's exported name.
// In Go to make a function/struct/const/var public the first letter must be upper.
func MakeExportedName(name string) string {
	if len(name) == 0 {
		return ""
	}

	return fmt.Sprintf(
		"%s%s",
		strings.ToUpper(name[:1]),
		name[1:],
	)
}
