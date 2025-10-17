// core/module/define.go
package module

import (
	"reflect"
	"strings"

	"github.com/leandroluk/ghast/core/provider"
)

type DeclarativeModule interface {
	Configure(b *Builder)
}

// Opcional: permitir nome custom além do inferido
type NamedModule interface {
	ModuleName() string
}

// Define cria *Module a partir de um "módulo-classe".
func Define[T DeclarativeModule]() *Module {
	var zero T

	// Nome: do tipo T, limpando ponteiro e garantindo sufixo "Module"
	rt := reflect.TypeOf(zero)
	if rt.Kind() == reflect.Pointer {
		rt = rt.Elem()
	}
	name := rt.Name()
	if name == "" {
		name = "AppModule"
	}
	if n, ok := any(zero).(NamedModule); ok {
		name = n.ModuleName()
	} else if !strings.HasSuffix(name, "Module") {
		name += "Module"
	}

	b := &Builder{
		module: &Module{
			Name:        name,
			Providers:   []*provider.Provider{},
			Controllers: []any{},
			Imports:     []*Module{},
		},
	}
	zero.Configure(b)
	return b.Build()
}
