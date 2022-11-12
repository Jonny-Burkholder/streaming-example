package netkit

import (
	"fmt"
	"strings"
)

func MatchPatternV1(p, path string) (string, string, int) {
	var key, val string
	var b, e, n, i, hc, m int
	for i < len(p) {
		if hc == 0 && m == 0 {
			n += strings.IndexByte(p, '{')
			if n == -1 {
				break
			}
			n++
			hc++
			// opening curly found, set beginning offset (b)
			b = n
		}
		if hc == 1 && m == 0 {
			// opening curly found, no match found
			if n <= len(p) && n <= len(path) {
				if p[:n-1] == path[:n-1] {
					m = n - 1 // // match found
				}
			}
		}
		if hc == 1 && m > 0 {
			n += strings.IndexByte(p[m:], '}')
			if n == -1 {
				break
			}
			n--
			hc++
			// ending curly found, set ending offset (e)
			e = n
		}
		if hc == 2 && m > 0 {
			// both curly braces found, and match established
			key = p[b:e]
			val = path[m:]
			break
		}
		i++
	}
	// we should have it all, return
	return key, val, m
}

func MatchPatternV2(p, path string) (string, string, int) {
	p1 := strings.Split(p, "/")
	p2 := strings.Split(path, "/")
	if len(p1) != len(p2) {
		return "", "", 0
	}
	var match int
	var key, val string
	for i := 0; i < len(p1); i++ {
		if p1[i] != p2[i] && i > 0 {
			if strings.IndexByte(p1[i], '{') != -1 || strings.IndexByte(p1[i], '}') != -1 {
				key = p1[i]
				val = p2[i]
				match = len(strings.Join(p1[:i-1], "/"))
			}
		}
	}
	return key, val, match
}

func MatchPatternV3(p, path string) {
	var n, i, hc int
	var key strings.Builder
	// n is the current character offset
	// c is the iteration count
	for i < 10 {
		fmt.Printf("n=%d, i=%d, hc=%d, p=%q, char=%c\n\n", n, i, hc, p[:n], p[n])

		// check to see if we are inside curly brackets
		if hc == 1 {
			n++                 // move to next character
			key.WriteByte(p[n]) // write char to key
			continue
		}

		// if we are not inside curly brackets, find next slash
		n += strings.IndexByte(p[n:], '/')
		if n == -1 {
			// no more slashes found
			break
		}
		n++ // found slash, skip, so we can do our other checks

		// check for left bracket
		if p[n] == '{' {
			n++                 // consume right curly brace
			hc++                // increase hit counter
			key.WriteByte(p[n]) // write char to key
			i++                 // increase iteration counter
			continue
		}

		// switch p[n] {
		// case '{':
		//
		// case '}':
		// default:
		// 	goto end
		// }
		i++ // increase the iteration count
	}
	// end:
	fmt.Println("done looping")
}
