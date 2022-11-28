package testing

import (
	"fmt"
	"net/http"
	"net/url"
	"regexp"
)

type reRoute struct {
	method  string
	pattern string
	re      *regexp.Regexp
	h       http.HandlerFunc
}

func (r *reRoute) String() string {
	return fmt.Sprintf("method=%q, pattern=%q, regex=%q, handler=%v\n", r.method, r.pattern, r.re, r.h)
}

type RegexURLMatcher struct {
	routes map[string]*reRoute
}

func NewRegexURLMatcher() *RegexURLMatcher {
	return &RegexURLMatcher{
		routes: make(map[string]*reRoute),
	}
}

func (re *RegexURLMatcher) HandleFunc(method string, pattern string, handler http.HandlerFunc) {
	compiled, err := regexp.Compile(`^` + pattern + `$`)
	if err != nil {
		panic(err)
	}
	re.routes[pattern] = &reRoute{
		method:  method,
		pattern: pattern,
		re:      compiled,
		h:       handler,
	}

}

func (re *RegexURLMatcher) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	for _, route := range re.routes {
		if route.method == r.Method && route.re.MatchString(r.URL.Path) {
			if sub := route.re.SubexpNames(); len(sub) > 1 {
				for at, s := range sub[1:] {
					if r.Form == nil {
						r.Form = url.Values{}
					}
					r.Form.Set(s, route.re.FindStringSubmatch(r.URL.Path)[at+1])
				}
			}
			route.h(w, r)
			return
		}
	}
	http.NotFound(w, r)
}
