# 🧩 Example — Middleware

Exemplo de uso do sistema de *middlewares* genérico do `core/middleware`.  
O modelo é **agnóstico de framework**, podendo ser adaptado para Fiber, Gin, Echo ou outro.

---

### 📦 Estruturas básicas

```go
package main

import (
	"fmt"

	"github.com/leandroluk/ghast/core/middleware"
)

// Context de exemplo simples
type SimpleContext struct {
	Path string
}

func (c *SimpleContext) Next() {
	fmt.Println("[context] moving to next middleware")
}
