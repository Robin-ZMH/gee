package gee

import (
	"fmt"
	"net/http"
)

type HandlerFunc func(http.ResponseWriter, *http.Request)

type Engine struct {
	router map[string]HandlerFunc
}

// NewEngine is the constructor of gee.Engine
func NewEngine() *Engine {
	return &Engine{router: make(map[string]HandlerFunc)}
}

func (e *Engine) register(method, pattern string, handler HandlerFunc) {
	e.router[method+"-"+pattern] = handler
}

// GET defines the method to add GET request
func (e Engine) GET(pattern string, handler HandlerFunc) {
	e.register("GET", pattern, handler)
}

// POST defines the method to add POST request
func (engine *Engine) POST(pattern string, handler HandlerFunc) {
	engine.register("POST", pattern, handler)
}

// Run defines the method to start a http server
func (engine *Engine) Run(addr string) (err error) {
	return http.ListenAndServe(addr, engine)
}

func (e *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	key := req.Method + "-" + req.URL.Path
	if handler, ok := e.router[key]; ok {
		handler(w, req)
	} else {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "404 Not Found: %s\n", req.URL)
	}
}
