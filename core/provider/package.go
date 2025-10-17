package provider

import "reflect"

// New executes the builder function immediately and returns the final provider.
//
// Example:
//
//	var LoggerProvider = provider.New(func(b *provider.Builder) {
//		b.Factory("Logger", func() (any, error) { return log.Default(), nil })
//	})
func New(fn func(b *Builder)) *Provider {
	b := &Builder{}
	fn(b)
	return b.Build()
}

// NameFromType generates a readable provider name from a reflect.Type.
func NameFromType(t reflect.Type) string {
	return t.PkgPath() + "." + t.Name()
}
