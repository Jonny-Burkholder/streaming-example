package radix

import (
	"fmt"
	"net/http"
	"path"
	"strings"
	"testing"
)

// func entryToKey(method string, path string, handler bool) string {
// 	i := strings.IndexByte(path, '{')
// 	j := strings.IndexByte(path, '}')
// 	if i > 0 && j > 0 {
// 		strings.Split(path, "{**}")
// 	}
// 	method + ""
// }

func BenchmarkTree_WalkPrefix(b *testing.B) {
	tree := NewTree()

	tree.Insert("/api", func() string { return "GET /api" })
	tree.Insert("/api/v1", func() string { return "GET /api/v1" })
	tree.Insert("/api/v1/users", func() string { return "GET /api/v1/users" })
	tree.Insert("/api/v1/users/{id}", func() string { return "GET /api/v1/users/{id}" })
	tree.Insert("/api/v1/users/{id}/resources", func() string { return "GET /api/v1/users/{id}/resources" })

	searchKey := "/api/v1/users/25"
	fmt.Printf("Looking for a match for: %q\n", searchKey)
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		tree.WalkPrefix(
			path.Dir(searchKey),
			func(k string, v any) bool {
				if strings.ContainsAny(k, "{}") {
					fmt.Printf("found match: %q (%v)\n", k, v)
					return true
				}
				return false
			},
		)
	}
}

func TestTree_Routes(t *testing.T) {

	tree := NewTree()

	tree.Insert("/api", func() string { return "GET /api" })
	tree.Insert("/api/v1", func() string { return "GET /api/v1" })
	tree.Insert("/api/v1/users", func() string { return "GET /api/v1/users" })
	tree.Insert("/api/v1/users/{id}", func() string { return "GET /api/v1/users/{id}" })
	tree.Insert("/api/v1/users/{id}/resources", func() string { return "GET /api/v1/users/{id}/resources" })

	searchKey := "/api/v1/users/25"
	fmt.Printf("Looking for a match for: %q\n", searchKey)
	tree.WalkPrefix(
		path.Dir(searchKey),
		func(k string, v any) bool {
			fmt.Printf("-> %q\n", k)
			// if strings.ContainsAny(k, "{}") {
			// 	fmt.Printf("found match: %q (%v)\n", k, v)
			// 	return true
			// }
			return false
		},
	)
	fmt.Println()

	searchKey = "/api/v1/users/25/resources"
	fmt.Printf("Looking for a match for: %q\n", searchKey)
	tree.WalkPath(
		searchKey,
		func(k string, v any) bool {
			fmt.Printf("-> %q\n", k)
			// if strings.ContainsAny(k, "{}") {
			// 	fmt.Printf("found match: %q (%v)\n", k, v)
			// 	return true
			// }
			return false
		},
	)
	fmt.Println()

	fmt.Printf("Looking for a match for: %q\n", searchKey)
	fmt.Println(tree.FindLongestPrefix(searchKey))
}

func TestNewTree(t *testing.T) {

	t.Logf("Creating new radix tree...")
	rt := NewTree()
	if rt.Len() != 0 {
		t.Fatalf("Bad length, expected %v, got %v", 0, rt.Len())
	}

	var entries = []struct {
		Pre string
		Key string
		Val any
	}{
		{http.MethodGet, "/api", true},
		{http.MethodGet, "/v1", true},
		{http.MethodGet, "/users", true},
		{http.MethodPost, "/users", true},
		{http.MethodGet, "/{id}", true},
		{http.MethodPost, "/{id}", true},
		{http.MethodGet, "/api/users/{id}", true},
	}

	t.Logf("Inserting a few keys and values....")
	for _, entry := range entries {
		rt.Insert(entry.Key, entry)
	}

	// t.Logf("Searching by prefix (1)...")
	// searchKey := entries[4].Pre + " " + entries[4].Key
	// foundKey, _, _ := rt.FindLongestPrefix(searchKey)
	// fmt.Printf("FindLongestPrefix(%q) => %q\n", searchKey, foundKey)

	t.Logf("Walking the tree (1)...")
	rt.Walk(
		func(k string, v any) bool {
			fmt.Printf("key: %q, value: %v\n", k, v)
			return false
		},
	)

	t.Logf("Walking the tree [prefix frist] (2)...")
	rt.WalkPrefix(
		"/api/users",
		func(k string, v any) bool {
			fmt.Printf("key: %q\n", k)
			return false
		},
	)

	t.Logf("Walking the tree [path first] (3)...")
	rt.WalkPath(
		"/api/users",
		func(k string, v any) bool {
			fmt.Printf("key: %q\n", k)
			return false
		},
	)

	t.Logf("Checking out the minimum, and the maximum...")
	outMin, _, _ := rt.Min()
	fmt.Printf("min: %q\n", outMin)
	// if outMin != min {
	// 	t.Fatalf("bad minimum: %v %v", outMin, min)
	// }
	outMax, _, _ := rt.Max()
	fmt.Printf("max: %q\n", outMax)
	// if outMax != max {
	// 	t.Fatalf("bad maximum: %v %v", outMax, max)
	// }

	t.Logf("Checking the length...\n")
	if rt.Len() != len(entries) {
		t.Fatalf("Bad length, expected %v, got %v", len(entries), rt.Len())
	}

	t.Logf("Deleting %q...\n", "www.example.com")
	rt.Delete("www.example.com")

	t.Logf("Checking the length again...\n")
	if rt.Len() != len(entries)-1 {
		t.Fatalf("Bad length, expected %v, got %v", len(entries)-1, rt.Len())
	}
}
