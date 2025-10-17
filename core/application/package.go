// core/application/package.go
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
	// cria o container global
	ctn := container.NewContainer()

	// cria o builder da aplicação com o container injetado
	b := &Builder{
		app: &Application{
			container: ctn,
			modules:   []*module.Module{},
		},
	}

	// executa o builder do usuário
	fn(b)

	return b.app
}
