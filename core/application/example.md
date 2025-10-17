# Ghast Application Example

``` go
package main

import (
    "fmt"
    "github.com/leandroluk/ghast/core/application"
    "github.com/leandroluk/ghast/core/controller"
    "github.com/leandroluk/ghast/core/middleware"
    "github.com/leandroluk/ghast/core/module"
)

// MockAdapter is a minimal implementation of the application.Adapter interface.
// It prints all routes being mounted by the Application.
type MockAdapter struct{}

func (m *MockAdapter) OnRoute(method controller.Method, path string, handler middleware.Handler) {
    fmt.Printf("[adapter] route registered: %s %s\n", method, path)
}

func main() {
    // Create a new Ghast application instance
    app := application.New()

    // Define a simple controller
    systemController := controller.New(func(b *controller.Builder) {
        b.BasePath("/system")
        b.Get("health", func(ctx middleware.Context) error {
            fmt.Println("health ok")
            return nil
        })
    })

    // Create a module containing the controller
    systemModule := module.NewBuilder("System").
        Controllers(systemController).
        Build()

    // Register the module in the app
    app.Register(systemModule)

    // Add a global middleware (example)
    mw := middleware.NewMiddleware("logger", func(ctx middleware.Context) error {
        fmt.Println("[middleware] before handler")
        err := ctx.Next()
        fmt.Println("[middleware] after handler")
        return err
    })
    app.Use(mw)

    // Mount all routes into a mock adapter (or Fiber/Gin in real usage)
    adapter := &MockAdapter{}
    app.Mount(adapter)
}
```

### Explanation

-   `application.New()` creates the main Ghast application.
-   Controllers are defined using `controller.NewController()`.
-   Modules group controllers and providers together.
-   `Application.Register()` wires the modules and their dependencies.
-   `Application.Mount()` walks all routes and binds them to an adapter
    (agnostic to HTTP framework).
-   You can implement `application.Adapter` for any framework (Fiber,
    Gin, Echo, etc.).

### Key Points

-   The `Application` manages dependency injection, guards,
    interceptors, and middleware execution.
-   No coupling to any HTTP library --- `Adapter` is the only
    abstraction layer.
-   Ideal for frameworks or custom server implementations.
