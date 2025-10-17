// core/application/request.go
package application

import (
	"github.com/leandroluk/ghast/core/container"
	"github.com/leandroluk/ghast/core/middleware"
)

func ResolveRequest[T any](app *Application, ctx middleware.Context) *T {
	return container.MustResolveRequest[T](app.container, ctx)
}
