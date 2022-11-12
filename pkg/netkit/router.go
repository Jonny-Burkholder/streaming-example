package netkit

import (
	"fmt"
	"net/http"
	"path"
	"sort"
	"strings"
	"sync"
)

type Route struct {
	Method  string
	Pattern string
	Handler http.Handler
}

type routeEntry struct {
	method  string
	pattern string
	handler http.Handler
}

func (m routeEntry) String() string {
	if m.method == http.MethodGet {
		return fmt.Sprintf("[%s]&nbsp;&nbsp;&nbsp;&nbsp;<a href=\"%s\">%s</a>", m.method, m.pattern, m.pattern)
	}
	if m.method == http.MethodPost {
		return fmt.Sprintf("[%s]&nbsp;&nbsp;&nbsp;%s", m.method, m.pattern)
	}
	if m.method == http.MethodPut {
		return fmt.Sprintf("[%s]&nbsp;&nbsp;&nbsp;&nbsp;%s", m.method, m.pattern)
	}
	if m.method == http.MethodDelete {
		return fmt.Sprintf("[%s]&nbsp;%s", m.method, m.pattern)
	}
	return fmt.Sprintf("[%s]&nbsp;%s", m.method, m.pattern)
}

type Config struct {
	StaticHandler http.Handler
	ErrHandler    http.Handler
	MetricsOn     bool
	LoggingLevel  int
}

var defaultConfig = &Config{
	StaticHandler: HandleStatic("/static/", "web/static/"),
	ErrHandler:    http.HandlerFunc(NotFound),
	MetricsOn:     false,
	LoggingLevel:  LevelInfo,
}

type Router struct {
	lock        sync.Mutex
	entryMap    map[string]routeEntry
	entrySet    []routeEntry
	logger      *Logger
	withLogging bool
}

func NewRouter(conf *Config) *Router {
	if conf == nil {
		conf = defaultConfig
	}
	mux := &Router{
		entryMap: make(map[string]routeEntry),
		entrySet: make([]routeEntry, 0),
		logger:   NewLogger(LevelInfo),
	}
	if conf.LoggingLevel < LevelOff {
		mux.logger = NewLogger(conf.LoggingLevel)
		mux.withLogging = true
	}
	if conf.StaticHandler != nil {
		mux.Get("/static/", conf.StaticHandler)
	}
	if conf.ErrHandler != nil {
		mux.Get("/error/", conf.ErrHandler)
	}
	if conf.MetricsOn {
		mux.Get("/metrics", HandleMetrics("Registered Entries", mux.entries()))
	}
	return mux
}

func (rm *Router) Group(name string) *Router {
	return rm
}

func (rm *Router) Handle(method string, pattern string, handler http.Handler) {
	rm.lock.Lock()
	defer rm.lock.Unlock()

	if pattern == "" {
		panic("http: invalid pattern")
	}
	if handler == nil {
		panic("http: nil handler")
	}
	if _, exist := rm.entryMap[pattern]; exist {
		panic("http: multiple registrations for " + pattern)
	}
	entry := routeEntry{
		method:  method,
		pattern: pattern,
		handler: handler,
	}
	rm.entryMap[pattern] = entry
	if pattern[len(pattern)-1] == '/' {
		rm.entrySet = appendSorted(rm.entrySet, entry)
	}
}

func (rm *Router) HandleFunc(method, pattern string, handler func(http.ResponseWriter, *http.Request)) {
	if handler == nil {
		panic("http: nil handler")
	}
	rm.Handle(method, pattern, http.HandlerFunc(handler))
}

func (rm *Router) Forward(oldpattern string, newpattern string) {
	rm.Handle(http.MethodGet, oldpattern, http.RedirectHandler(newpattern, http.StatusTemporaryRedirect))
}

func (rm *Router) Get(pattern string, handler http.Handler) {
	rm.Handle(http.MethodGet, pattern, handler)
}

func (rm *Router) Post(pattern string, handler http.Handler) {
	rm.Handle(http.MethodPost, pattern, handler)
}

func (rm *Router) Put(pattern string, handler http.Handler) {
	rm.Handle(http.MethodPut, pattern, handler)
}

func (rm *Router) Delete(pattern string, handler http.Handler) {
	rm.Handle(http.MethodDelete, pattern, handler)
}

func (rm *Router) Static(pattern string, path string) {
	staticHandler := http.StripPrefix(pattern, http.FileServer(http.Dir(path)))
	rm.Handle(http.MethodGet, pattern, staticHandler)
}

func (rm *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.RequestURI == "*" {
		if r.ProtoAtLeast(1, 1) {
			w.Header().Set("Connection", "close")
		}
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	meth, _, hdlr := rm.match(r.URL.Path)
	if meth != r.Method && meth != "*" {
		hdlr = http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				code := http.StatusMethodNotAllowed
				http.Error(w, http.StatusText(code), code)
			},
		)
	}
	if hdlr == nil {
		hdlr = http.NotFoundHandler()
	}
	if rm.withLogging {
		// if logging is configured, then log, otherwise skip
		hdlr = HandleWithLogging(rm.logger, hdlr)
	}
	hdlr.ServeHTTP(w, r)
}

func (rm *Router) Len() int {
	return len(rm.entrySet)
}

func (rm *Router) Less(i, j int) bool {
	return rm.entrySet[i].pattern < rm.entrySet[j].pattern
}

func (rm *Router) Swap(i, j int) {
	rm.entrySet[j], rm.entrySet[i] = rm.entrySet[i], rm.entrySet[j]
}

func (rm *Router) Search(x string) int {
	return sort.Search(
		len(rm.entrySet), func(i int) bool {
			return rm.entrySet[i].pattern >= x
		},
	)
}

func (rm *Router) entries() []string {
	rm.lock.Lock()
	defer rm.lock.Unlock()
	var entries []string
	for _, entry := range rm.entryMap {
		entries = append(entries, fmt.Sprintf("%s %s\n", entry.method, entry.pattern))
	}
	return entries
}

// match attempts to locate a handler on a handler map given a
// path string; most-specific (longest) pattern wins
func (rm *Router) match(path string) (string, string, http.Handler) {
	// first, check for exact match
	e, ok := rm.entryMap[path]
	if ok {
		return e.method, e.pattern, e.handler
	}
	// then, check for longest valid match. mux.entrySet
	// contains all patterns that end in "/" sorted
	// from longest to shortest
	for _, e = range rm.entrySet {
		if strings.HasPrefix(path, e.pattern) {
			return e.method, e.pattern, e.handler
		}
	}
	return "", "", nil
}

// cleanPath returns the canonical path for p, eliminating . and .. elements
func cleanPath(p string) string {
	if p == "" {
		return "/"
	}
	if p[0] != '/' {
		p = "/" + p
	}
	np := path.Clean(p)
	// path.Clean removes trailing slash except for root;
	// put the trailing slash back if necessary.
	if p[len(p)-1] == '/' && np != "/" {
		// Fast path for common case of p being the string we want:
		if len(p) == len(np)+1 && strings.HasPrefix(p, np) {
			np = p
		} else {
			np += "/"
		}
	}
	return np
}

func appendSorted(es []routeEntry, e routeEntry) []routeEntry {
	n := len(es)
	i := sort.Search(
		n, func(i int) bool {
			return len(es[i].pattern) < len(e.pattern)
		},
	)
	if i == n {
		return append(es, e)
	}
	// we now know that i points at where we want to insert
	es = append(es, routeEntry{}) // try to grow the slice in place, any entry works.
	copy(es[i+1:], es[i:])        // Move shorter entries down
	es[i] = e
	return es
}
