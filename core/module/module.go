// core/module/module.go
package module

import (
	"fmt"
	"reflect"

	"github.com/leandroluk/ghast/core/container"
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
}

// Register registra recursivamente todos os providers do módulo
// (e módulos importados) no container e inicializa os singletons.
func (m *Module) Register(c *container.Container) {
	m.Container = c

	// Registra módulos importados primeiro
	for _, imported := range m.Imports {
		imported.Register(c)
	}

	// Registra os providers deste módulo
	for _, p := range m.Providers {
		c.Register(p)
	}

	// Resolve e inicializa imediatamente os providers Singleton
	for _, p := range m.Providers {
		if p.Scope == provider.Singleton {
			var instance any

			defer func() {
				if r := recover(); r != nil {
					fmt.Printf("[module] failed to resolve provider: %s (%v)\n", p.Name, r)
				}
			}()

			// Usa o tipo direto (sem ponteiro extra) — compatível com o registro
			providerType := reflect.Zero(p.Type).Interface()
			instance = c.Resolve(providerType)

			if instance == nil {
				fmt.Printf("[module] provider %s returned nil\n", p.Name)
				continue
			}

			// Executa hooks de ciclo de vida (redundância segura)
			if initable, ok := instance.(provider.OnInit); ok {
				initable.OnInit()
			}
			if afterInit, ok := instance.(provider.AfterInit); ok {
				afterInit.AfterInit()
			}

			fmt.Printf("[module] initialized provider: %s (%s)\n", p.Name, p.Scope)
		}
	}

	fmt.Printf(
		"[module] %s registered (%d providers, %d controllers)\n",
		m.Name,
		len(m.Providers),
		len(m.Controllers),
	)
}
