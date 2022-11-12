package netkit

import (
	"fmt"
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
