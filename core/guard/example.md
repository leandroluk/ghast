# 🧩 Example — Guard

```go
package auth

import (
	"github.com/leandroluk/ghast/core/guard"
	"github.com/leandroluk/ghast/core/middleware"
)

type AuthGuard struct{}

func (AuthGuard) CanActivate(ctx middleware.Context) bool {
	token := ctx.Headers()["Authorization"]
	return token == "Bearer 123"
}

// Exemplo de uso
var AuthRequired = guard.Function(func(ctx middleware.Context) bool {
	return ctx.Headers()["X-Auth"] == "true"
})
```

E dentro de um controller futuramente:

```go
if !guard.Execute(ctx, AuthGuard{}, AuthRequired) {
	ctx.Status(403).Text("Forbidden")
	return
}
```