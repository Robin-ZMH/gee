package gee

import (
	"log"
	"reflect"
	"testing"
)

func newTestRouter() *router {
	r := newRouter()
	r.register("GET", "/", nil)
	r.register("GET", "/hello/:name", nil)
	r.register("GET", "/hello/b/c", nil)
	r.register("GET", "/hi/:name", nil)
	r.register("GET", "/assets/*filepath", nil)
	return r
}

func TestParsePattern(t *testing.T) {
	ok := reflect.DeepEqual(parsePattern("/p/:name"), []string{"p", ":name"})
	ok = ok && reflect.DeepEqual(parsePattern("/p/*"), []string{"p", "*"})
	ok = ok && reflect.DeepEqual(parsePattern("/p/*name/*"), []string{"p", "*name"})
	if !ok {
		t.Fatal("test parsePattern failed")
	}
}

func TestMatchPattern(t *testing.T) {
	r := newTestRouter()
	pattern := r.matchPattern("GET", "/hello/geektutu")

	if pattern != "/hello/:name" {
		t.Fatalf("result is %s\n, should match /hello/:name\n", pattern)
	}

	log.Printf("matched path: %s\n", pattern)
}

func TestGetParams(t *testing.T) {
	r := newTestRouter()
	params := r.getParams("/hello/geektutu", "/hello/:name")
	if params["name"] != "geektutu" {
		t.Fatal("name should be equal to 'geektutu'")
	}
	log.Printf("params['name']: %s\n", params["name"])
}

func TestMatchPattern2(t *testing.T) {
	r := newTestRouter()
	pattern := r.matchPattern("GET", "/assets/file1.txt")

	ok1 := pattern == "/assets/*filepath"
	if !ok1 {
		t.Fatal("pattern shoule be /assets/*filepath & filepath shoule be file1.txt")
	}
	pattern = r.matchPattern("GET", "/assets/css/test.css")
	ok2 := pattern == "/assets/*filepath"
	if !ok2 {
		t.Fatal("pattern shoule be /assets/*filepath & filepath shoule be css/test.css")
	}
}
