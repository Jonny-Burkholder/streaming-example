package netkit

import (
	"fmt"
	"mime"
	"net/http"
	"sort"
	"strings"
	"sync"
)

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
	if m.method == "*" {
		return fmt.Sprintf("[%s]&nbsp;&nbsp;&nbsp;&nbsp;<a href=\"%s\">%s</a>", "ANY", m.pattern, m.pattern)
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
	ErrHandler:    http.HandlerFunc(HandleErrors),
	MetricsOn:     true,
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
		mux.Handle(http.MethodGet, "/static/", conf.StaticHandler)
	}
	if conf.ErrHandler != nil {
		mux.Handle(http.MethodGet, "/error/", conf.ErrHandler)
	}
	if conf.MetricsOn {
		mux.Handle(http.MethodGet, "/metrics", mux.handleMetrics())
	}
	return mux
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

func (rm *Router) Get(pattern string, handler http.HandlerFunc) {
	rm.HandleFunc(http.MethodGet, pattern, handler)
}

func (rm *Router) Post(pattern string, handler http.HandlerFunc) {
	rm.HandleFunc(http.MethodPost, pattern, handler)
}

func (rm *Router) Put(pattern string, handler http.HandlerFunc) {
	rm.HandleFunc(http.MethodPut, pattern, handler)
}

func (rm *Router) Delete(pattern string, handler http.HandlerFunc) {
	rm.HandleFunc(http.MethodDelete, pattern, handler)
}

func (rm *Router) Static(pattern string, path string) {
	staticHandler := http.StripPrefix(pattern, http.FileServer(http.Dir(path)))
	rm.Handle(http.MethodGet, pattern, staticHandler)
}

func (rm *Router) handleMetrics() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//
		// Write page heading
		sb := new(strings.Builder)
		sb.WriteString("<h2>All route entries</h2>")
		//
		// Collect base routes (routes that have sub-routes)
		var base []routeEntry
		sb.WriteString("<h4>Base routes:</h4>")
		rm.lock.Lock()
		for _, entry := range rm.entrySet {
			base = append(base, entry)
		}
		rm.lock.Unlock()
		//
		// Sort and write base routes
		sort.Slice(base, func(i, j int) bool { return base[i].pattern < base[j].pattern })
		for _, ent := range base {
			sb.WriteString(ent.String())
			sb.WriteString("<br>")
		}
		//
		//
		//
		// Collect sub routes (routes that are a continuation of a base route)
		var sub []routeEntry
		sb.WriteString("<h4>Sub routes:</h4>")
		rm.lock.Lock()
		for _, entry := range rm.entryMap {
			p := entry.pattern
			if p[len(p)-1] == '/' {
				continue
			}
			sub = append(sub, entry)
		}
		rm.lock.Unlock()
		//
		// Sort and write base routes
		sort.Slice(sub, func(i, j int) bool { return sub[i].pattern < sub[j].pattern })
		for _, ent := range sub {
			sb.WriteString(ent.String())
			sb.WriteString("<br>")
		}
		//
		// Write Content-Type header, and write everything to the http.ResponseWriter
		w.Header().Set("Content-Type", mime.TypeByExtension(".html"))
		w.WriteHeader(200)
		fmt.Fprintf(w, "%s", sb)
	})
}

func (rm *Router) _ServeHTTP(w http.ResponseWriter, r *http.Request) {
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

func (rm *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.RequestURI == "*" {
		if r.ProtoAtLeast(1, 1) {
			w.Header().Set("Connection", "close")
		}
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	found, entry := rm.matchV0(r)
	if found {
		if entry.method != r.Method && entry.method != "*" {
			entry.handler.ServeHTTP(w, r)
			return
		}
	}
	http.NotFound(w, r)
	return
	// meth, _, hdlr := rm.match(r.URL.Path)
	// if meth != r.Method && meth != "*" {
	// 	hdlr = http.HandlerFunc(
	// 		func(w http.ResponseWriter, r *http.Request) {
	// 			code := http.StatusMethodNotAllowed
	// 			http.Error(w, http.StatusText(code), code)
	// 		},
	// 	)
	// }
	// if hdlr == nil {
	// 	hdlr = http.NotFoundHandler()
	// }
	// if rm.withLogging {
	// 	// if logging is configured, then log, otherwise skip
	// 	hdlr = HandleWithLogging(rm.logger, hdlr)
	// }
	// hdlr.ServeHTTP(w, r)
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
	for _, entry := range rm.entrySet {
		entries = append(entries, fmt.Sprintf("%s %s\n", entry.method, entry.pattern))
	}
	return entries
}

// match attempts to locate a handler on a handler map given a path string. The
// most-specific (longest) pattern wins. It returns the method, pattern and handler.
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
		// inline check for same prefix has prefix
		if len(path) >= len(e.pattern) && path[0:len(e.pattern)] == e.pattern {
			return e.method, e.pattern, e.handler
		}
	}
	return "", "", nil
}

func (rm *Router) matchV0(r *http.Request) (bool, *routeEntry) {
	// check each entry in the map for a match against
	// the current request.
	for path, entry := range rm.entryMap {
		// check for the presence of an opening path key
		// byte, and if we cannot find one, we will do
		// a simple path match.
		beg := strings.IndexByte(path, '{')
		if beg == -1 {
			// check for a basic path match, and break
			return path == r.URL.Path, &entry
		}
		// otherwise, we found an opening path key
		// byte, so next, we will check for a closed
		// path key byte.
		end := strings.IndexByte(path, '}')
		if end == -1 {
			// we have not found a closed path
			// key byte, which means there must
			// be an error with p1, so a match
			// will be impossible at this point,
			// try next path entry...
			continue
		}
		// next, make certain we have a match
		if path[:beg] != r.URL.Path[:beg] {
			// we don't so we continue to the next one
			continue
		}
		// we have found a closed path key byte, so
		// we should now extract the key from path
		key := path[beg+1 : end]
		// and then find the value that maps to our key.
		i := strings.IndexByte(r.URL.Path[beg:], '/')
		if i == -1 {
			r.URL.Query().Add(key, r.URL.Path[beg:])
		} else {
			r.URL.Query().Add(key, r.URL.Path[beg:beg+i])
		}
		// now return our match, and entry
		return true, &entry
	}
	return false, nil
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
