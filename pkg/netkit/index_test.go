package netkit

import (
	"bytes"
	"strings"
	"testing"
)

const position = 3
const repeat = 10000

var data = []byte(strings.Repeat("1", position) + "\t" + "dsahfksdfhkdsj")

// avg: 35000 ns/op (2nd place)
func BenchmarkIndexByte(b *testing.B) {
	for n := 0; n < b.N; n++ {
		for i := 0; i < repeat; i++ {
			if pos := bytes.IndexByte(data, '\t'); pos != position {
				b.Fatalf("Wrong position %d, %d expected", pos, position)
			}
		}
	}
}

// avg: 25000 ns/op (it gets inlined, 1st place)
func BenchmarkLoop(b *testing.B) {
	for n := 0; n < b.N; n++ {
		for i := 0; i < repeat; i++ {
			pos := -1
			for j, c := range data {
				if c == '\t' {
					pos = j
					break
				}
			}
			if pos != position {
				b.Fatalf("Wrong position %d, %d expected", pos, position)
			}
		}
	}
}

// avg: 45000 ns/op (it does not get inlined because of go:noinline below)
func BenchmarkLoopNoInline(b *testing.B) {
	for n := 0; n < b.N; n++ {
		for i := 0; i < repeat; i++ {
			pos := indexByte(data, '\t')
			if pos != position {
				b.Fatalf("Wrong position %d, %d expected", pos, position)
			}
		}
	}
}

//go:noinline
func indexByte(data []byte, ch byte) int {
	pos := -1
	for j, c := range data {
		if c == ch {
			pos = j
			break
		}
	}
	return pos
}
