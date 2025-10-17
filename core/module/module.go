// core/module/module.go
package module

import (
	"fmt"
	"reflect"

	"github.com/leandroluk/ghast/core/container"
	"github.com/leandroluk/ghast/core/controller"
	"github.com/leandroluk/ghast/core/logger"
	"github.com/leandroluk/ghast/core/provider"
)

// Module representa uma unidade lógica da aplicação.
// Agrupa providers, controllers e módulos importados.
type Module struct {
	Name        string
	Providers   []*provider.Provider
	Controllers []any
	Imports     []*Module
	Container   *container.Container
	Logger      logger.Logger
}

// Register registra recursivamente todos os providers do módulo
// (e módulos importados) no container e inicializa os singletons.
func (m *Module) Register(c *container.Container) {
	m.Container = c

	for _, imported := range m.Imports {
		// herda logger
		imported.Logger = m.Logger
		imported.Register(c)
	}

	for _, p := range m.Providers {
		c.Register(p)
	}

	for _, p := range m.Providers {
		if p.Scope == provider.Singleton {
			defer func() {
				if r := recover(); r != nil {
					fmt.Printf("[module] failed to resolve provider: %s (%v)\n", p.Name, r)
				}
			}()
			providerType := reflect.Zero(p.Type).Interface()
			_ = c.Resolve(providerType)
		}
	}

	for i, ctrl := range m.Controllers {
		if ctr, ok := ctrl.(*controller.Controller); ok {
			m.Controllers[i] = ctr.EnsureBuilt(c)
		}
	}

	if m.Logger != nil {
		m.Logger.Logf("InstanceLoader", "%s dependencies initialized", m.Name)
	}
}
