// adapter/gin/adapter.go
package adapter

import (
	"context"  // <- novo
	"net/http" // <- novo

	"github.com/gin-gonic/gin"
	"github.com/leandroluk/ghast/core/controller"
	"github.com/leandroluk/ghast/core/middleware"
)

type Adapter struct {
	engine *gin.Engine
	srv    *http.Server // <- para shutdown gracioso
}

func New(engine *gin.Engine) *Adapter { return &Adapter{engine: engine} }

func (a *Adapter) OnRoute(path string, method controller.Method, handler middleware.Handler) {
	wrap := func(c *gin.Context) { _ = handler(NewContext(c)) }
	switch method {
	case controller.GET:
		a.engine.GET(path, wrap)
	case controller.POST:
		a.engine.POST(path, wrap)
	case controller.PUT:
		a.engine.PUT(path, wrap)
	case controller.PATCH:
		a.engine.PATCH(path, wrap)
	case controller.DELETE:
		a.engine.DELETE(path, wrap)
	case controller.OPTIONS:
		a.engine.OPTIONS(path, wrap)
	case controller.HEAD:
		a.engine.HEAD(path, wrap)
	default:
		panic("[github.com/leandroluk/ghast/gin.Adapter] unsupported method: " + string(method))
	}
}

// Start com http.Server para permitir Shutdown(ctx).
func (a *Adapter) Start(addr string) error {
	a.srv = &http.Server{
		Addr:    addr,
		Handler: a.engine,
	}
	// ListenAndServe bloqueia até Shutdown(ctx) ser chamado
	err := a.srv.ListenAndServe()
	// Quando Shutdown é chamado, o retorno costuma ser http.ErrServerClosed — não é “erro fatal”.
	if err == http.ErrServerClosed {
		return nil
	}
	return err
}

// Opcional: manter por compatibilidade, mas não é usado pelo core:
// func (a *Adapter) Listen(addr string) error { return a.Start(addr) }

// Implementa a interface stoppable do Application.
func (a *Adapter) Shutdown(ctx context.Context) error {
	if a.srv != nil {
		return a.srv.Shutdown(ctx)
	}
	return nil
}
