package application

import (
	"github.com/leandroluk/ghast/core/exception"
	"github.com/leandroluk/ghast/core/guard"
	"github.com/leandroluk/ghast/core/interceptor"
	"github.com/leandroluk/ghast/core/middleware"
)

func wrapWithGuards(final middleware.Handler, gs []guard.Guard) middleware.Handler {
	return func(ctx middleware.Context) error {
		for _, g := range gs {
			if !g.CanActivate(ctx) {
				// 403 por padrão
				ctx.Res().Status(403)
				return ctx.Res().Json(map[string]any{
					"statusCode": 403,
					"message":    "Forbidden",
				})
			}
		}
		return final(ctx)
	}
}

func wrapWithInterceptors(final middleware.Handler, is []interceptor.Interceptor) middleware.Handler {
	return func(ctx middleware.Context) error {
		out := interceptor.Execute(ctx, is, func(nextCtx middleware.Context) any {
			return final(nextCtx) // retorna error (ou nil)
		})
		if err, ok := out.(error); ok {
			return err
		}
		return nil
	}
}

func wrapWithExceptions(final middleware.Handler, fs []exception.Filter) middleware.Handler {
	return func(ctx middleware.Context) (err error) {
		defer func() {
			if r := recover(); r != nil {
				err = exception.FromRecover(r)
			}
			if err != nil {
				for _, f := range fs {
					f.Catch(ctx, err)
				}
				// erro já tratado pelos filtros
				err = nil
			}
		}()
		err = final(ctx)
		return err
	}
}
