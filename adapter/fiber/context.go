// adapter/fiber/context.go
package adapter

import (
	"io"
	"mime/multipart"

	"github.com/gofiber/fiber/v2"
	"github.com/leandroluk/ghast/core/middleware"
)

// Context adapta o *fiber.Ctx para a interface middleware.Context do Ghast.
type Context struct {
	c       *fiber.Ctx
	aborted bool
}

// NewContext cria um novo wrapper Context a partir do *fiber.Ctx.
func NewContext(c *fiber.Ctx) middleware.Context {
	return &Context{c: c}
}

// Next continua para o próximo middleware/handler.
func (ctx *Context) Next() error {
	if ctx.aborted {
		return nil
	}
	return ctx.c.Next()
}

// Abort interrompe imediatamente a cadeia de handlers/middlewares.
func (ctx *Context) Abort() {
	ctx.aborted = true
	// opcional: ctx.c.Status(fiber.StatusOK)
}

// ==== Implementação da interface middleware.Context ====

func (ctx *Context) Req() middleware.Request  { return &fiberRequest{ctx.c} }
func (ctx *Context) Res() middleware.Response { return &fiberResponse{ctx.c} }
func (ctx *Context) Locals() middleware.ContextLocals {
	return &fiberLocals{ctx.c}
}

// ======================================================================
// Request
// ======================================================================

type fiberRequest struct{ c *fiber.Ctx }

func (r *fiberRequest) Params() middleware.RequestParams   { return &fiberParams{r.c} }
func (r *fiberRequest) Query() middleware.RequestQuery     { return &fiberQuery{r.c} }
func (r *fiberRequest) Headers() middleware.RequestHeaders { return &fiberHeaders{r.c} }
func (r *fiberRequest) Cookies() middleware.RequestCookies { return &fiberCookies{r.c} }
func (r *fiberRequest) Body() middleware.RequestBody       { return &fiberBody{r.c} }

// --- Params ---

type fiberParams struct{ c *fiber.Ctx }

func (p *fiberParams) Get(key string) string {
	return p.c.Params(key)
}

func (p *fiberParams) All() map[string]string {
	m := make(map[string]string)
	for _, name := range p.c.Route().Params {
		m[name] = p.c.Params(name)
	}
	return m
}

// --- Query ---

type fiberQuery struct{ c *fiber.Ctx }

func (q *fiberQuery) Get(key string) string {
	return q.c.Query(key)
}

func (q *fiberQuery) All() map[string][]string {
	args := q.c.Context().QueryArgs()
	m := make(map[string][]string)
	args.VisitAll(func(k, v []byte) {
		m[string(k)] = append(m[string(k)], string(v))
	})
	return m
}

// --- Headers ---

type fiberHeaders struct{ c *fiber.Ctx }

func (h *fiberHeaders) Get(key string) string {
	return h.c.Get(key)
}

func (h *fiberHeaders) All() map[string]string {
	m := make(map[string]string)
	h.c.Request().Header.VisitAll(func(k, v []byte) {
		m[string(k)] = string(v)
	})
	return m
}

// --- Cookies ---

type fiberCookies struct{ c *fiber.Ctx }

func (ck *fiberCookies) Get(key string) string {
	return ck.c.Cookies(key)
}

func (ck *fiberCookies) All() map[string]string {
	m := make(map[string]string)
	ck.c.Request().Header.VisitAllCookie(func(k, v []byte) {
		m[string(k)] = string(v)
	})
	return m
}

// --- Body ---

type fiberBody struct{ c *fiber.Ctx }

func (b *fiberBody) Text() (string, error) {
	return string(b.c.Body()), nil
}

func (b *fiberBody) Raw() ([]byte, error) {
	return b.c.Body(), nil
}

func (b *fiberBody) Bind(v any) error {
	return b.c.BodyParser(v)
}

func (b *fiberBody) Form() (map[string][]string, error) {
	form, err := b.c.MultipartForm()
	if err != nil && err != fiber.ErrUnprocessableEntity {
		return nil, err
	}
	if form == nil {
		return nil, nil
	}
	return form.Value, nil
}

func (b *fiberBody) Multipart() (*multipart.Form, error) {
	return b.c.MultipartForm()
}

func (b *fiberBody) Stream() io.Reader {
	return b.c.Context().RequestBodyStream()
}

// ======================================================================
// Response
// ======================================================================

type fiberResponse struct{ c *fiber.Ctx }

func (r *fiberResponse) Status(code int) {
	r.c.Status(code)
}

func (r *fiberResponse) Header(key, value string) {
	r.c.Set(key, value)
}

func (r *fiberResponse) Json(v any) error {
	return r.c.JSON(v)
}

func (r *fiberResponse) Text(body string) error {
	return r.c.SendString(body)
}

func (r *fiberResponse) Stream(reader io.Reader) error {
	_, err := io.Copy(r.c, reader)
	return err
}

func (r *fiberResponse) Redirect(url string, code int) error {
	return r.c.Redirect(url, code)
}

// ======================================================================
// Locals
// ======================================================================

type fiberLocals struct{ c *fiber.Ctx }

func (l *fiberLocals) Set(key string, value any) {
	l.c.Locals(key, value)
}

func (l *fiberLocals) Get(key string) any {
	return l.c.Locals(key)
}

func (l *fiberLocals) All() map[string]any {
	// Fiber não expõe todos os Locals publicamente
	return map[string]any{}
}
