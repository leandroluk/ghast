// core/provider/provider.go
package provider

import (
	"fmt"
	"reflect"
)

// Kind represents the creation strategy of a provider.
type Kind string

const (
	UseClass   Kind = "class"
	UseFactory Kind = "factory"
	UseValue   Kind = "value"
)

// Scope represents the lifecycle of a provider instance.
type Scope string

const (
	Singleton Scope = "singleton"
	Transient Scope = "transient"
	Request   Scope = "request"
)

// Provider defines a dependency that can be resolved by the container.
type Provider struct {
	Name     string
	Kind     Kind
	Scope    Scope
	Type     reflect.Type
	Value    any
	Factory  func() (any, error)        // synchronous or async-safe factory
	FactoryC func() (<-chan any, error) // optional fully async factory
}

// Lifecycle hooks
type OnInit interface {
	OnInit()
}

type AfterInit interface {
	AfterInit()
}

// Builder helps construct providers declaratively.
type Builder struct {
	provider *Provider
}

// Class sets this provider to use a struct type (UseClass).
func (b *Builder) Class(instanceType any) *Builder {
	t := reflect.TypeOf(instanceType)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	if t.Kind() != reflect.Struct {
		panic(fmt.Sprintf("[provider] class provider expects a struct type, got %s", t.Kind()))
	}

	b.provider = &Provider{
		Name:  NameFromType(t),
		Kind:  UseClass,
		Type:  t,
		Scope: Singleton,
	}
	return b
}

// Factory sets this provider to use a (possibly async) factory function.
func (b *Builder) Factory(name string, fn func() (any, error)) *Builder {
	b.provider = &Provider{
		Name:    name,
		Kind:    UseFactory,
		Scope:   Singleton,
		Factory: fn,
	}
	return b
}

// FactoryAsync defines an asynchronous provider factory.
func (b *Builder) FactoryAsync(name string, fn func() (<-chan any, error)) *Builder {
	b.provider = &Provider{
		Name:     name,
		Kind:     UseFactory,
		Scope:    Singleton,
		FactoryC: fn,
	}
	return b
}

// Value sets this provider to use a static value (UseValue).
func (b *Builder) Value(name string, val any) *Builder {
	b.provider = &Provider{
		Name:  name,
		Kind:  UseValue,
		Scope: Singleton,
		Value: val,
	}
	return b
}

// Scoped sets the lifecycle scope of this provider.
func (b *Builder) Scoped(scope Scope) *Builder {
	if b.provider == nil {
		panic("[provider] Scoped() called before defining provider type")
	}
	b.provider.Scope = scope
	return b
}

// Build finalizes and returns the constructed provider.
func (b *Builder) Build() *Provider {
	if b.provider == nil {
		panic("[provider] Build() called before defining provider type")
	}
	return b.provider
}
