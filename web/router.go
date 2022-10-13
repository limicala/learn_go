package web

import (
	"fmt"
	"net/http"
)

type HandleFunc func(http.ResponseWriter, *http.Request)
type router struct {
	handlers map[string]HandleFunc
}

func newRouter() *router {
	return &router{handlers: make(map[string]HandleFunc)}
}

func (r *router) addRoute(method string, pattern string, handler HandleFunc) {
	key := method + ":" + pattern
	r.handlers[key] = handler
}

func (r *router) handle(rsp http.ResponseWriter, req *http.Request) {
	key := req.Method + ":" + req.URL.Path
	if handler, ok := r.handlers[key]; ok {
		handler(rsp, req)
	} else {
		fmt.Fprintf(rsp, "404 NOT FOUND: %s\n", req.URL)
	}
}
