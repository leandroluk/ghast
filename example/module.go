// example/module.go
package example

import (
	"github.com/leandroluk/ghast/core/module"
	"github.com/leandroluk/ghast/core/provider"
)

var UserModule = module.New(func(b *module.Builder) {
	b.AddProviders(
		provider.New(func(p *provider.Builder) { p.Class(UserService{}) }),
	)
	b.AddControllers(Controller)
})
