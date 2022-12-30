package gee

import "net/http"

// HandlerFunc defines the request handler used by gee
type HandlerFunc func(*Context)

type Engine struct {
	router *routerGroup
	groups []*routerGroup
}

// New is the constructor of gee.Engine
func New() *Engine {
	e := &Engine{router: newRouterGroup()}
	e.router.Engine = e
	e.groups = []*routerGroup{e.router}
	return e
}

func (e *Engine) Group(prefix string) *routerGroup {
	newGroup := &routerGroup{
		router: e.router.router,
		prefix: e.router.prefix + prefix,
		parent: e.router,
		Engine: e,
	}
	e.groups = append(e.groups, newGroup)

	return newGroup
}

// GET defines the method to add GET request
func (e *Engine) GET(pattern string, handler HandlerFunc) {
	e.router.register("GET", pattern, handler)
}

// POST defines the method to add POST request
func (e *Engine) POST(pattern string, handler HandlerFunc) {
	e.router.register("POST", pattern, handler)
}

// Run defines the method to start a http server
func (engine *Engine) Run(addr string) (err error) {
	return http.ListenAndServe(addr, engine)
}

func (e *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	c := NewContext(w, req)
	e.router.handle(c)
}
