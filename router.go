package gee

import (
	"log"
	"net/http"
	"strings"
)

type router struct {
	tries map[string]*trie
}

func newRouter() *router {
	return &router{tries: make(map[string]*trie)}
}

func (r *router) register(method string, pattern string, middlewares ...HandlerFunc) {
	log.Printf("ADD route %4s - %s\n", method, pattern)

	if _, ok := r.tries[method]; !ok {
		r.tries[method] = newTrie()
	}
	r.tries[method].insert(pattern, middlewares)
}

func (r *router) route(method, path string) *node {
	if t, ok := r.tries[method]; ok {
		node := t.search(path)
		return node
	}
	return nil
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

// handle will parse the url path to find the handler function,
// and add it into *(Context).handlers, then call *(Context).Next()
func (r *router) handle(c *Context) {
	node := r.route(c.Method, c.Path)
	if node != nil {
		pattern, handlers := node.pattern, node.handlers
		params := r.getParams(c.Path, pattern)
		c.handlers = append(c.handlers, handlers...)
		c.Params = params
	} else {
		c.handlers = append(c.handlers, func(ctx *Context) {
			c.String(http.StatusNotFound, "404 NOT FOUND: %s\n", c.Path)
		})
	}
	c.Next()
}
