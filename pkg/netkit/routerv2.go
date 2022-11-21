package netkit

import (
	"fmt"
	"log"
	"mime"
	"net/http"
	"path/filepath"
	"strings"
	"sync"

	"github.com/Jonny-Burkholder/streaming-example/pkg/trees/radix"
)

type RouterV2 struct {
	lock   sync.Mutex
	routes *radix.Tree
	groups []string
}

func NewRouterV2() *RouterV2 {
	return &RouterV2{
		routes: radix.NewTree(),
		groups: make([]string, 0),
	}
}

func urlMapping(r rune) rune {
	switch {
	case 'a' <= r && r <= 'z':
		return r - 'a' - 'A'
	}
	return r
}

func sanitize(method string, pattern string) string {
	return method + pattern //strings.Map(urlMapping, pattern)
}

func (rt *RouterV2) Handle(method string, pattern string, handler http.Handler) {
	entry := routeEntry{
		method:  method,
		pattern: pattern,
		handler: handler,
	}
	rt.routes.Insert(sanitize(method, pattern), entry)
}

func (rt *RouterV2) HandleFunc(method, pattern string, handler func(http.ResponseWriter, *http.Request)) {
	entry := routeEntry{
		method:  method,
		pattern: pattern,
		handler: http.HandlerFunc(handler),
	}
	rt.routes.Insert(sanitize(method, pattern), entry)
}

func (rt *RouterV2) Get(pattern string, handler http.HandlerFunc) {
	rt.HandleFunc(http.MethodGet, pattern, handler)
}

func (rt *RouterV2) Post(pattern string, handler http.HandlerFunc) {
	rt.HandleFunc(http.MethodPost, pattern, handler)
}

func (rt *RouterV2) Put(pattern string, handler http.HandlerFunc) {
	rt.HandleFunc(http.MethodPut, pattern, handler)
}

func (rt *RouterV2) Delete(pattern string, handler http.HandlerFunc) {
	rt.HandleFunc(http.MethodDelete, pattern, handler)
}

func (rt *RouterV2) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.RequestURI == "*" {
		if r.ProtoAtLeast(1, 1) {
			w.Header().Set("Connection", "close")
		}
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	matched, entry, found := rt.routes.FindLongestPrefix(sanitize(r.Method, r.URL.Path))
	if !found {
		http.NotFound(w, r)
		return
	}
	log.Printf("path: %q, matched: %q\n", r.URL.Path, matched)
	entry.(routeEntry).handler.ServeHTTP(w, r)
	return
}

func (rt *RouterV2) MetricsHandler(w http.ResponseWriter, r *http.Request) {
	//
	// Write page heading
	sb := new(strings.Builder)
	sb.WriteString("<h2>All route entries</h2>")
	//
	// Walk all routes
	sb.WriteString("<h4>Routes:</h4>")
	rt.lock.Lock()
	rt.routes.Walk(func(k string, v any) bool {
		if ent, castOkay := v.(routeEntry); castOkay {
			sb.WriteString(ent.String())
			sb.WriteString("<br>")
		}
		return false
	})
	rt.lock.Unlock()
	//
	// Write Content-Type header, and write everything to the http.ResponseWriter
	w.Header().Set("Content-Type", mime.TypeByExtension(".html"))
	w.WriteHeader(200)
	fmt.Fprintf(w, "%s", sb)
}

func (rt *RouterV2) NewGroup(group string) *RouterV2Group {
	if len(group) > 0 && group[0] != '/' {
		group = "/" + group
	}
	if group[len(group)-1] != '/' {
		group += "/"
	}
	rt.groups = append(rt.groups, group)
	return &RouterV2Group{
		prefix: group,
		router: rt,
	}
}

type RouterV2Group struct {
	prefix string
	router *RouterV2
}

// joinGroup cleans and joins the group path with the pattern and
// returns the joined string
func joinGroup(group, pattern string) string {
	s := filepath.ToSlash(filepath.Join(group, pattern))
	if pattern[len(pattern)-1] == '/' {
		s += "/"
	}
	return s
}

func (rg *RouterV2Group) Handle(method string, pattern string, handler http.Handler) {
	rg.router.Handle(method, joinGroup(rg.prefix, pattern), http.StripPrefix(rg.prefix, handler))
}

func (rg *RouterV2Group) HandleFunc(method, pattern string, handler func(http.ResponseWriter, *http.Request)) {
	rg.router.Handle(method, joinGroup(rg.prefix, pattern), http.StripPrefix(rg.prefix, http.HandlerFunc(handler)))
}

func (rg *RouterV2Group) Get(pattern string, handler http.HandlerFunc) {
	rg.router.Handle(http.MethodGet, joinGroup(rg.prefix, pattern), http.StripPrefix(rg.prefix, handler))
}

func (rg *RouterV2Group) Post(pattern string, handler http.HandlerFunc) {
	rg.router.Handle(http.MethodPost, joinGroup(rg.prefix, pattern), http.StripPrefix(rg.prefix, handler))
}

func (rg *RouterV2Group) Put(pattern string, handler http.HandlerFunc) {
	rg.router.Handle(http.MethodPut, joinGroup(rg.prefix, pattern), http.StripPrefix(rg.prefix, handler))
}

func (rg *RouterV2Group) Delete(pattern string, handler http.HandlerFunc) {
	rg.router.Handle(http.MethodDelete, joinGroup(rg.prefix, pattern), http.StripPrefix(rg.prefix, handler))
}
