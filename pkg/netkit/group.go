package netkit

import (
	"fmt"
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
	// standardize the group path
	if len(group) > 0 && group[0] != '/' {
		group = "/" + group
	}
	if group[len(group)-1] == '/' {
		group = group[:len(group)-1]
	}
	// create a new group
	g := &Group{
		mux:   rm,
		group: group,
	}
	// register base path for group
	rm.Handle("*", join(group, "/"), http.StripPrefix(g.group, g.handleGroupRoot()))
	// return new group
	return g
}

// join cleans and joins the group path with the pattern and
// returns the joined string
func join(group, pattern string) string {
	s := filepath.ToSlash(filepath.Join(group, pattern))
	if pattern[len(pattern)-1] == '/' {
		s += "/"
	}
	return s
}

func (g *Group) Get(pattern string, handler http.HandlerFunc) {
	g.mux.Handle(http.MethodGet, join(g.group, pattern), http.StripPrefix(g.group, handler))
}

func (g *Group) Post(pattern string, handler http.HandlerFunc) {
	g.mux.Handle(http.MethodPost, join(g.group, pattern), http.StripPrefix(g.group, handler))
}

func (g *Group) Put(pattern string, handler http.HandlerFunc) {
	g.mux.Handle(http.MethodPut, join(g.group, pattern), http.StripPrefix(g.group, handler))
}

func (g *Group) Delete(pattern string, handler http.HandlerFunc) {
	g.mux.Handle(http.MethodDelete, join(g.group, pattern), http.StripPrefix(g.group, handler))
}

func (g *Group) handleGroupRoot() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		WriteRaw(w, r, 200, []byte(fmt.Sprintf("group: %s", g.group)))
	})
}
