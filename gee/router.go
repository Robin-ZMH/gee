package gee

import (
	"log"
	"net/http"
	"strings"
)

type router struct {
	tries    map[string]*trie
	handlers map[string]HandlerFunc
}

func newRouter() *router {
	return &router{
		tries:    make(map[string]*trie),
		handlers: make(map[string]HandlerFunc),
	}
}

func (r *router) register(method string, pattern string, handler HandlerFunc) {
	log.Printf("ADD route %4s - %s", method, pattern)

	if _, ok := r.tries[method]; !ok {
		r.tries[method] = newTrie()
	}
	r.tries[method].insert(pattern)

	key := method + "-" + pattern
	r.handlers[key] = handler
}

func (r *router) matchPattern(method, path string) (pattern string) {
	if t, ok := r.tries[method]; ok {
		node := t.search(path)
		if node != nil {
			return node.pattern
		}
	}
	return
}

func (r *router) getParams(path, pattern string) map[string]string {
	params := make(map[string]string)
	path_parts := parsePattern(path)
	for i, part := range parsePattern(pattern) {
		switch part[0] {
		case ':':
			params[part[1:]] = path_parts[i]
		case '*':
			params[part[1:]] = strings.Join(path_parts[i:], "/")
			return params
		}
	}
	return params
}

func (r *router) handle(c *Context) {
	pattern := r.matchPattern(c.Method, c.Path)
	if pattern != "" {
		params := r.getParams(c.Path, pattern)
		c.Params = params
		key := c.Method + "-" + pattern
		r.handlers[key](c)
	} else {
		c.String(http.StatusNotFound, "404 NOT FOUND: %s\n", c.Path)
	}
}
