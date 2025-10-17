// core/internal/naming/naming.go  (novo)
package naming

import (
	"reflect"
)

func TypeNameOf(v any) string {
	if v == nil {
		return "<nil>"
	}
	t := reflect.TypeOf(v)
	if t.Kind() == reflect.Pointer {
		t = t.Elem()
	}
	if t.PkgPath() != "" && t.Name() != "" {
		return t.Name()
	}
	return t.String()
}
