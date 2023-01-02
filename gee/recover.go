package gee

import (
	"fmt"
	"log"
	"net/http"
	"runtime"
	"strings"
)

func trace(message string) string {
	var str strings.Builder
	str.WriteString(message + "TraceBack:\n")
	pcs := make([]uintptr, 64)
	count := runtime.Callers(0, pcs)
	for _, pc := range pcs[:count] {
		fn := runtime.FuncForPC(pc)
		file, line := fn.FileLine(pc)
		str.WriteString(fmt.Sprintf("\t%s:%d\n", file, line))
	}
	return str.String()
}

func Recovery() HandlerFunc {
	return func(ctx *Context) {
		defer func() {
			if err := recover(); err != nil {
				message := fmt.Sprintf("%s\n", err)
				log.Printf("%s\n", trace(message))
				ctx.Fail(http.StatusInternalServerError, "Internal Server Error")
			}
		}()
		ctx.Next()
	}
}
