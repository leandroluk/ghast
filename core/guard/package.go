// core/guard/package.go
package guard

import "github.com/leandroluk/ghast/core/middleware"

// Execute runs a sequence of guards in order.
// Returns false if any guard fails (CanActivate == false).
func Execute(ctx middleware.Context, guards ...Guard) bool {
	for _, g := range guards {
		if !g.CanActivate(ctx) {
			return false
		}
	}
	return true
}

// New creates a Guard from a function, equivalent to guard.Function(fn).
func New(fn func(ctx middleware.Context) bool) Guard {
	return Function(fn)
}
