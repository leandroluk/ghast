# 🧱 Example — Module

Exemplo de uso do sistema de *módulos* do `core/module`.  
Módulos agrupam **providers**, **controllers** e **imports** — como no NestJS.

---

### 📦 Estruturas simulando serviços

```go
package main

import (
	"fmt"

	"github.com/leandroluk/ghast/core/container"
	"github.com/leandroluk/ghast/core/module"
	"github.com/leandroluk/ghast/core/provider"
)

type LoggerService struct{}

func (l *LoggerService) Log(msg string) {
	fmt.Println("[log]", msg)
}

type UserService struct {
	Logger *LoggerService `inject:"github.com/leandroluk/ghast/core/example.LoggerService"`
}

func (u *UserService) CreateUser(name string) {
	u.Logger.Log(fmt.Sprintf("Created user: %s", name))
}
```

---

### ⚙️ Criando módulos

```go
func main() {
	loggerProvider := provider.NewClass[LoggerService]()
	userProvider := provider.NewClass[UserService]()

	// Cria o módulo de usuários
	userModule := module.NewBuilder("UserModule").
		AddProviders(loggerProvider, userProvider).
		Build()

	// Registra no container
	c := container.New()
	userModule.Register(c)

	// Resolve o serviço
	userService := c.Resolve(&UserService{}).(*UserService)
	userService.CreateUser("Alice")
}
```

---

### 🧩 Saída esperada

```
[module] UserModule registered with 2 providers and 0 controllers
[log] Created user: Alice
```

---

### 🧠 Explicação rápida

| Conceito | Descrição |
|-----------|------------|
| `Module` | Agrupa providers, controllers e imports |
| `ModuleBuilder` | Builder fluente (`AddProviders`, `AddControllers`, `Import`) |
| `Register()` | Registra todos os providers (recursivamente) no container |
| `Container` | Gerencia instâncias e injeções automaticamente |

---

💡 **Dica:**  
Módulos podem importar outros módulos (`Import()`) para compartilhar dependências, exatamente como no NestJS.
