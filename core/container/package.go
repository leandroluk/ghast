package container

import (
	"fmt"
	"reflect"
)

// MustResolve automatically resolves and casts the requested generic type.
// If the type is not found or the cast fails, it causes a panic.
func MustResolve[T any](c *Container) *T {
	var zero T
	name := reflect.TypeOf(zero).String()
	raw := c.Resolve(&zero)
	instance, ok := raw.(*T)
	if !ok {
		panic(fmt.Sprintf("[container] incorrect type for %s", name))
	}
	return instance
}
