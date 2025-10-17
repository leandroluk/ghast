package module

import "github.com/leandroluk/ghast/core/provider"

// Builder provides a fluent API for building modules declaratively.
type Builder struct {
	module *Module
}

// AddProviders registers providers for this module.
func (b *Builder) AddProviders(providers ...*provider.Provider) *Builder {
	b.module.Providers = append(b.module.Providers, providers...)
	return b
}

// AddControllers registers controllers for this module.
func (b *Builder) AddControllers(controllers ...any) *Builder {
	b.module.Controllers = append(b.module.Controllers, controllers...)
	return b
}

// Import adds other modules as dependencies.
func (b *Builder) Import(modules ...*Module) *Builder {
	b.module.Imports = append(b.module.Imports, modules...)
	return b
}

// Build finalizes and returns the module instance.
func (b *Builder) Build() *Module {
	return b.module
}
