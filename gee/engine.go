package gee

import (
	"net/http"
)

// HandlerFunc defines the request handler used by gee
type HandlerFunc func(*Context)

type Engine struct {
	*routerGroup
	groups map[string]*routerGroup
}

// New is the constructor of gee.Engine
func New() *Engine {
	e := &Engine{routerGroup: newRouterGroup()}
	e.routerGroup.engine = e
	e.groups = map[string]*routerGroup{"": e.routerGroup}
	return e
}

// ServeHTTP can support middlewares
func (e *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	c := NewContext(e, w, req)
	e.router.handle(c)
}

// Run defines the method to start a http server
func (engine *Engine) Run(addr string) (err error) {
	return http.ListenAndServe(addr, engine)
}
