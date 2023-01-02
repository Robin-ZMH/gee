package gee

type routerGroup struct {
	router      *router
	prefix      string
	parent      *routerGroup
	engine      *Engine
	middlewares []HandlerFunc
}

func newRouterGroup() *routerGroup {
	return &routerGroup{router: newRouter()}
}

func (g *routerGroup) scanMiddleWares() (middlewares []HandlerFunc) {
	node := g
	for node != nil {
		middlewares = append(node.middlewares, middlewares...)
		node = node.parent
	}
	return
}

func (g *routerGroup) addRoute(method, pattern string, handler HandlerFunc) {
	pattern = g.prefix + pattern
	middlewares := g.scanMiddleWares()
	middlewares = append(middlewares, handler)
	g.router.register(method, pattern, middlewares...)
}

// GET defines the method to add GET request
func (g *routerGroup) GET(pattern string, handler HandlerFunc) {
	g.addRoute("GET", pattern, handler)
}

// POST defines the method to add POST request
func (g *routerGroup) POST(pattern string, handler HandlerFunc) {
	g.addRoute("POST", pattern, handler)
}

// Group create a new group that inherit from the parent group
func (g *routerGroup) Group(prefix string) *routerGroup {
	newGroup := &routerGroup{
		router: g.router,
		prefix: g.prefix + prefix,
		parent: g,
		engine: g.engine,
	}
	g.engine.groups[prefix] = newGroup

	return newGroup
}

// Use adds the middlewares to the current level group
func (g *routerGroup) Use(middlewares ...HandlerFunc) {
	g.middlewares = append(g.middlewares, middlewares...)
}
