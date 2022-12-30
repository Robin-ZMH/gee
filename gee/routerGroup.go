package gee

type routerGroup struct {
	// embedding router
	*router
	// new fields
	prefix string
	parent *routerGroup
	engine *Engine
}

func newRouterGroup() *routerGroup {
	return &routerGroup{router: newRouter()}
}

func (r *routerGroup) addRoute(method, pattern string, handler HandlerFunc) {
	pattern = r.prefix + pattern
	r.register(method, pattern, handler)
}

func (r *routerGroup) Group(prefix string) *routerGroup {
	newGroup := &routerGroup{
		router: r.router,
		prefix: r.prefix + prefix,
		parent: r,
		engine: r.engine,
	}
	r.engine.groups = append(r.engine.groups, newGroup)

	return newGroup
}

// GET defines the method to add GET request
func (g *routerGroup) GET(pattern string, handler HandlerFunc) {
	g.addRoute("GET", pattern, handler)
}

// POST defines the method to add POST request
func (g *routerGroup) POST(pattern string, handler HandlerFunc) {
	g.addRoute("POST", pattern, handler)
}
