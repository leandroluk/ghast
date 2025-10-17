// example/fiber/main.go
package main

import (
	"github.com/gofiber/fiber/v2"
	adapter "github.com/leandroluk/ghast/adapter/fiber"
	"github.com/leandroluk/ghast/core/application"
	"github.com/leandroluk/ghast/example"
)

func main() {
	fiberApp := fiber.New()
	application.New(func(b *application.Builder) {
		b.Register(example.UserModule)
		b.Mount(adapter.New(fiberApp))
	}).Start(":3000")
}
