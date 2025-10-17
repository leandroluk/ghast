// core/middleware/middleware.go
package middleware

// Package middleware defines the minimal middleware interface
// and chaining mechanism for the Ghast framework.
type Middleware struct {
	Name    string
	Handler Handler
}
