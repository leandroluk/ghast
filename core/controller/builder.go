package controller

import (
	"fmt"
	"reflect"

	"github.com/leandroluk/ghast/core/middleware"
)

// Builder provides a fluent API for defining controllers and their routes.
// Each method modifies the builder state and returns itself for chaining.
type Builder struct {
	name   string
	base   string
	routes []*Route
}

// BasePath sets the base path for all routes defined in this controller.
// The leading slash is optional.
func (b *Builder) BasePath(path string) *Builder {
	if path == "" {
		b.base = ""
		return b
	}
	if path[0] != '/' {
		path = "/" + path
	}
	b.base = path
	return b
}

// Get registers a GET route.
func (b *Builder) Get(path string, handler middleware.Handler) *Builder {
	return b.add(GET, path, handler)
}

// Post registers a POST route.
func (b *Builder) Post(path string, handler middleware.Handler) *Builder {
	return b.add(POST, path, handler)
}

// Put registers a PUT route.
func (b *Builder) Put(path string, handler middleware.Handler) *Builder {
	return b.add(PUT, path, handler)
}

// Patch registers a PATCH route.
func (b *Builder) Patch(path string, handler middleware.Handler) *Builder {
	return b.add(PATCH, path, handler)
}

// Delete registers a DELETE route.
func (b *Builder) Delete(path string, handler middleware.Handler) *Builder {
	return b.add(DELETE, path, handler)
}

// Options registers a OPTIONS route.
func (b *Builder) Options(path string, handler middleware.Handler) *Builder {
	return b.add(OPTIONS, path, handler)
}

// Head registers a HEAD route.
func (b *Builder) Head(path string, handler middleware.Handler) *Builder {
	return b.add(HEAD, path, handler)
}

// Trace registers a TRACE route.
func (b *Builder) Trace(path string, handler middleware.Handler) *Builder {
	return b.add(TRACE, path, handler)
}

// Connect registers a CONNECT route.
func (b *Builder) Connect(path string, handler middleware.Handler) *Builder {
	return b.add(CONNECT, path, handler)
}

// Use attaches middleware(s) to the most recently defined route.
func (b *Builder) Use(mw ...middleware.Middleware) *Builder {
	if len(b.routes) == 0 {
		panic("controller.Use() called before any route definition")
	}
	last := b.routes[len(b.routes)-1]
	last.Middlewares = append(last.Middlewares, mw...)
	return b
}

// Guard attaches guard(s) to the most recently defined route.
func (b *Builder) Guard(gs ...any) *Builder {
	if len(b.routes) == 0 {
		panic("controller.Guard() called before any route definition")
	}
	last := b.routes[len(b.routes)-1]
	last.Guards = append(last.Guards, gs...)
	return b
}

// Intercept attaches interceptor(s) to the most recently defined route.
func (b *Builder) Intercept(is ...any) *Builder {
	if len(b.routes) == 0 {
		panic("controller.Intercept() called before any route definition")
	}
	last := b.routes[len(b.routes)-1]
	last.Interceptors = append(last.Interceptors, is...)
	return b
}

// Build finalizes the controller and returns its instance.
// It infers the name automatically when not provided.
func (b *Builder) Build() *Controller {
	name := b.name
	if name == "" {
		name = inferName()
	}
	return &Controller{
		Name:   name,
		Base:   b.base,
		Routes: b.routes,
	}
}

// add registers a new route under the given method and path.
func (b *Builder) add(method Method, path string, handler middleware.Handler) *Builder {
	if path == "" {
		path = "/"
	}
	if path[0] != '/' {
		path = "/" + path
	}
	full := b.base + path
	b.routes = append(b.routes, &Route{
		Method:  method,
		Path:    full,
		Handler: handler,
	})
	return b
}

// inferName provides a generic fallback name when none is defined.
// Used mainly for anonymous controllers.
func inferName() string {
	pc := reflect.TypeOf(Builder{})
	return fmt.Sprintf("Controller_%v", pc.NumMethod())
}
