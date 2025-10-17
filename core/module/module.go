// core/module/module.go
package module

import (
	"fmt"
	"reflect"

	"github.com/leandroluk/ghast/core/container"
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

	// 1) registra imports primeiro
	for _, imported := range m.Imports {
		imported.Logger = m.Logger
		imported.Register(c)
	}

	// 2) registra providers
	for _, p := range m.Providers {
		c.Register(p)
	}

	// 3) resolve singletons (instancia) e chama OnModuleInit (se houver)
	for _, p := range m.Providers {
		if p.Scope != provider.Singleton {
			continue
		}

		defer func(name string) {
			if r := recover(); r != nil {
				fmt.Printf("[module] failed to resolve provider: %s (%v)\n", name, r)
			}
		}(p.Name)

		providerType := reflect.Zero(p.Type).Interface()
		inst := c.Resolve(providerType) // singleton cacheado

		if hook, ok := inst.(provider.OnModuleInit); ok {
			hook.OnModuleInit()
		}
	}

	// (controllers são “built” pelo Module.New/Controller.EnsureBuilt fora daqui)

	if m.Logger != nil {
		m.Logger.Logf("InstanceLoader", "%s dependencies initialized", m.Name)
	}
}
