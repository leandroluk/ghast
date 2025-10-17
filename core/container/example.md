# 🧩 Example — Container

Exemplo básico de uso do `core/container` para registrar providers e resolver dependências automaticamente.

---

### 📦 Estruturas simulando serviços

```go
package main

import (
	"fmt"

	"github.com/leandroluk/ghast/core/container"
	"github.com/leandroluk/ghast/core/provider"
)

// LoggerService é um serviço simples de logging.
type LoggerService struct{}

func (l *LoggerService) Log(msg string) {
	fmt.Println("[log]", msg)
}

// UserService depende do LoggerService.
type UserService struct {
	Logger *LoggerService `inject:"github.com/leandroluk/ghast/core/example.LoggerService"`
}

func (u *UserService) CreateUser(name string) {
	u.Logger.Log(fmt.Sprintf("Created user: %s", name))
}
