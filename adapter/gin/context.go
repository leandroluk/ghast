// adapter/gin/context.go
package adapter

import (
	"io"
	"mime/multipart"

	"github.com/gin-gonic/gin"
	"github.com/leandroluk/ghast/core/middleware"
)

// Context implementa middleware.Context para o framework Gin.
type Context struct {
	c       *gin.Context
	aborted bool
}

// NewContext cria um novo wrapper Context a partir do *gin.Context.
func NewContext(c *gin.Context) middleware.Context {
	return &Context{c: c}
}

// Next continua para o próximo middleware/handler.
func (ctx *Context) Next() error {
	if ctx.aborted {
		return nil
	}
	ctx.c.Next()
	return nil
}

// Abort interrompe imediatamente a cadeia de handlers/middlewares.
func (ctx *Context) Abort() {
	ctx.aborted = true
	ctx.c.Abort()
}

// ==== Implementação da interface middleware.Context ====

func (ctx *Context) Req() middleware.Request  { return &ginRequest{ctx.c} }
func (ctx *Context) Res() middleware.Response { return &ginResponse{ctx.c} }
func (ctx *Context) Locals() middleware.ContextLocals {
	return &ginLocals{ctx.c}
}

// ======================================================================
// Request
// ======================================================================

type ginRequest struct{ c *gin.Context }

func (r *ginRequest) Params() middleware.RequestParams   { return &ginParams{r.c} }
func (r *ginRequest) Query() middleware.RequestQuery     { return &ginQuery{r.c} }
func (r *ginRequest) Headers() middleware.RequestHeaders { return &ginHeaders{r.c} }
func (r *ginRequest) Cookies() middleware.RequestCookies { return &ginCookies{r.c} }
func (r *ginRequest) Body() middleware.RequestBody       { return &ginBody{r.c} }

// --- Params ---

type ginParams struct{ c *gin.Context }

func (p *ginParams) Get(key string) string {
	return p.c.Param(key)
}

func (p *ginParams) All() map[string]string {
	m := make(map[string]string)
	for _, param := range p.c.Params {
		m[param.Key] = param.Value
	}
	return m
}

// --- Query ---

type ginQuery struct{ c *gin.Context }

func (q *ginQuery) Get(key string) string {
	return q.c.Query(key)
}

func (q *ginQuery) All() map[string][]string {
	return q.c.Request.URL.Query()
}

// --- Headers ---

type ginHeaders struct{ c *gin.Context }

func (h *ginHeaders) Get(key string) string {
	return h.c.GetHeader(key)
}

func (h *ginHeaders) All() map[string]string {
	m := make(map[string]string)
	for k, v := range h.c.Request.Header {
		if len(v) > 0 {
			m[k] = v[0]
		}
	}
	return m
}

// --- Cookies ---

type ginCookies struct{ c *gin.Context }

func (ck *ginCookies) Get(key string) string {
	val, err := ck.c.Cookie(key)
	if err != nil {
		return ""
	}
	return val
}

func (ck *ginCookies) All() map[string]string {
	m := make(map[string]string)
	for _, cookie := range ck.c.Request.Cookies() {
		m[cookie.Name] = cookie.Value
	}
	return m
}

// --- Body ---

type ginBody struct{ c *gin.Context }

func (b *ginBody) Text() (string, error) {
	data, err := io.ReadAll(b.c.Request.Body)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func (b *ginBody) Raw() ([]byte, error) {
	data, err := io.ReadAll(b.c.Request.Body)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (b *ginBody) Bind(v any) error {
	return b.c.ShouldBindJSON(v)
}

func (b *ginBody) Form() (map[string][]string, error) {
	if err := b.c.Request.ParseForm(); err != nil {
		return nil, err
	}
	return b.c.Request.Form, nil
}

func (b *ginBody) Multipart() (*multipart.Form, error) {
	if err := b.c.Request.ParseMultipartForm(32 << 20); err != nil { // 32MB default
		return nil, err
	}
	return b.c.Request.MultipartForm, nil
}

func (b *ginBody) Stream() io.Reader {
	return b.c.Request.Body
}

// ======================================================================
// Response
// ======================================================================

type ginResponse struct{ c *gin.Context }

func (r *ginResponse) Status(code int) {
	r.c.Status(code)
}

func (r *ginResponse) Header(key, value string) {
	r.c.Header(key, value)
}

func (r *ginResponse) Json(v any) error {
	r.c.JSON(r.c.Writer.Status(), v)
	return nil
}

func (r *ginResponse) Text(body string) error {
	r.c.String(r.c.Writer.Status(), body)
	return nil
}

func (r *ginResponse) Stream(reader io.Reader) error {
	_, err := io.Copy(r.c.Writer, reader)
	return err
}

func (r *ginResponse) Redirect(url string, code int) error {
	r.c.Redirect(code, url)
	return nil
}

// ======================================================================
// Locals
// ======================================================================

type ginLocals struct{ c *gin.Context }

func (l *ginLocals) Set(key string, value any) {
	l.c.Set(key, value)
}

func (l *ginLocals) Get(key string) any {
	v, _ := l.c.Get(key)
	return v
}

func (l *ginLocals) All() map[string]any {
	// Gin não expõe todos os locals internamente
	return map[string]any{}
}
