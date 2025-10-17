package middleware

// Context represents a minimal interface for middleware execution.
// It should be implemented by the HTTP adapter (Fiber, Gin, etc.).
type Context interface {
	Next() error // Continue to the next middleware or handler
	Abort()      // Stop the chain immediately
	Set(key string, val any)
	Get(key string) any
}
