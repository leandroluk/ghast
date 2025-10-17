package guard

import "github.com/leandroluk/ghast/core/middleware"

// Package guard defines the Guard interface used to control
// route-level authorization within the Ghast framework.
type Guard interface {
	// CanActivate is called before a controller handler executes.
	// It returns true to allow execution, or false to block it.
	CanActivate(ctx middleware.Context) bool
}

// Function is an adapter that allows a plain function to act as a Guard.
type Function func(ctx middleware.Context) bool

// CanActivate implements the Guard interface for Function.
func (fn Function) CanActivate(ctx middleware.Context) bool {
	return fn(ctx)
}
