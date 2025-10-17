// core/controller/router.go
package controller

import "github.com/leandroluk/ghast/core/middleware"

// Method defines the supported HTTP methods for controllers.
type Method string

const (
	GET     Method = "GET"
	POST    Method = "POST"
	PUT     Method = "PUT"
	PATCH   Method = "PATCH"
	DELETE  Method = "DELETE"
	OPTIONS Method = "OPTIONS"
	HEAD    Method = "HEAD"
)

// Route represents a single controller route definition,
// including its method, path, handler, guards, interceptors, and middlewares.
type Route struct {
	Method       Method
	Path         string
	Handler      middleware.Handler
	Guards       []any
	Interceptors []any
	Middlewares  []middleware.Middleware
	Filters      []any
}
