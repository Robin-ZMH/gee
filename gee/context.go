package gee

import (
	"encoding/json"
	"fmt"
	"net/http"
)

/*
Context exposes interfaces of obtaining information about the request,
Context provides methods to render the different data type
*/
type Context struct {
	// origin objects
	Writer http.ResponseWriter
	Req    *http.Request

	// request info
	Path   string
	Method string
	Params map[string]string

	// response info
	StatusCode int

	// middlewares
	handlers []HandlerFunc
	index    int
}

// NewContext will create a new Context object, and it will parse registered middlewares
// of the Engine object, add them to the "hanlers" field.
func NewContext(e *Engine, w http.ResponseWriter, req *http.Request) *Context {
	context := &Context{
		Writer: w,
		Req:    req,
		Path:   req.URL.Path,
		Method: req.Method,
		index:  -1,
	}
	path := req.URL.Path
	prefixes := parsePattern(path)
	for _, prefix := range prefixes {
		if group, ok := e.groups[prefix]; ok {
			context.handlers = append(context.handlers, group.middlewares...)
		}
	}
	return context
}

func (c *Context) Fail(code int, err string) {
	c.index = len(c.handlers)
	c.JSON(code, H{"message": err})
}

func (c *Context) Param(key string) string {
	value := c.Params[key]
	return value
}

func (c *Context) PostForm(key string) string {
	return c.Req.FormValue(key)
}

/*
Query returns the corresponding values of the desired key:

https://example.org/?a=1&a=2&b=&=3&&&&

[1, 2] <-Context.Query(a).
*/
func (c *Context) Query(key string) string {
	return c.Req.URL.Query().Get(key)
}

func (c *Context) Status(code int) {
	c.StatusCode = code
	c.Writer.WriteHeader(code)
}

func (c *Context) SetHeader(key string, val string) {
	c.Writer.Header().Set(key, val)
}

func (c *Context) String(code int, format string, data ...any) {
	c.SetHeader("Content-Type", "text/plain")
	c.Status(code)
	write_data := fmt.Sprintf(format, data...)
	c.Writer.Write([]byte(write_data))
}

func (c *Context) JSON(code int, obj interface{}) {
	c.SetHeader("Content-Type", "application/json")
	c.Status(code)
	encoder := json.NewEncoder(c.Writer)
	if err := encoder.Encode(obj); err != nil {
		panic(err)
	}
}

func (c *Context) Data(code int, data []byte) {
	c.Status(code)
	c.Writer.Write(data)
}

func (c *Context) HTML(code int, html string) {
	c.SetHeader("Content-Type", "text/html")
	c.Status(code)
	c.Writer.Write([]byte(html))
}

// Next will sequentially execute the middlewares(handlers).
func (c *Context) Next() {
	// every time calls Next() must increament the index
	c.index++
	for ; c.index < len(c.handlers); c.index++ {
		c.handlers[c.index](c)
	}
}
