// core/container/request.go
package container

import (
	"fmt"

	"github.com/leandroluk/ghast/core/middleware"
)

const reqIDKey = "__ghast_reqid"

func MustResolveRequest[T any](c *Container, ctx middleware.Context) *T {
	id, _ := ctx.Locals().Get(reqIDKey).(string)
	if id == "" {
		panic("[container] request scope not initialized")
	}
	var zero T
	raw := c.ResolveWithReq(&zero, id)
	inst, ok := raw.(*T)
	if !ok {
		name := fmt.Sprintf("%T", zero)
		panic(fmt.Sprintf("[container] incorrect type for %s", name))
	}
	return inst
}
