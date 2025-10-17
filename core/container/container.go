// core/container/container.go
package container

import (
	"fmt"
	"reflect"
	"sync"
	"time"

	"github.com/leandroluk/ghast/core/provider"
)

// Container manages dependency resolution and provider lifecycle.
type Container struct {
	mu        sync.RWMutex
	providers map[string]*provider.Provider
	instances map[string]any
}

// NewContainer creates a new instance of the container.
func NewContainer() *Container {
	return &Container{
		providers: make(map[string]*provider.Provider),
		instances: make(map[string]any),
	}
}

// Register adds a provider to the container.
// Panics if a provider with the same name already exists.
func (c *Container) Register(p *provider.Provider) {
	c.mu.Lock()
	defer c.mu.Unlock()

	fmt.Println("[container] registering:", p.Name)

	if _, exists := c.providers[p.Name]; exists {
		panic(fmt.Sprintf("[container] provider '%s' already registered", p.Name))
	}

	c.providers[p.Name] = p
}

// Resolve returns an instance of the requested dependency.
// It respects the provider's lifecycle scope (Singleton, Transient, Request).
func (c *Container) Resolve(target any) any {
	targetType := reflect.TypeOf(target)
	if targetType.Kind() == reflect.Pointer {
		targetType = targetType.Elem()
	}

	name := provider.TypeName(targetType)

	c.mu.RLock()
	p, ok := c.providers[name]
	c.mu.RUnlock()
	if !ok {
		panic(fmt.Sprintf("[container] provider not found: %s", name))
	}

	switch p.Scope {
	case provider.Transient:
		return c.instantiate(p, make(map[string]bool))

	case provider.Request:
		fmt.Printf("[container] warning: request scope not yet implemented for %s\n", p.Name)
		return c.instantiate(p, make(map[string]bool))

	default: // Singleton
		c.mu.RLock()
		instance, exists := c.instances[name]
		c.mu.RUnlock()
		if exists {
			return instance
		}

		instance = c.instantiate(p, make(map[string]bool))

		c.mu.Lock()
		c.instances[name] = instance
		c.mu.Unlock()

		return instance
	}
}

// instantiate creates an instance of a provider and injects dependencies recursively.
func (c *Container) instantiate(p *provider.Provider, seen map[string]bool) any {
	if p.Kind == provider.UseFactory {
		return c.instantiateFactory(p)
	}
	if p.Kind == provider.UseValue {
		return p.Value
	}
	return c.instantiateClass(p, seen)
}

// instantiateFactory handles both sync and async factories.
func (c *Container) instantiateFactory(p *provider.Provider) any {
	if p.FactoryC != nil {
		ch, err := p.FactoryC()
		if err != nil {
			panic(fmt.Sprintf("[container] async factory error for %s: %v", p.Name, err))
		}
		select {
		case instance := <-ch:
			c.callLifecycleHooks(p, instance)
			return instance
		case <-time.After(10 * time.Second):
			panic(fmt.Sprintf("[container] async factory timeout for %s", p.Name))
		}
	}

	if p.Factory != nil {
		instance, err := p.Factory()
		if err != nil {
			panic(fmt.Sprintf("[container] factory error for %s: %v", p.Name, err))
		}
		c.callLifecycleHooks(p, instance)
		return instance
	}

	panic(fmt.Sprintf("[container] factory provider %s has no factory function", p.Name))
}

// instantiateClass creates struct instances and injects dependencies recursively.
func (c *Container) instantiateClass(p *provider.Provider, seen map[string]bool) any {
	if seen[p.Name] {
		panic(fmt.Sprintf("[container] circular dependency detected: %s", p.Name))
	}
	seen[p.Name] = true

	providerType := p.Type
	instanceValue := reflect.New(providerType).Elem()

	for i := 0; i < instanceValue.NumField(); i++ {
		field := instanceValue.Type().Field(i)
		tag := field.Tag.Get("inject")
		if tag == "" {
			continue
		}

		depName := tag
		if depName == "" {
			depName = field.Type.String()
		}

		c.mu.RLock()
		depProvider, ok := c.providers[depName]
		c.mu.RUnlock()
		if !ok {
			panic(fmt.Sprintf("[container] dependency not found: %s", depName))
		}

		depInstance := c.instantiate(depProvider, seen)
		depVal := reflect.ValueOf(depInstance)

		if depVal.Type().AssignableTo(field.Type) && instanceValue.Field(i).CanSet() {
			instanceValue.Field(i).Set(depVal)
		} else {
			panic(fmt.Sprintf(
				"[container] cannot assign dependency %s to field %s (expected %s, got %s)",
				depName, field.Name, field.Type, depVal.Type(),
			))
		}
	}

	instance := instanceValue.Addr().Interface()
	c.callLifecycleHooks(p, instance)
	return instance
}

// callLifecycleHooks executes OnInit and AfterInit hooks if implemented.
func (c *Container) callLifecycleHooks(p *provider.Provider, instance any) {
	if hook, ok := instance.(provider.OnInit); ok {
		hook.OnInit()
	}
	if hook, ok := instance.(provider.AfterInit); ok {
		hook.AfterInit()
	}
}
