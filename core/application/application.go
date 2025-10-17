// core/application/application.go  (alterar)
package application

import (
	"fmt"
	"strings"

	"github.com/leandroluk/ghast/core/container"
	"github.com/leandroluk/ghast/core/controller"
	"github.com/leandroluk/ghast/core/guard"
	"github.com/leandroluk/ghast/core/interceptor"
	"github.com/leandroluk/ghast/core/internal/naming"
	"github.com/leandroluk/ghast/core/logger"
	"github.com/leandroluk/ghast/core/middleware"
	"github.com/leandroluk/ghast/core/module"
)

type Application struct {
	name               string
	logger             logger.Logger
	container          *container.Container
	modules            []*module.Module
	globalMiddlewares  []middleware.Middleware
	globalInterceptors []interceptor.Interceptor
	globalGuards       []guard.Guard
	adapters           []Adapter
}

type Builder struct{ app *Application }

func (b *Builder) Name(name string) *Builder {
	b.app.name = name
	return b
}

func (b *Builder) Logger(l logger.Logger) *Builder {
	if l != nil {
		b.app.logger = l
	}
	return b
}

// já existentes...
func (b *Builder) Use(mw ...middleware.Middleware) *Builder {
	b.app.globalMiddlewares = append(b.app.globalMiddlewares, mw...)
	return b
}
func (b *Builder) Intercept(ic ...interceptor.Interceptor) *Builder {
	b.app.globalInterceptors = append(b.app.globalInterceptors, ic...)
	return b
}
func (b *Builder) Guard(g ...guard.Guard) *Builder {
	b.app.globalGuards = append(b.app.globalGuards, g...)
	return b
}
func (b *Builder) Register(mods ...*module.Module) *Builder {
	b.app.modules = append(b.app.modules, mods...)
	return b
}
func (b *Builder) Mount(adapters ...Adapter) *Builder {
	b.app.adapters = append(b.app.adapters, adapters...)
	return b
}

func (a *Application) Start(addr string) error {
	log := a.logger.WithAppTag(a.name)

	log.Logf("NestFactory", "Starting Ghast application...")

	for _, m := range a.modules {
		m.Logger = log
		m.Register(a.container)
	}
	log.Logf("NestFactory", "All modules registered (%d)", len(a.modules))

	for _, ad := range a.adapters {
		if err := a.mountAdapter(log, ad); err != nil {
			return err
		}
	}

	for _, ad := range a.adapters {
		if starter, ok := ad.(interface{ Start(addr string) error }); ok {
			if err := starter.Start(addr); err != nil {
				return fmt.Errorf("failed to start adapter: %w", err)
			}
		}
	}

	log.Logf("NestApplication", "Ghast application successfully started")
	log.Logf("bootstrap", "🌎 started on %s", addr)
	return nil
}

func (a *Application) mountAdapter(log logger.Logger, ad Adapter) error {
	// nomes dos globais (pra log bonito)
	globalGuardNames := mapSlice(a.globalGuards, func(g guard.Guard) string { return naming.TypeNameOf(g) })
	globalInterNames := mapSlice(a.globalInterceptors, func(i interceptor.Interceptor) string { return naming.TypeNameOf(i) })
	globalMwNames := mapSlice(a.globalMiddlewares, func(m middleware.Middleware) string { return safeMwName(m) })

	for _, m := range a.modules {
		for _, ctrl := range m.Controllers {
			c, ok := ctrl.(*controller.Controller)
			if !ok {
				continue
			}
			log.Logf("RoutesResolver", "%s {%s}:", c.Name, c.Base)

			for _, route := range c.Routes {
				// nomes específicos da rota
				routeGuardNames := mapAny(route.Guards, naming.TypeNameOf)
				routeInterNames := mapAny(route.Interceptors, naming.TypeNameOf)
				routeMwNames := make([]string, 0, len(route.Middlewares))
				for _, mw := range route.Middlewares {
					routeMwNames = append(routeMwNames, safeMwName(mw))
				}

				// print estilo Nest (extra)
				if len(globalGuardNames) > 0 || len(routeGuardNames) > 0 {
					log.Logf("RouterExplorer", "[Guards] %s",
						joinTwo(globalGuardNames, routeGuardNames))
				}
				if len(globalInterNames) > 0 || len(routeInterNames) > 0 {
					log.Logf("RouterExplorer", "[Interceptors] %s",
						joinTwo(globalInterNames, routeInterNames))
				}
				if len(globalMwNames) > 0 || len(routeMwNames) > 0 {
					log.Logf("RouterExplorer", "[Middlewares] %s",
						joinTwo(globalMwNames, routeMwNames))
				}

				log.Logf("RouterExplorer", "Mapped {%s, %s} route", route.Path, route.Method)
				ad.OnRoute(route.Path, route.Method, route.Handler)
			}
		}
	}

	return nil
}

func mapSlice[T any, R any](in []T, f func(T) R) []R {
	out := make([]R, 0, len(in))
	for _, v := range in {
		out = append(out, f(v))
	}
	return out
}

func mapAny(in []any, f func(any) string) []string {
	out := make([]string, 0, len(in))
	for _, v := range in {
		out = append(out, f(v))
	}
	return out
}

func joinTwo(a, b []string) string {
	parts := []string{}
	if len(a) > 0 {
		parts = append(parts, "global: "+strings.Join(a, ", "))
	}
	if len(b) > 0 {
		parts = append(parts, "route: "+strings.Join(b, ", "))
	}
	if len(parts) == 0 {
		return ""
	}
	return strings.Join(parts, " | ")
}

func safeMwName(m middleware.Middleware) string {
	if m.Name != "" {
		return m.Name
	}
	return "<anonymous>"
}
