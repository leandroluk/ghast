package application

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"reflect"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/leandroluk/ghast/core/container"
	"github.com/leandroluk/ghast/core/controller"
	"github.com/leandroluk/ghast/core/exception"
	"github.com/leandroluk/ghast/core/guard"
	"github.com/leandroluk/ghast/core/interceptor"
	"github.com/leandroluk/ghast/core/internal/naming"
	"github.com/leandroluk/ghast/core/logger"
	"github.com/leandroluk/ghast/core/middleware"
	"github.com/leandroluk/ghast/core/module"
	"github.com/leandroluk/ghast/core/provider"
)

var reqSeq uint64

// Application representa o runtime root do Ghast.
type Application struct {
	name               string
	logger             logger.Logger
	container          *container.Container
	modules            []*module.Module
	globalMiddlewares  []middleware.Middleware
	globalInterceptors []interceptor.Interceptor
	globalGuards       []guard.Guard
	globalFilters      []exception.Filter
	adapters           []Adapter

	defaultFilter exception.Filter
}

// Adapter opcionalmente pode implementar Shutdown(ctx).
type stoppable interface {
	Shutdown(ctx context.Context) error
}

// Builder fornece a API fluente de configuração da app.
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

func (b *Builder) DefaultFilter(f exception.Filter) *Builder {
	b.app.defaultFilter = f
	return b
}

func (b *Builder) DisableDefaultFilter() *Builder {
	b.app.defaultFilter = nil
	return b
}

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

func (b *Builder) Filter(fs ...exception.Filter) *Builder {
	b.app.globalFilters = append(b.app.globalFilters, fs...)
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

// Start registra módulos, monta rotas, inicia adapters e gerencia shutdown.
func (a *Application) Start(addr string) error {
	log := a.logger.WithAppTag(a.name)
	log.Logf("NestFactory", "Starting Ghast application...")

	// registra módulos
	for _, m := range a.modules {
		m.Logger = log
		m.Register(a.container)
	}
	log.Logf("NestFactory", "All modules registered (%d)", len(a.modules))

	// mapeia rotas
	for _, ad := range a.adapters {
		if err := a.mountAdapter(log, ad); err != nil {
			return err
		}
	}

	// inicia adapters em goroutine
	errCh := make(chan error, len(a.adapters))
	for _, ad := range a.adapters {
		go func(ad Adapter) {
			if starter, ok := ad.(interface{ Start(addr string) error }); ok {
				errCh <- starter.Start(addr)
			} else {
				errCh <- nil
			}
		}(ad)
	}

	log.Logf("NestApplication", "Ghast application successfully started")
	log.Logf("bootstrap", "🌎 started on %s", addr)

	// espera sinal de SO ou erro de adapter
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)

	var shutReason string
	select {
	case sig := <-sigCh:
		shutReason = sig.String()
		log.Logf("NestApplication", "Shutting down (%s)...", shutReason)
	case err := <-errCh:
		if err != nil {
			shutReason = "adapter-error"
			log.Errorf("NestApplication", "Adapter error: %v", err)
		} else {
			shutReason = "adapter-exit"
			log.Logf("NestApplication", "Adapter exited")
		}
	}

	// tenta parar adapters com Shutdown(ctx)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	for _, ad := range a.adapters {
		if s, ok := ad.(stoppable); ok {
			_ = s.Shutdown(ctx)
		}
	}

	// hooks de shutdown nos providers singleton (sem instanciar novos)
	for _, m := range a.modules {
		for _, p := range m.Providers {
			if p.Scope != provider.Singleton {
				continue
			}
			// pega a instância cacheada (Resolve de singleton retorna a existente)
			inst := a.container.Resolve(reflect.Zero(p.Type).Interface())

			if appHook, ok := inst.(provider.OnApplicationShutdown); ok {
				// melhor esforço: passa ctx com timeout
				func() {
					defer func() { _ = recover() }()
					appHook.OnApplicationShutdown(ctx)
				}()
			}
			if destroy, ok := inst.(provider.OnModuleDestroy); ok {
				func() {
					defer func() { _ = recover() }()
					destroy.OnModuleDestroy()
				}()
			}
		}
	}

	log.Logf("NestApplication", "Shutdown complete (%s)", shutReason)
	return nil
}

