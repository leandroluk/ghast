# 🧩 Example — Provider

Exemplo de uso do sistema de *providers* do `core/provider`.  
O provider define **como um recurso é criado ou fornecido** dentro do container.

---

### 📦 Estruturas simulando serviços

```go
package main

import (
	"fmt"

	"github.com/leandroluk/ghast/core/provider"
)

// LoggerService é um serviço simples.
type LoggerService struct{}

func (l *LoggerService) Log(msg string) {
	fmt.Println("[log]", msg)
}

// Config é um valor fixo (UseValue)
type Config struct {
	Env string
}
```

---

### ⚙️ Criando providers

```go
func main() {
	// Provider baseado em classe (instanciado via reflexão)
	classProvider := provider.NewClass[LoggerService]()

	// Provider baseado em factory (instanciado via função)
	factoryProvider := provider.NewFactory(func() *LoggerService {
		return &LoggerService{}
	})

	// Provider baseado em valor fixo
	valueProvider := provider.NewValue(Config{Env: "production"})

	fmt.Println(classProvider.Name)
	fmt.Println(factoryProvider.Name)
	fmt.Println(valueProvider.Value)
}
```

---

### 🧠 Explicação rápida

| Tipo de Provider | Método           | Descrição |
|------------------|------------------|------------|
| `UseClass`       | `NewClass()`     | Cria instância via `reflect.New()` |
| `UseFactory`     | `NewFactory()`   | Cria instância via função customizada |
| `UseValue`       | `NewValue()`     | Fornece valor constante já instanciado |

---

### 🧩 Saída esperada

```
github.com/leandroluk/ghast/core/example.LoggerService
github.com/leandroluk/ghast/core/example.LoggerService
{production}
```

---

💡 **Dica:**  
O campo `Name` de cada provider é derivado automaticamente do tipo usando `provider.NameFromType()`.  
Isso garante nomes únicos e consistentes dentro do container.
