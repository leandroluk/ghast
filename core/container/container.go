// core/container/container.go
package container

import (
	"fmt"
	"reflect"
	"sync"
	"time"

	"github.com/leandroluk/ghast/core/provider"
)

// interface local para cleanup por request (não exige declarar no provider pkg)
type onRequestDestroy interface{ OnRequestDestroy() }

// Container manages dependency resolution and provider lifecycle.
type Container struct {
	mu        sync.RWMutex
	providers map[string]*provider.Provider
	instances map[string]any

	// request-scope
	reqMu   sync.RWMutex
	reqInst map[string]map[string]any // reqID -> providerName -> instance
}

// NewContainer creates a new instance of the container.
func NewContainer() *Container {
	return &Container{
		providers: make(map[string]*provider.Provider),
		instances: make(map[string]any),
		reqInst:   make(map[string]map[string]any),
	}
}

// Register adds a provider to the container.
// Panics if a provider with the same name already exists.
func (c *Container) Register(p *provider.Provider) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if _, exists := c.providers[p.Name]; exists {
		panic(fmt.Sprintf("[container] provider '%s' already registered", p.Name))
	}
	c.providers[p.Name] = p
}

// Resolve returns an instance of the requested dependency.
// It respects the provider's lifecycle scope (Singleton, Transient, Request).
// Para Request-scope sem reqID => panic (uso incorreto).
func (c *Container) Resolve(target any) any {
	return c.resolveInternal(target, "", make(map[string]bool))
}

// ResolveWithReq resolve considerando o cache por requisição (reqID).
func (c *Container) ResolveWithReq(target any, reqID string) any {
	return c.resolveInternal(target, reqID, make(map[string]bool))
}

// resolveInternal unifica a resolução e propaga reqID para TODA a cadeia de dependências.
func (c *Container) resolveInternal(target any, reqID string, seen map[string]bool) any {
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

	// ⚠️ sempre centraliza aqui
	return c.resolveProvider(p, reqID, seen)
}

// Inject preenche campos com tag `inject:"..."` em uma instância já criada.
// (Sem reqID: não suporta request-scoped aqui. Se tentar, vai panicar — correto.)
func (c *Container) Inject(target any) {
	v := reflect.ValueOf(target)
	if v.Kind() != reflect.Pointer || v.Elem().Kind() != reflect.Struct {
		panic("[container] Inject expects pointer to struct")
	}
	inst := v.Elem()
	seen := make(map[string]bool)

	for i := 0; i < inst.NumField(); i++ {
		field := inst.Type().Field(i)
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

		depInstance := c.resolveProvider(depProvider, "", seen)
		depVal := reflect.ValueOf(depInstance)

		if depVal.Type().AssignableTo(field.Type) && inst.Field(i).CanSet() {
			inst.Field(i).Set(depVal)
		} else {
			panic(fmt.Sprintf(
				"[container] cannot assign dependency %s to field %s (expected %s, got %s)",
				depName, field.Name, field.Type, depVal.Type(),
			))
		}
	}
}

// adiciona no arquivo (mesmo package)
func (c *Container) resolveProvider(p *provider.Provider, reqID string, seen map[string]bool) any {
	name := p.Name

	switch p.Scope {
	case provider.Request:
		if reqID == "" {
			panic(fmt.Sprintf("[container] provider %s is request-scoped but no request is active", p.Name))
		}

		// read path
		c.reqMu.RLock()
		if perReq, ok := c.reqInst[reqID]; ok {
			if inst, ok := perReq[name]; ok {
				c.reqMu.RUnlock()
				return inst
			}
		}
		c.reqMu.RUnlock()

		// cria
		instance := c.instantiate(p, seen, reqID)

		// write path (double-check)
		c.reqMu.Lock()
		if _, ok := c.reqInst[reqID]; !ok {
			c.reqInst[reqID] = make(map[string]any)
		}
		if existing, ok := c.reqInst[reqID][name]; ok {
			c.reqMu.Unlock()
			return existing
		}
		c.reqInst[reqID][name] = instance
		c.reqMu.Unlock()
		return instance

	case provider.Transient:
		return c.instantiate(p, seen, reqID)

	default: // Singleton
		// 1º read lock
		c.mu.RLock()
		if inst, ok := c.instances[name]; ok {
			c.mu.RUnlock()
			return inst
		}
		c.mu.RUnlock()

		// cria
		instance := c.instantiate(p, seen, reqID)

		// escreve com double-check
		c.mu.Lock()
		if existing, ok := c.instances[name]; ok {
			c.mu.Unlock()
			return existing
		}
		c.instances[name] = instance
		c.mu.Unlock()
		return instance
	}
}

// BeginRequest inicia o escopo de request para um reqID.
func (c *Container) BeginRequest(reqID string) {
	c.reqMu.Lock()
	if _, ok := c.reqInst[reqID]; !ok {
		c.reqInst[reqID] = make(map[string]any)
	}
	c.reqMu.Unlock()
}

// EndRequest finaliza o escopo: executa OnRequestDestroy() (se existir) e limpa o cache.
func (c *Container) EndRequest(reqID string) {
	c.reqMu.Lock()
	if perReq, ok := c.reqInst[reqID]; ok {
		for _, inst := range perReq {
			if d, ok := inst.(onRequestDestroy); ok {
				func() { defer func() { _ = recover() }(); d.OnRequestDestroy() }()
			}
		}
		delete(c.reqInst, reqID)
	}
	c.reqMu.Unlock()
}

// instantiate creates an instance of a provider and injects dependencies recursively.
func (c *Container) instantiate(p *provider.Provider, seen map[string]bool, reqID string) any {
	if p.Kind == provider.UseFactory {
		return c.instantiateFactory(p, reqID)
	}
	if p.Kind == provider.UseValue {
		return p.Value
	}
	return c.instantiateClass(p, seen, reqID)
}

// instantiateFactory handles both sync and async factories.
func (c *Container) instantiateFactory(p *provider.Provider, _ string) any {
	if p.FactoryC != nil {
		ch, err := p.FactoryC()
		if err != nil {
			panic(fmt.Sprintf("[container] async factory error for %s: %v", p.Name, err))
		}
		select {
		case instance := <-ch:
			c.callLifecycleHooks(instance)
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
		c.callLifecycleHooks(instance)
		return instance
	}

	panic(fmt.Sprintf("[container] factory provider %s has no factory function", p.Name))
}

// instantiateClass creates struct instances and injects dependencies recursively.
func (c *Container) instantiateClass(p *provider.Provider, seen map[string]bool, reqID string) any {
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

		// 🔴 AQUI é o pulo do gato:
		depInstance := c.resolveProvider(depProvider, reqID, seen)

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
	c.callLifecycleHooks(instance)
	return instance
}

// callLifecycleHooks executes OnInit and AfterInit hooks if implemented.
func (c *Container) callLifecycleHooks(instance any) {
	if hook, ok := instance.(provider.OnInit); ok {
		hook.OnInit()
	}
	if hook, ok := instance.(provider.AfterInit); ok {
		hook.AfterInit()
	}
}
