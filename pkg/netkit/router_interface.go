package netkit

import (
	"net/http"
)

type RouterInterface interface {
	Handle(method string, pattern string, handler http.Handler)
	HandleFunc(method, pattern string, handler func(http.ResponseWriter, *http.Request))
	Get(pattern string, handler http.HandlerFunc)
	Post(pattern string, handler http.HandlerFunc)
	Put(pattern string, handler http.HandlerFunc)
	Delete(pattern string, handler http.HandlerFunc)
	ServeHTTP(w http.ResponseWriter, r *http.Request)
}
