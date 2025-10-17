// core/controller/controller.go
package controller

import (
	"sync"

	"github.com/leandroluk/ghast/core/container"
)

type Controller struct {
	Name   string
	Base   string
	Routes []*Route

	// lazy build
	lazy  func(ctn *container.Container) *Controller
	once  sync.Once
	built bool
}

// EnsureBuilt resolve o controller quando o container existir.
func (c *Controller) EnsureBuilt(ctn *container.Container) *Controller {
	if c == nil || c.built || c.lazy == nil {
		return c
	}
	var built *Controller
	c.once.Do(func() {
		built = c.lazy(ctn)
		c.built = true
		c.lazy = nil
	})
	if built != nil {
		c.Name = built.Name
		c.Base = built.Base
		c.Routes = built.Routes
	}
	return c
}

// AddRoute permanece igual
func (c *Controller) AddRoute(r *Route) {
	c.Routes = append(c.Routes, r)
}