// mountAdapter aplica o pipeline real e mapeia rotas no adapter.
func (a *Application) mountAdapter(log logger.Logger, ad Adapter) error {
	globalGuardNames := mapSlice(a.globalGuards, func(g guard.Guard) string { return naming.TypeNameOf(g) })
	globalInterNames := mapSlice(a.globalInterceptors, func(i interceptor.Interceptor) string { return naming.TypeNameOf(i) })
	globalMwNames := mapSlice(a.globalMiddlewares, func(m middleware.Middleware) string { return safeMwName(m) })
	globalFilterNames := mapSlice(a.globalFilters, func(f exception.Filter) string { return naming.TypeNameOf(f) })

	for _, m := range a.modules {
		for _, ctrl := range m.Controllers {
			c, ok := ctrl.(*controller.Controller)
			if !ok {
				continue
			}
			log.Logf("RoutesResolver", "%s {%s}:", c.Name, c.Base)

			for _, rt := range c.Routes {
				// nomes por rota (log)
				routeGuardNames := mapAny(rt.Guards, naming.TypeNameOf)
				routeInterNames := mapAny(rt.Interceptors, naming.TypeNameOf)
				routeMwNames := make([]string, 0, len(rt.Middlewares))
				for _, mw := range rt.Middlewares {
					routeMwNames = append(routeMwNames, safeMwName(mw))
				}
				routeFilterNames := mapAny(rt.Filters, naming.TypeNameOf)

				if len(globalGuardNames) > 0 || len(routeGuardNames) > 0 {
					log.Logf("RouterExplorer", "[Guards] %s", joinTwo(globalGuardNames, routeGuardNames))
				}
				if len(globalInterNames) > 0 || len(routeInterNames) > 0 {
					log.Logf("RouterExplorer", "[Interceptors] %s", joinTwo(globalInterNames, routeInterNames))
				}
				if len(globalMwNames) > 0 || len(routeMwNames) > 0 {
					log.Logf("RouterExplorer", "[Middlewares] %s", joinTwo(globalMwNames, routeMwNames))
				}
				if len(globalFilterNames) > 0 || len(routeFilterNames) > 0 {
					log.Logf("RouterExplorer", "[Filters] %s", joinTwo(globalFilterNames, routeFilterNames))
				}

				// === PIPELINE REAL ===
				h := rt.Handler

				// Interceptors
				allInterceptors := append([]interceptor.Interceptor{}, a.globalInterceptors...)
				for _, it := range rt.Interceptors {
					if v, ok := it.(interceptor.Interceptor); ok {
						allInterceptors = append(allInterceptors, v)
					}
				}
				h = wrapWithInterceptors(h, allInterceptors)

				// Guards
				allGuards := append([]guard.Guard{}, a.globalGuards...)
				for _, g := range rt.Guards {
					if v, ok := g.(guard.Guard); ok {
						allGuards = append(allGuards, v)
					}
				}
				h = wrapWithGuards(h, allGuards)

				// Middlewares
				allMws := append([]middleware.Middleware{}, a.globalMiddlewares...)
				allMws = append(allMws, rt.Middlewares...)
				h = middleware.Compose(h, allMws...)

				// Exceptions
				allFilters := append([]exception.Filter{}, a.globalFilters...)
				for _, f := range rt.Filters {
					if v, ok := f.(exception.Filter); ok {
						allFilters = append(allFilters, v)
					}
				}
				if a.defaultFilter != nil {
					allFilters = append(allFilters, a.defaultFilter)
				} else {
					allFilters = append(allFilters, exception.NewDefault())
				}
				h = wrapWithExceptions(h, allFilters)

				// Request scope (mais externo)
				h = a.wrapWithRequestScope(h)

				// Mapear
				log.Logf("RouterExplorer", "Mapped {%s, %s} route", rt.Path, rt.Method)
				ad.OnRoute(rt.Path, rt.Method, h)
			}
		}
	}
	return nil
}

// wrapWithRequestScope abre um escopo de request no container e limpa no fim.
func (a *Application) wrapWithRequestScope(final middleware.Handler) middleware.Handler {
	return func(ctx middleware.Context) error {
		id := fmt.Sprintf("%d", atomic.AddUint64(&reqSeq, 1))
		a.container.BeginRequest(id)
		ctx.Locals().Set("__ghast_reqid", id)
		defer a.container.EndRequest(id)
		return final(ctx)
	}
}
