package exception

import "github.com/leandroluk/ghast/core/middleware"

// Filter captura erros e transforma em resposta HTTP.
type Filter interface {
	Catch(ctx middleware.Context, err error)
}

// DefaultFilter: HTTP 4xx/5xx padrão, sem vazar stack.
type DefaultFilter struct{}

func NewDefault() *DefaultFilter { return &DefaultFilter{} }

func (f *DefaultFilter) Catch(ctx middleware.Context, err error) {
	if err == nil {
		return
	}
	if he, ok := err.(*HttpError); ok {
		ctx.Res().Status(he.Status)
		_ = ctx.Res().Json(map[string]any{
			"statusCode": he.Status,
			"message":    he.Message,
		})
		return
	}
	ctx.Res().Status(500)
	_ = ctx.Res().Json(map[string]any{
		"statusCode": 500,
		"message":    "Internal Server Error",
	})
}
