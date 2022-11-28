package testing

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/Jonny-Burkholder/streaming-example/pkg/netkit"
)

func getPathParam(uri string) (string, bool) {
	if uri[len(uri)-1] == '/' {
		uri = uri[:len(uri)-1]
	}
	i := len(uri) - 1
	for i >= 0 && uri[i] != '/' {
		i--
	}
	param := uri[i+1:]
	return param, param != "" && param != uri[(len(uri)-1)/2:]
}

func getPathID(uri string) int {
	param, found := getPathParam(uri)
	if !found {
		return -1
	}
	id, err := strconv.ParseInt(param, 10, 0)
	if err != nil {
		return -1
	}
	return int(id)
}

func Benchmark_HTTP_StandardLibrary(b *testing.B) {
	mux := http.NewServeMux()
	mux.HandleFunc(
		"/", func(w http.ResponseWriter, r *http.Request) {
			http.RedirectHandler("/api", http.StatusTemporaryRedirect)
			return
		},
	)
	mux.HandleFunc(
		"/api", func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, "API ROOT")
			return
		},
	)
	mux.HandleFunc(
		"/api/users", func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, "USERS LIST ROOT")
			return
		},
	)
	mux.HandleFunc(
		"/api/users/", func(w http.ResponseWriter, r *http.Request) {
			id := getPathID(r.URL.Path)
			if id < 0 {
				http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
				return
			}
			fmt.Fprintf(w, "USERS %d ROOT", id)
			return
		},
	)
	ts := httptest.NewServer(mux)
	defer ts.Close()

	testURLs := []struct {
		URL      string
		Response string
	}{
		{"/", ""},
		{"/api", "API ROOT"},
		{"/api/users", "USERS LIST ROOT"},
		{"/api/users/1", "USERS 1 ROOT"},
		{"/api/users/16", "USERS 16 ROOT"},
		{"/api/users/32", "USERS 32 ROOT"},
		{"/api/users/256", "USERS 256 ROOT"},
	}

	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		for _, testURL := range testURLs {
			res, err := http.Get(ts.URL + testURL.URL)
			if err != nil {
				b.Errorf("%s", err)
				log.Fatal(err)
			}
			resp, err := io.ReadAll(res.Body)
			res.Body.Close()
			if err != nil {
				b.Errorf("%s", err)
				log.Fatal(err)
			}
			if string(resp) != testURL.Response {
				b.Errorf("bad response: got=%q, wanted=%q\n", string(resp), testURL.Response)
			}
		}
	}
}

func Benchmark_HTTP_NetKitRouter(b *testing.B) {
	mux := netkit.NewRouter(
		&netkit.Config{
			StaticHandler: nil,
			ErrHandler:    nil,
			MetricsOn:     false,
			LoggingLevel:  netkit.LevelOff,
		},
	)
	mux.HandleFunc(
		http.MethodGet, "/", func(w http.ResponseWriter, r *http.Request) {
			http.RedirectHandler("/api", http.StatusTemporaryRedirect)
			return
		},
	)
	mux.HandleFunc(
		http.MethodGet,
		"/api", func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, "API ROOT")
			return
		},
	)
	mux.HandleFunc(
		http.MethodGet,
		"/api/users", func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, "USERS LIST ROOT")
			return
		},
	)
	mux.HandleFunc(
		http.MethodGet,
		"/api/users/", func(w http.ResponseWriter, r *http.Request) {
			id := getPathID(r.URL.Path)
			if id < 0 {
				http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
				return
			}
			fmt.Fprintf(w, "USERS %d ROOT", id)
			return
		},
	)
	ts := httptest.NewServer(mux)
	defer ts.Close()

	testURLs := []struct {
		URL      string
		Response string
	}{
		{"/", ""},
		{"/api", "API ROOT"},
		{"/api/users", "USERS LIST ROOT"},
		{"/api/users/1", "USERS 1 ROOT"},
		{"/api/users/16", "USERS 16 ROOT"},
		{"/api/users/32", "USERS 32 ROOT"},
		{"/api/users/256", "USERS 256 ROOT"},
	}

	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		for _, testURL := range testURLs {
			res, err := http.Get(ts.URL + testURL.URL)
			if err != nil {
				b.Errorf("%s", err)
				log.Fatal(err)
			}
			resp, err := io.ReadAll(res.Body)
			res.Body.Close()
			if err != nil {
				b.Errorf("%s", err)
				log.Fatal(err)
			}
			if string(resp) != testURL.Response {
				b.Errorf("bad response: got=%q, wanted=%q\n", string(resp), testURL.Response)
			}
		}
	}
}

func Benchmark_HTTP_CustomRegexRouter(b *testing.B) {
	mux := NewRegexURLMatcher()
	mux.HandleFunc(
		http.MethodGet, `/`, func(w http.ResponseWriter, r *http.Request) {
			http.RedirectHandler("/api", http.StatusTemporaryRedirect)
			return
		},
	)
	mux.HandleFunc(
		http.MethodGet,
		`/api`, func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, "API ROOT")
			return
		},
	)
	mux.HandleFunc(
		http.MethodGet,
		`/api/users`, func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, "USERS LIST ROOT")
			return
		},
	)
	// `/api/users/(?P<id>\d+)`

	mux.HandleFunc(
		http.MethodGet,
		`/api/users/{id}`, func(w http.ResponseWriter, r *http.Request) {
			id := r.FormValue("id")
			if id == "" {
				http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
				return
			}
			fmt.Fprintf(w, "USERS %v ROOT", id)
			return
		},
	)
	ts := httptest.NewServer(mux)
	defer ts.Close()

	testURLs := []struct {
		URL      string
		Response string
	}{
		{"/", ""},
		{"/api", "API ROOT"},
		{"/api/users", "USERS LIST ROOT"},
		{"/api/users/1", "USERS 1 ROOT"},
		{"/api/users/16", "USERS 16 ROOT"},
		{"/api/users/32", "USERS 32 ROOT"},
		{"/api/users/256", "USERS 256 ROOT"},
	}

	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		for _, testURL := range testURLs {
			res, err := http.Get(ts.URL + testURL.URL)
			if err != nil {
				b.Errorf("%s", err)
				log.Fatal(err)
			}
			resp, err := io.ReadAll(res.Body)
			res.Body.Close()
			if err != nil {
				b.Errorf("%s", err)
				log.Fatal(err)
			}
			if string(resp) != testURL.Response {
				b.Errorf("bad response: got=%q, wanted=%q\n", string(resp), testURL.Response)
			}
		}
	}
}
