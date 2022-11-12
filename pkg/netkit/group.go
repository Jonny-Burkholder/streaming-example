package netkit

import (
	"log"
	"net/http"
	"path/filepath"
)

type RouterGroup interface {
	Get(pattern string, handler http.HandlerFunc)
	Post(pattern string, handler http.HandlerFunc)
	Put(pattern string, handler http.HandlerFunc)
	Delete(pattern string, handler http.HandlerFunc)
}

type Group struct {
	mux   *Router
	group string
}

func (rm *Router) NewGroup(group string) *Group {
	return &Group{
		mux:   rm,
		group: sanitize(group),
	}
}

// sanitize makes cleans the path "foo" into "/foo"
func sanitize(s string) string {
	if len(s) > 0 && s[0] != '/' {
		s = "/" + s
	}
	if s[len(s)-1] == '/' {
		s = s[:len(s)-1]
	}
	return s
}

func (g *Group) clean(pattern string) string {
	gpattern := filepath.Join(g.group, pattern)
	if pattern[len(pattern)-1] == '/' {
		gpattern += "/"
	}
	return gpattern
}

func (g *Group) Get(pattern string, handler http.HandlerFunc) {
	log.Printf("g.mux.Handle(%s, g.clean(%q), http.StripPrefix(%q, handler))\n", http.MethodGet, pattern, g.group)
	g.mux.Handle(http.MethodGet, g.clean(pattern), http.StripPrefix(g.group, handler))
}

func (g *Group) Post(pattern string, handler http.HandlerFunc) {
	g.mux.Post(g.clean(pattern), handler)
}

func (g *Group) Put(pattern string, handler http.HandlerFunc) {
	g.mux.Put(g.clean(pattern), handler)
}

func (g *Group) Delete(pattern string, handler http.HandlerFunc) {
	g.mux.Delete(g.clean(pattern), handler)
}

func (g *Group) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	g.mux.ServeHTTP(w, r)
}
