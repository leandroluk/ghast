package middleware

import "fmt"

// Chain composes multiple middlewares into a single handler.
// It executes them sequentially using the Context's Next() mechanism.
func Chain(middlewares ...*Middleware) Handler {
	return func(ctx Context) error {
		index := -1

		var next func() error
		next = func() error {
			index++
			if index >= len(middlewares) {
				return nil
			}

			mw := middlewares[index]
			return mw.Handler(&chainedContext{
				Context: ctx,
				next:    next,
			})
		}

		return next()
	}
}

// chainedContext wraps a Context to control middleware execution order.
type chainedContext struct {
	Context
	next func() error
}

// New creates a middleware instance with a given name and handler.
// Equivalent to NewMiddleware(name, fn).
func New(name string, fn Handler) *Middleware {
	if fn == nil {
		panic(fmt.Sprintf("[middleware] '%s' handler cannot be nil", name))
	}
	return &Middleware{
		Name:    name,
		Handler: fn,
	}
}
