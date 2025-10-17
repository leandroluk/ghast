// core/application/package.go
package application

import (
	"github.com/leandroluk/ghast/core/container"
	"github.com/leandroluk/ghast/core/exception"
	"github.com/leandroluk/ghast/core/logger"
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
	ctn := container.NewContainer()
	b := &Builder{
		app: &Application{
			name:          "ghast",
			logger:        logger.NewConsole(),
			container:     ctn,
			modules:       []*module.Module{},
			defaultFilter: exception.NewDefault(),
		},
	}
	fn(b)
	return b.app
}
