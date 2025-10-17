// core/provider/package.go
package provider

import "reflect"

// New executa o builder e retorna o provider final.
func New(fn func(b *Builder)) *Provider {
	b := &Builder{}
	fn(b)
	return b.Build()
}

// NameFromType gera o nome público do provider com suporte a genéricos.
func NameFromType(t reflect.Type) string {
	return TypeName(t)
}

// TypeName retorna o nome totalmente qualificado de um tipo (compatível com genéricos).
func TypeName(t reflect.Type) string {
	if t.Kind() == reflect.Pointer {
		t = t.Elem()
	}

	// Go "normal": tem PkgPath + Name
	if t.Name() != "" && t.PkgPath() != "" {
		return t.PkgPath() + "." + t.Name()
	}

	// Go genérico: Name() vazio, String() retorna o formato correto
	// exemplo: "github.com/leandroluk/ghast/example.Repo[int]"
	return t.String()
}
