package interceptor

import (
	"context"

	"github.com/leandroluk/ghast/core/middleware"
)

// NextFunc defines the next function in the interceptor chain.
type NextFunc func(ctx middleware.Context) any

// Package interceptor provides a mechanism for wrapping and extending controller execution.
// Interceptors can transform input/output or add cross-cutting behavior.
type Interceptor interface {
	// Intercept receives the request context and the next handler in the chain.
	// It can run code before and/or after calling `next(ctx)`.
	Intercept(ctx middleware.Context, next NextFunc) any
}

// Function allows plain functions to act as interceptors.
type Function func(ctx middleware.Context, next NextFunc) any

// Intercept implements the Interceptor interface for Function.
func (fn Function) Intercept(ctx middleware.Context, next NextFunc) any {
	return fn(ctx, next)
}

// ContextKey is used for metadata injection in the interceptor chain.
type ContextKey string

// WithValue safely injects metadata into the context.
func WithValue(ctx context.Context, key ContextKey, value any) context.Context {
	return context.WithValue(ctx, key, value)
}

// Value retrieves metadata from the context.
func Value(ctx context.Context, key ContextKey) any {
	return ctx.Value(key)
}
