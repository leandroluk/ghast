// adapter/gin/adapter.go
package adapter

import (
	"github.com/gin-gonic/gin"
	"github.com/leandroluk/ghast/core/controller"
	"github.com/leandroluk/ghast/core/middleware"
)

// Adapter integrates Ghast with Gin.
type Adapter struct {
	engine *gin.Engine
}

// New creates a new Gin adapter instance.
func New(engine *gin.Engine) *Adapter {
	return &Adapter{engine: engine}
}

// OnRoute registers a Ghast controller route into Gin.
func (a *Adapter) OnRoute(path string, method controller.Method, handler middleware.Handler) {
	switch method {
	case controller.GET:
		a.engine.GET(path, func(c *gin.Context) { handler(NewContext(c)) })
	case controller.POST:
		a.engine.POST(path, func(c *gin.Context) { handler(NewContext(c)) })
	case controller.PUT:
		a.engine.PUT(path, func(c *gin.Context) { handler(NewContext(c)) })
	case controller.PATCH:
		a.engine.PATCH(path, func(c *gin.Context) { handler(NewContext(c)) })
	case controller.DELETE:
		a.engine.DELETE(path, func(c *gin.Context) { handler(NewContext(c)) })
	case controller.OPTIONS:
		a.engine.OPTIONS(path, func(c *gin.Context) { handler(NewContext(c)) })
	case controller.HEAD:
		a.engine.HEAD(path, func(c *gin.Context) { handler(NewContext(c)) })
	default:
		panic("[github.com/leandroluk/ghast/gin.Adapter] unsupported method: " + string(method))
	}
}

// Listen starts the Gin server.
func (a *Adapter) Listen(addr string) error {
	return a.engine.Run(addr)
}

func (a *Adapter) Start(addr string) error {
	return a.engine.Run(addr)
}
