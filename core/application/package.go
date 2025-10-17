package application

import (
	"github.com/leandroluk/ghast/core/container"
	"github.com/leandroluk/ghast/core/module"
)

// New executes the builder immediately and returns a fully configured Application.
//
// Example:
//
//	var App = application.New(func(b *application.Builder) {
//	    b.Use(LoggerMiddleware)
//	    b.Intercept(MetricsInterceptor)
//	    b.Register(SystemModule)
//	})
func New(fn func(b *Builder)) *Application {
	b := &Builder{
		app: &Application{
			container: container.NewContainer(),
			modules:   []*module.Module{},
		},
	}
	fn(b)
	return b.app
}
