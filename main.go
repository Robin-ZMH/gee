package main

// $ curl http://localhost:9999/
// URL.Path = "/"
// $ curl http://localhost:9999/hello
// Header["Accept"] = ["*/*"]
// Header["User-Agent"] = ["curl/7.54.0"]

import (
	"gee"
	"log"
)

func main() {
	engine := gee.NewEngine()
	log.Fatal(engine.Run(":9000"))
}
