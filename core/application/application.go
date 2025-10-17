package application

import (
	"fmt"

	"github.com/leandroluk/ghast/core/container"
	"github.com/leandroluk/ghast/core/controller"
	"github.com/leandroluk/ghast/core/guard"
	"github.com/leandroluk/ghast/core/interceptor"
	"github.com/leandroluk/ghast/core/middleware"
	"github.com/leandroluk/ghast/core/module"
)

// Application represents the Ghast root runtime container.
// It coordinates modules, providers, controllers and global middlewares.
type Application struct {
	container          *container.Container
	modules            []*module.Module
	globalMiddlewares  []middleware.Middleware
	globalInterceptors []interceptor.Interceptor
	globalGuards       []guard.Guard
	adapters           []Adapter
}

// Builder provides a fluent interface for building an Application declaratively.
type Builder struct {
	app *Application
}

// Use registers global middlewares to be executed for every request.
func (b *Builder) Use(mw ...middleware.Middleware) *Builder {
	b.app.globalMiddlewares = append(b.app.globalMiddlewares, mw...)
	return b
}

// Intercept registers global interceptors applied to all controllers.
func (b *Builder) Intercept(ic ...interceptor.Interceptor) *Builder {
	b.app.globalInterceptors = append(b.app.globalInterceptors, ic...)
	return b
}

// Guard registers global guards applied before any controller handler.
func (b *Builder) Guard(g ...guard.Guard) *Builder {
	b.app.globalGuards = append(b.app.globalGuards, g...)
	return b
}

// Register adds one or more modules to the Application.
func (b *Builder) Register(mods ...*module.Module) *Builder {
	b.app.modules = append(b.app.modules, mods...)
	return b
}

// Mount attaches an adapter (HTTP, gRPC, etc.) to the Application.
func (b *Builder) Mount(adapters ...Adapter) *Builder {
	b.app.adapters = append(b.app.adapters, adapters...)
	return b
}

// Start bootstraps the application and runs all attached adapters.
func (a *Application) Start() error {
	for _, m := range a.modules {
		m.Register(a.container)
	}
	fmt.Printf("[app] registered %d modules\n", len(a.modules))

	for _, adapter := range a.adapters {
		if err := a.mountAdapter(adapter); err != nil {
			return err
		}
	}
	fmt.Println("[app] started successfully")
	return nil
}

// mountAdapter binds all controllers from modules into the given adapter.
func (a *Application) mountAdapter(adapter Adapter) error {
	for _, m := range a.modules {
		for _, ctrl := range m.Controllers {
			c, ok := ctrl.(*controller.Controller)
			if !ok {
				continue
			}
			for _, route := range c.Routes {
				adapter.OnRoute(c.Base, route.Path, route.Method, route.Handler)
			}
		}
	}
	return nil
}

// Container exposes the underlying DI container.
func (a *Application) Container() *container.Container {
	return a.container
}

// Snapshot returns a reflective view of the Application composition.
func (a *Application) Snapshot() *Snapshot {
	return snapshotOf(a)
}
