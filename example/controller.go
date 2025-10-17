// example/controller.go
package example

import (
	"fmt"
	"strconv"

	"github.com/leandroluk/ghast/core/controller"
	"github.com/leandroluk/ghast/core/middleware"
)

// UserController lida com todas as rotas relacionadas a usuários.
type UserController struct {
	Service *UserService `inject:"github.com/leandroluk/ghast/example/example.UserService"`
}

// ListUser GET /users
func (c *UserController) ListUser(ctx middleware.Context) error {
	users := c.Service.FindAll()
	fmt.Printf("GET /users - returned %d users\n", len(users))
	return ctx.Res().Json(users)
}

// GetUser GET /users/:id
func (c *UserController) GetUser(ctx middleware.Context) error {
	idStr := ctx.Req().Params().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		fmt.Println("GET /users/:id - invalid id")
		ctx.Res().Status(400)
		return ctx.Res().Text("invalid id")
	}

	user := c.Service.FindByID(id)
	if user == nil {
		fmt.Printf("GET /users/%d - not found\n", id)
		ctx.Res().Status(404)
		return ctx.Res().Text("user not found")
	}

	fmt.Printf("GET /users/%d - found %+v\n", id, user)
	return ctx.Res().Json(user)
}

// CreateUser POST /users
func (c *UserController) CreateUser(ctx middleware.Context) error {
	var body struct {
		Name  string `json:"name"`
		Email string `json:"email"`
	}

	if err := ctx.Req().Body().Bind(&body); err != nil {
		ctx.Res().Status(400)
		return ctx.Res().Text("invalid request body")
	}

	user := c.Service.Create(body.Name, body.Email)
	fmt.Printf("POST /users - created %+v\n", user)

	ctx.Res().Status(201)
	return ctx.Res().Json(user)
}

// UpdateUser PUT /users/:id
func (c *UserController) UpdateUser(ctx middleware.Context) error {
	idStr := ctx.Req().Params().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		ctx.Res().Status(400)
		return ctx.Res().Text("invalid id")
	}

	var body struct {
		Name  string `json:"name"`
		Email string `json:"email"`
	}
	if err := ctx.Req().Body().Bind(&body); err != nil {
		ctx.Res().Status(400)
		return ctx.Res().Text("invalid body")
	}

	user := c.Service.Update(id, body.Name, body.Email)
	if user == nil {
		ctx.Res().Status(404)
		return ctx.Res().Text("user not found")
	}

	fmt.Printf("PUT /users/%d - updated %+v\n", id, user)
	return ctx.Res().Json(user)
}

// DeleteUser DELETE /users/:id
func (c *UserController) DeleteUser(ctx middleware.Context) error {
	idStr := ctx.Req().Params().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		ctx.Res().Status(400)
		return ctx.Res().Text("invalid id")
	}

	deleted := c.Service.Delete(id)
	if deleted {
		fmt.Printf("DELETE /users/%d - deleted\n", id)
		ctx.Res().Status(204)
		return nil
	}

	fmt.Printf("DELETE /users/%d - not found\n", id)
	ctx.Res().Status(404)
	return ctx.Res().Text("user not found")
}

// Controller registra todas as rotas de usuário no Ghast.
var Controller = controller.New(func(b *controller.Builder, c *UserController) {
	b.BasePath("/users")
	b.Get("", c.ListUser)
	b.Get("/:id", c.GetUser)
	b.Post("", c.CreateUser)
	b.Put("/:id", c.UpdateUser)
	b.Delete("/:id", c.DeleteUser)
})
