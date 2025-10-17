package application

import (
	"strings"

	"github.com/leandroluk/ghast/core/middleware"
)

func mapSlice[T any, R any](in []T, f func(T) R) []R {
	out := make([]R, 0, len(in))
	for _, v := range in {
		out = append(out, f(v))
	}
	return out
}

func mapAny(in []any, f func(any) string) []string {
	out := make([]string, 0, len(in))
	for _, v := range in {
		out = append(out, f(v))
	}
	return out
}

func joinTwo(a, b []string) string {
	parts := []string{}
	if len(a) > 0 {
		parts = append(parts, "global: "+strings.Join(a, ", "))
	}
	if len(b) > 0 {
		parts = append(parts, "route: "+strings.Join(b, ", "))
	}
	if len(parts) == 0 {
		return ""
	}
	return strings.Join(parts, " | ")
}

func safeMwName(m middleware.Middleware) string {
	if m.Name != "" {
		return m.Name
	}
	return "<anonymous>"
}
