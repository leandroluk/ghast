package module

import (
	"fmt"

	"github.com/leandroluk/ghast/core/container"
	"github.com/leandroluk/ghast/core/provider"
)

// Module represents a logical unit of the application.
// It groups related providers, controllers and imported modules.
type Module struct {
	Name        string
	Providers   []*provider.Provider
	Controllers []any
	Imports     []*Module
}

// Register recursively registers all module providers (including imported ones)
// into the dependency container. This method is called by the Application during bootstrap.
func (m *Module) Register(c *container.Container) {
	// Register imported modules first
	for _, imported := range m.Imports {
		imported.Register(c)
	}

	// Register this module's providers
	for _, p := range m.Providers {
		c.Register(p)
	}

	fmt.Printf(
		"[module] %s registered (%d providers, %d controllers)\n",
		m.Name,
		len(m.Providers),
		len(m.Controllers),
	)
}
