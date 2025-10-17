// core/controller/package.go
package controller

import (
	"fmt"
	"reflect"

	"github.com/leandroluk/ghast/core/container"
)

// New executa o builder function e garante que o controller tenha receivers corretos.
// Exemplo:
//
//	var SystemController = controller.New(func(b *controller.Builder, c *MyController) {
//		b.BasePath("/system")
//		b.Get("health", c.HealthCheck)
//	})
func New[T any](fn func(b *Builder, c *T)) *Controller {
	t := reflect.TypeOf((*T)(nil)).Elem()
	validateControllerType(t)

	// devolve um Controller "vazio" com closure de build
	return &Controller{
		lazy: func(ctn *container.Container) *Controller {
			b := &Builder{}
			b.name = t.Name()
			instance := new(T)
			ctn.Inject(instance)
			fn(b, instance)
			return b.Build()
		},
	}
}

// validateControllerType garante que todos os métodos do controller
// tenham receiver por ponteiro (*Controller), não por valor.
func validateControllerType(t reflect.Type) {
	ptrType := reflect.PointerTo(t)
	for i := 0; i < t.NumMethod(); i++ {
		m := t.Method(i)
		if _, ok := ptrType.MethodByName(m.Name); !ok {
			panic(fmt.Sprintf(
				"[controller] method %s must have a pointer receiver (*%s)",
				m.Name, t.Name(),
			))
		}
	}
}
