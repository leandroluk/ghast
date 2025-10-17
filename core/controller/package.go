package controller

// New executes the builder function immediately
// and returns the finalized Controller instance.
//
// Example:
//
//	var SystemController = controller.New(func(b *controller.Builder) {
//		b.BasePath("/system")
//		b.Get("health", func(ctx middleware.Context) { ... })
//	})
func New(fn func(b *Builder)) *Controller {
	b := &Builder{}
	fn(b)
	return b.Build()
}
