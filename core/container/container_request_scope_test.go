package container_test

import (
	"sync"
	"sync/atomic"
	"testing"

	"github.com/leandroluk/ghast/core/container"
	"github.com/leandroluk/ghast/core/provider"
)

// =====================
// Tipos de teste
// =====================

// Request-scoped B
type ReqDepB struct {
	ID int32
}

var reqDepBGen int32

// OnInit é chamado pelo container após instanciar
func (b *ReqDepB) OnInit() {
	atomic.AddInt32(&reqDepBGen, 1)
	b.ID = reqDepBGen
}

// Request-scoped A -> depende de B
// ATENÇÃO: a tag deve bater com o nome do provider de B
// Nome esperado: "github.com/leandroluk/ghast/core/container.ReqDepB"
type ReqDepA struct {
	Dep *ReqDepB `inject:"github.com/leandroluk/ghast/core/container_test.ReqDepB"`
}

// Singleton que (erroneamente) depende de request -> deve panicar ao resolver sem reqID
type SingletonS struct {
	Dep *ReqDepB `inject:"github.com/leandroluk/ghast/core/container_test.ReqDepB"`
}

// Request-scoped com cleanup por request (OnRequestDestroy)
type CleanerReq struct{}

var cleanerCalls int32

// satisfaz a interface onRequestDestroy (método com essa assinatura)
func (c *CleanerReq) OnRequestDestroy() {
	atomic.AddInt32(&cleanerCalls, 1)
}

// =====================
// Helpers de registro
// =====================

func registerReqB(t *testing.T, c *container.Container) {
	t.Helper()
	p := provider.New(func(b *provider.Builder) {
		b.Class(ReqDepB{}).Scoped(provider.Request) // <- TEM que ser Request
	})
	if p.Scope != provider.Request {
		t.Fatalf("ReqDepB Provider scope incorreto: %s (esperado: request)", p.Scope)
	}
	c.Register(p)
}

func registerReqA(t *testing.T, c *container.Container) {
	t.Helper()
	p := provider.New(func(b *provider.Builder) {
		b.Class(ReqDepA{}).Scoped(provider.Request)
	})
	if p.Scope != provider.Request {
		t.Fatalf("ReqDepA Provider scope incorreto: %s (esperado: request)", p.Scope)
	}
	c.Register(p)
}

func registerSingletonS(t *testing.T, c *container.Container) {
	t.Helper()
	p := provider.New(func(b *provider.Builder) {
		b.Class(SingletonS{}) // default singleton
	})
	if p.Scope != provider.Singleton {
		t.Fatalf("SingletonS Provider scope incorreto: %s (esperado: singleton)", p.Scope)
	}
	c.Register(p)
}

func registerCleanerReq(c *container.Container) {
	c.Register(provider.New(func(b *provider.Builder) {
		b.Class(CleanerReq{}).Scoped(provider.Request)
	}))
}

// =====================
// Tests
// =====================

func TestRequestScope_DependencyGraphAndCache(t *testing.T) {
	c := container.NewContainer()
	registerReqB(t, c)
	registerReqA(t, c)

	// -------- r1 --------
	req1 := "r1"
	c.BeginRequest(req1)

	a1 := c.ResolveWithReq(&ReqDepA{}, req1).(*ReqDepA)
	a2 := c.ResolveWithReq(&ReqDepA{}, req1).(*ReqDepA)

	if a1 == nil || a2 == nil || a1.Dep == nil || a2.Dep == nil {
		t.Fatalf("nil instance(s): a1=%v a2=%v a1.Dep=%v a2.Dep=%v", a1, a2, a1.Dep, a2.Dep)
	}
	t.Logf("r1: a1=%p a2=%p | a1.Dep=%p(ID=%d) a2.Dep=%p(ID=%d)", a1, a2, a1.Dep, a1.Dep.ID, a2.Dep, a2.Dep.ID)

	// mesmo objeto dentro da mesma request
	if a1 != a2 {
		t.Fatalf("expected same ReqDepA instance within request %s", req1)
	}
	if a1.Dep != a2.Dep {
		t.Fatalf("expected same ReqDepB instance within request %s (via A.Dep)", req1)
	}

	c.EndRequest(req1)

	// -------- r2 --------
	req2 := "r2"
	c.BeginRequest(req2)

	a3 := c.ResolveWithReq(&ReqDepA{}, req2).(*ReqDepA)
	if a3 == nil || a3.Dep == nil {
		t.Fatalf("nil instance(s) in r2: a3=%v a3.Dep=%v", a3, a3.Dep)
	}
	t.Logf("r2: a3=%p | a3.Dep=%p(ID=%d)", a3, a3.Dep, a3.Dep.ID)

	// A deve ser diferente entre requests
	if a3 == a1 {
		t.Fatalf("expected different ReqDepA instance across requests (r1 vs r2)")
	}
	// B deve ser diferente entre requests — compara por ID (robusto p/ zero-sized)
	if a3.Dep.ID == a1.Dep.ID {
		t.Fatalf("expected different ReqDepB instance across requests (r1 vs r2): a1.Dep.ID=%d a3.Dep.ID=%d", a1.Dep.ID, a3.Dep.ID)
	}

	c.EndRequest(req2)
}

func TestRequestScope_PanicWhenSingletonDependsOnRequest(t *testing.T) {
	c := container.NewContainer()
	registerReqB(t, c)
	registerSingletonS(t, c)

	defer func() {
		if r := recover(); r == nil {
			t.Fatalf("expected panic when resolving Singleton depending on Request-scoped provider without request")
		}
	}()

	// Deve panicar: SingletonS -> ReqDepB (request) sem reqID
	_ = c.Resolve(&SingletonS{})
}

func TestRequestScope_ConcurrentSameInstance(t *testing.T) {
	c := container.NewContainer()
	registerReqB(t, c)

	req := "concurrent"
	c.BeginRequest(req)

	const N = 20
	var wg sync.WaitGroup
	wg.Add(N)

	ptrs := make(chan *ReqDepB, N)

	for i := 0; i < N; i++ {
		go func() {
			defer wg.Done()
			inst := c.ResolveWithReq(&ReqDepB{}, req).(*ReqDepB)
			ptrs <- inst
		}()
	}

	wg.Wait()
	close(ptrs)

	var first *ReqDepB
	for p := range ptrs {
		if first == nil {
			first = p
			continue
		}
		if p != first {
			t.Fatalf("expected all goroutines to receive the same ReqDepB instance for request %s", req)
		}
	}

	c.EndRequest(req)
}

func TestRequestScope_OnRequestDestroyCalled(t *testing.T) {
	atomic.StoreInt32(&cleanerCalls, 0)

	c := container.NewContainer()
	registerCleanerReq(c)

	req := "cleanup"
	c.BeginRequest(req)
	_ = c.ResolveWithReq(&CleanerReq{}, req).(*CleanerReq)
	c.EndRequest(req)

	if got := atomic.LoadInt32(&cleanerCalls); got != 1 {
		t.Fatalf("expected OnRequestDestroy to be called once, got %d", got)
	}
}
