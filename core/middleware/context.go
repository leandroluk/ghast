// core/middleware/context.go
package middleware

import (
	"io"
	"mime/multipart"
)

// Context representa o ciclo de vida de uma requisição HTTP no Ghast.
// Inspirado no modelo Express.js, separa entrada (Req) e saída (Res).
type Context interface {
	Next() error
	Abort()

	Req() Request
	Res() Response
	Locals() ContextLocals
}

// ======================================================================
// Request
// ======================================================================

// Request encapsula todos os dados de entrada da requisição HTTP.
type Request interface {
	Params() RequestParams
	Query() RequestQuery
	Headers() RequestHeaders
	Cookies() RequestCookies
	Body() RequestBody
}

// RequestParams representa os parâmetros de rota.
type RequestParams interface {
	Get(key string) string
	All() map[string]string
}

// RequestQuery representa os parâmetros de query string.
type RequestQuery interface {
	Get(key string) string
	All() map[string][]string
}

// RequestHeaders representa os cabeçalhos HTTP.
type RequestHeaders interface {
	Get(key string) string
	All() map[string]string
}

// RequestCookies representa os cookies recebidos.
type RequestCookies interface {
	Get(key string) string
	All() map[string]string
}

// RequestBody representa o corpo da requisição.
type RequestBody interface {
	Text() (string, error)
	Raw() ([]byte, error)
	Bind(v any) error
	Form() (map[string][]string, error)
	Multipart() (*multipart.Form, error)
	Stream() io.Reader
}

// ======================================================================
// Response
// ======================================================================

// Response encapsula a saída HTTP (status, headers, body, etc.).
type Response interface {
	Status(code int)
	Header(key, value string)
	Json(v any) error
	Text(body string) error
	Stream(r io.Reader) error
	Redirect(url string, code int) error
}

// ======================================================================
// Locals
// ======================================================================

// ContextLocals permite armazenar dados temporários durante o ciclo da requisição.
type ContextLocals interface {
	Set(key string, value any)
	Get(key string) any
	All() map[string]any
}
