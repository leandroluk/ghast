// core/module/package.go
package module

import "github.com/leandroluk/ghast/core/provider"

// New executes the builder function immediately and
// returns the finalized Module instance.
//
// Example:
//
//	var SystemModule = module.New(func(b *module.Builder) {
//	    b.AddControllers(SystemController)
//	    b.AddProviders(DbProvider)
//	})
func New(fn func(b *Builder)) *Module {
	b := &Builder{
		module: &Module{
			Providers:   []*provider.Provider{},
			Controllers: []any{},
			Imports:     []*Module{},
		},
	}
	fn(b)
	return b.Build()
}
