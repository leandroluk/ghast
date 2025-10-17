package interceptor

import "github.com/leandroluk/ghast/core/middleware"

// Execute runs a chain of interceptors sequentially.
// Each interceptor calls the next until the final handler executes.
func Execute(ctx middleware.Context, interceptors []Interceptor, final NextFunc) any {
	if len(interceptors) == 0 {
		return final(ctx)
	}

	// Recursive chain builder
	chain := func(current middleware.Context) any {
		if len(interceptors) == 0 {
			return final(current)
		}
		head := interceptors[0]
		tail := interceptors[1:]
		return head.Intercept(current, func(nextCtx middleware.Context) any {
			return Execute(nextCtx, tail, final)
		})
	}

	return chain(ctx)
}

// New creates a new Interceptor from a function.
func New(fn func(ctx middleware.Context, next NextFunc) any) Interceptor {
	return Function(fn)
}
