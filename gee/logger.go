package gee

import (
	"log"
	"time"
)

func startTimer(c *Context) func() {
	t := time.Now()
	return func() {
		d := time.Since(t)
		log.Println("=================================")
		defer log.Println("=================================")
		log.Printf("[%d] %s in %v", c.StatusCode, c.Req.RequestURI, d)
	}
}

func Logger() HandlerFunc {
	return func(c *Context) {
		// start timer
		stop := startTimer(c)
		defer stop()
		// Process request
		c.Next()
	}
}
