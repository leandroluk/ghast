// example/gin/main.go
package main

import (
	"github.com/gin-gonic/gin"
	adapter "github.com/leandroluk/ghast/adapter/gin"
	"github.com/leandroluk/ghast/core/application"
	"github.com/leandroluk/ghast/example"
)

func main() {
	ginApp := gin.New()
	application.New(func(b *application.Builder) {
		b.Register(example.UserModule)
		b.Mount(adapter.New(ginApp))
	}).Start(":3000")
}
