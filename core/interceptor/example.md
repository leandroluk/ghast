# 🧩 Example — Interceptor

```go
package metrics

import (
	"fmt"
	"time"

	"github.com/leandroluk/ghast/core/interceptor"
	"github.com/leandroluk/ghast/core/middleware"
)

var LoggingInterceptor = interceptor.Function(func(ctx middleware.Context, next interceptor.NextFunc) any {
	start := time.Now()
	result := next(ctx)
	fmt.Printf("[interceptor] took %v\n", time.Since(start))
	return result
})

var TransformInterceptor = interceptor.Function(func(ctx middleware.Context, next interceptor.NextFunc) any {
	data := next(ctx)
	return map[string]any{"data": data, "success": true}
})
```

E o pipeline de execução:

```go
result := interceptor.Execute(ctx, []interceptor.Interceptor{
	LoggingInterceptor,
	TransformInterceptor,
}, func(ctx middleware.Context) any {
	return "Hello, world!"
})
```

Resultado:

```json
{"data": "Hello, world!", "success": true}
```