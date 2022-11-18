package netkit

import (
	"fmt"
	"net/url"
	"strings"
)

func locate(p []byte, t1, t2 byte) (int, int) {
	var n1, n2 int
	for i, c := range p {
		if n1 == 0 && c == t1 {
			n1 = i
		}
		if n2 == 0 && c == t2 {
			n2 = i
		}
		if n1 != 0 && n2 != 0 {
			break
		}
	}
	return n1, n2
}

func parseV1(p, path string) (url.Values, bool) {
	q := make(url.Values)
	var i, j int
	for i < len(path) {
		switch {
		case j >= len(p):
			if p != "/" && len(p) > 0 && p[len(p)-1] == '/' {
				return q, true
			}
			return nil, false
		case p[j] == ':':
			var name, val string
			var nextc byte
			name, nextc, j = match(p, isBoth, j+1)
			val, _, i = match(path, byteParse(nextc), i)
			q.Add(":"+name, val)
		case path[i] == p[j]:
			i++
			j++
		default:
			return nil, false
		}
	}
	if j != len(p) {
		return nil, false
	}
	return q, true
}

func parseV2(p, path string) (url.Values, bool) {
	q := make(url.Values)
	var i, j, k int
	for i < len(path) {
		switch {
		case j >= len(p):
			if p != "/" && len(p) > 0 && p[len(p)-1] == '/' {
				return q, true
			}
			return nil, false
		case p[j] == ':':
			var name, val string
			var next byte
			//
			// start matching for p
			k = j + 1
			for k < len(p) && isBoth(p[k]) {
				k++
			}
			if k < len(p) {
				next = p[k]
			}
			name = p[j+1 : k]
			j = k
			// done with p
			//
			// start matching for path
			k = i
			for k < len(path) && path[k] != next && path[k] != '/' {
				k++
			}
			if k < len(path) {
				next = path[k]
			}
			val = path[i:k]
			i = k
			// done with path
			//
			q.Add(":"+name, val)
		case path[i] == p[j]:
			i++
			j++
		default:
			return nil, false
		}
	}
	if j != len(p) {
		return nil, false
	}
	return q, true
}

func byteParse(b byte) func(byte) bool {
	return func(c byte) bool {
		return c != b && c != '/'
	}
}

// match path with registered handler
func match(s string, f func(byte) bool, i int) (matched string, next byte, j int) {
	j = i
	for j < len(s) && f(s[j]) {
		j++
	}
	if j < len(s) {
		next = s[j]
	}
	return s[i:j], next, j
}

func isAlpha(c byte) bool { return 'a' <= c && c <= 'z' || 'A' <= c && c <= 'Z' || c == '_' }
func isDigit(c byte) bool { return '0' <= c && c <= '9' }
func isBoth(c byte) bool {
	return ('a' <= c && c <= 'z' || 'A' <= c && c <= 'Z' || c == '_') || ('0' <= c && c <= '9')
}

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

func matchV0(p1, p2 string) (match int, key, val string) {
	// check for the presence of an opening path key
	// byte, and if we cannot find one, we will do
	// a simple path match.
	beg := strings.IndexByte(p1, '{')
	if beg == -1 {
		// check for a basic path match, and return
		if p1 == p2 {
			match = 1
		}
		return
	}
	// otherwise, we found an opening path key
	// byte, so next, we will check for a closed
	// path key byte.
	end := strings.IndexByte(p1, '}')
	if end == -1 {
		// we have not found a closed path
		// key byte, which means there must
		// be an error with p1, so a match
		// will be impossible at this point
		return
	}
	// we have found a closed path key byte, so
	// we should now extract the key from p1
	key = p1[beg+1 : end]
	// now on to do a path match verification check.
	if p1[:beg] == p2[:beg] {
		match = 1
	}
	// and then find the value that maps to our key.
	i := strings.IndexByte(p2[beg:], '/')
	if i == -1 {
		val = p2[beg:]
		return
	}
	val = p2[beg : beg+i]
	return
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
