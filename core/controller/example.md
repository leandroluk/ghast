# Controller Example

``` go
package example

import (
    "fmt"
    "github.com/leandroluk/ghast/core/controller"
    "github.com/leandroluk/ghast/core/middleware"
)

var SystemController = controller.NewController(func(b *controller.Builder) {
    b.BasePath("/system")

    b.Get("health", func(ctx middleware.Context) {
        fmt.Println("health ok")
    })

    b.Get("info", func(ctx middleware.Context) {
        fmt.Println("system info")
    }).Use(loggingMW)
})
```

-   A função passada para `NewController` é executada imediatamente.
-   O `Builder` mantém estado interno e gera automaticamente um
    `*controller.Type`.
-   Rota final: `/system/health`, método `GET`.
