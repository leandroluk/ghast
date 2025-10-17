package controller

// Controller represents a fully built controller with its routes.
// It is the final structure returned by the builder and registered in modules.
type Controller struct {
	Name   string
	Base   string
	Routes []*Route
}

// AddRoute appends a route directly to this controller.
// Mostly used internally or for dynamic route registration.
func (c *Controller) AddRoute(r *Route) {
	c.Routes = append(c.Routes, r)
}
