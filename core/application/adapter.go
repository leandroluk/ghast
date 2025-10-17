package application

import (
	"github.com/leandroluk/ghast/core/controller"
	"github.com/leandroluk/ghast/core/middleware"
)

// Adapter defines a transport integration (HTTP, gRPC, etc.)
// responsible for binding routes from controllers into its routing system.
type Adapter interface {
	OnRoute(base string, path string, method controller.Method, handler middleware.Handler)
}
