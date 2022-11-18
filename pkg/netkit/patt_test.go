package netkit

import (
	"fmt"
	"net/url"
	"testing"
)

func checkMatch(print bool, key string, val string, match int) bool {
	if print {
		fmt.Printf("key=%q, val=%q, match=%d\n", key, val, match)
	}
	return key != "" && val != "" && match != 0
}

var res1 any
var res2 any
var res3 any

// average: 180 ns/op, 32 B/op, 2 allocs/op
func BenchmarkMatchV0(b *testing.B) {
	var m int
	var k, v string
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		m, k, v = matchV0("/api/users/jobs/{jobID}", "/api/users/jobs/12")
		if !checkMatch(false, k, v, m) {
			b.Fatal(m)
		}
		m, k, v = matchV0("/api/users/jobs/{jobID}", "/api/users/jobs/1234")
		if !checkMatch(false, k, v, m) {
			b.Fatal(m)
		}
		m, k, v = matchV0("/api/users/{id}", "/api/users/12")
		if !checkMatch(false, k, v, m) {
			b.Fatal(m)
		}
		m, k, v = matchV0("/api/users/{id}", "/api/users/123456")
		if !checkMatch(false, k, v, m) {
			b.Fatal(m)
		}
		m, k, v = matchV0("/api/users/{id}/", "/api/users/79/")
		if !checkMatch(false, k, v, m) {
			b.Fatal(m)
		}
		m, k, v = matchV0("/api/users/{id}/foo", "/api/users/123456/foo")
		if !checkMatch(false, k, v, m) {
			b.Fatal(m)
		}
		res1 = k
		res2 = v
		res3 = m
	}
	_ = res1
	_ = res2
	_ = res3
}

// average: 180 ns/op, 32 B/op, 2 allocs/op
func BenchmarkMatchPatternV1(b *testing.B) {
	var m int
	var k, v string
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		k, v, m = MatchPatternV1("/api/users/jobs/{jobID}", "/api/users/jobs/12")
		if !checkMatch(false, k, v, m) {
			b.Fatal(m)
		}
		k, v, m = MatchPatternV1("/api/users/jobs/{jobID}", "/api/users/jobs/1234")
		if !checkMatch(false, k, v, m) {
			b.Fatal(m)
		}
		k, v, m = MatchPatternV1("/api/users/{id}", "/api/users/12")
		if !checkMatch(false, k, v, m) {
			b.Fatal(m)
		}
		k, v, m = MatchPatternV1("/api/users/{id}", "/api/users/123456")
		if !checkMatch(false, k, v, m) {
			b.Fatal(m)
		}
		k, v, m = MatchPatternV1("/api/users/{id}/", "/api/users/79/")
		if !checkMatch(false, k, v, m) {
			b.Fatal(m)
		}
		k, v, m = MatchPatternV1("/api/users/{id}/foo", "/api/users/123456/foo")
		if !checkMatch(false, k, v, m) {
			b.Fatal(m)
		}
		res1 = k
		res2 = v
		res3 = m
	}
	_ = res1
	_ = res2
	_ = res3
}

// average: 2159 ns/op, 992 B/op, 20 allocs/op
func BenchmarkMatchPatternV2(b *testing.B) {
	var m int
	var k, v string
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		k, v, m = MatchPatternV2("/api/users/jobs/{jobID}", "/api/users/jobs/12")
		if !checkMatch(false, k, v, m) {
			b.Fatal(m)
		}
		k, v, m = MatchPatternV2("/api/users/jobs/{jobID}", "/api/users/jobs/1234")
		if !checkMatch(false, k, v, m) {
			b.Fatal(m)
		}
		k, v, m = MatchPatternV2("/api/users/{id}", "/api/users/12")
		if !checkMatch(false, k, v, m) {
			b.Fatal(m)
		}
		k, v, m = MatchPatternV2("/api/users/{id}", "/api/users/123456")
		if !checkMatch(false, k, v, m) {
			b.Fatal(m)
		}
		k, v, m = MatchPatternV2("/api/users/{id}/", "/api/users/79/")
		if !checkMatch(false, k, v, m) {
			b.Fatal(m)
		}
		k, v, m = MatchPatternV2("/api/users/{id}/foo", "/api/users/123456/foo")
		if !checkMatch(false, k, v, m) {
			b.Fatal(m)
		}
		res1 = k
		res2 = v
		res3 = m
	}
	_ = res1
	_ = res2
	_ = res3
}

// average: 2450 ns/op	2515 B/op	24 allocs/op
func BenchmarkMatch_parseV1(b *testing.B) {
	var ok bool
	var v url.Values
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		v, ok = parseV1("/api/users/jobs/:id", "/api/users/jobs/12")
		if v == nil || !ok {
			b.Fatal(v)
		}
		v, ok = parseV1("/api/users/jobs/:id", "/api/users/jobs/1234")
		if v == nil || !ok {
			b.Fatal(v)
		}
		v, ok = parseV1("/api/users/:id", "/api/users/12")
		if v == nil || !ok {
			b.Fatal(v)
		}
		v, ok = parseV1("/api/users/:id", "/api/users/123456")
		if v == nil || !ok {
			b.Fatal(v)
		}
		v, ok = parseV1("/api/users/:id/", "/api/users/79/")
		if v == nil || !ok {
			b.Fatal(v)
		}
		v, ok = parseV1("/api/users/:id/foo", "/api/users/123456/foo")
		if v == nil || !ok {
			b.Fatal(v)
		}
		res1 = v
		res2 = ok
	}
	_ = res1
	_ = res2
}

// average: 2450 ns/op	2515 B/op	24 allocs/op
func BenchmarkMatch_parseV2(b *testing.B) {
	var ok bool
	var v url.Values
	// b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		v, ok = parseV2("/api/users/jobs/:id", "/api/users/jobs/12")
		if v == nil || !ok {
			b.Fatal(v)
		}
		v, ok = parseV2("/api/users/jobs/:id", "/api/users/jobs/1234")
		if v == nil || !ok {
			b.Fatal(v)
		}
		v, ok = parseV2("/api/users/:id", "/api/users/12")
		if v == nil || !ok {
			b.Fatal(v)
		}
		v, ok = parseV2("/api/users/:id", "/api/users/123456")
		if v == nil || !ok {
			b.Fatal(v)
		}
		v, ok = parseV2("/api/users/:id/", "/api/users/79/")
		if v == nil || !ok {
			b.Fatal(v)
		}
		v, ok = parseV2("/api/users/:id/foo", "/api/users/123456/foo")
		if v == nil || !ok {
			b.Fatal(v)
		}
		res1 = v
		res2 = ok
	}
	_ = res1
	_ = res2
}
