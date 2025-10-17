// adapter/fiber/adapter.go
package adapter

import (
	"github.com/gofiber/fiber/v2"
	"github.com/leandroluk/ghast/core/controller"
	"github.com/leandroluk/ghast/core/middleware"
)

type Adapter struct {
	app *fiber.App
}

func New(app *fiber.App) *Adapter {
	return &Adapter{app: app}
}

func (a *Adapter) OnRoute(path string, method controller.Method, handler middleware.Handler) {
	fiberHandler := func(c *fiber.Ctx) error {
		ctx := NewContext(c)
		return handler(ctx)
	}

	switch method {
	case controller.GET:
		a.app.Get(path, fiberHandler)
	case controller.POST:
		a.app.Post(path, fiberHandler)
	case controller.PUT:
		a.app.Put(path, fiberHandler)
	case controller.PATCH:
		a.app.Patch(path, fiberHandler)
	case controller.DELETE:
		a.app.Delete(path, fiberHandler)
	case controller.OPTIONS:
		a.app.Options(path, fiberHandler)
	case controller.HEAD:
		a.app.Head(path, fiberHandler)
	default:
		panic("[github.com/leandroluk/ghast/fiber.Adapter] unsupported method: " + string(method))
	}
}

func (a *Adapter) Start(addr string) error {
	return a.app.Listen(addr)
}
