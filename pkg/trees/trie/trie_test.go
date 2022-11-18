package trie

import (
	"fmt"
	"testing"
)

func InSet(k string, set []string) bool {
	for _, s := range set {
		if k == s {
			return true
		}
	}
	return false
}

func Test_NewTrieNode(t *testing.T) {
	n := newTrieNode()
	if n.terminal != false {
		t.Errorf("expexted: %v, got: %v\n", false, n.terminal)
	}
	for i := range n.children {
		if n.children[i] != nil {
			t.Errorf("expected: %v, got: %v\n", nil, n.children[i])
		}
	}
}

func TestNewTrie(t *testing.T) {
	var tt *Trie
	if tt != nil {
		t.Errorf("expected: %v, got: %v\n", nil, tt)
	}
	tt = NewTrie()
	if tt == nil {
		t.Errorf("expected: %v, got: %v\n", new(Trie), tt)
	}
}

func TestTrie_Insert(t *testing.T) {
	t.Logf("[INSERTING] a few words...")
	words := []string{"going", "and", "go", "a", "golang", "angler", "mango", "angle", "man"}
	tt := NewTrie()
	for _, word := range words {
		tt.Insert(word)
	}
	t.Logf("[SEARCHING] our words...")
	for _, word := range words {
		contains := tt.Search(word)
		if !contains {
			t.Errorf("expected: %v, got: %v\n", true, contains)
		}
		t.Logf("contains the word %q: %v\n", word, contains)
	}
}

func TestTrie_Search(t *testing.T) {
	t.Logf("[INSERTING] a few words...")
	words := []string{"going", "and", "go", "a", "golang", "angler", "mango", "angle", "man"}
	tt := NewTrie()
	for _, word := range words {
		tt.Insert(word)
	}
	t.Logf("[SEARCHING] our words...")
	for _, word := range words {
		contains := tt.Search(word)
		if !contains {
			t.Errorf("expected: %v, got: %v\n", true, contains)
		}
		t.Logf("contains the word %q: %v\n", word, contains)
	}
}

func TestTrie_Delete(t *testing.T) {
	t.Logf("[INSERTING] a few words...")
	words := []string{"going", "and", "go", "a", "golang", "angler", "mango", "angle", "man"}
	tt := NewTrie()
	for _, word := range words {
		tt.Insert(word)
	}
	t.Logf("[SEARCHING] our words...")
	for _, word := range words {
		contains := tt.Search(word)
		if !contains {
			t.Errorf("expected: %v, got: %v\n", true, contains)
		}
		t.Logf("contains the word %q: %v\n", word, contains)
	}
	t.Logf("[DELETING] a few words...")
	removeWords := []string{"go", "angler", "mango", "and"}
	for _, word := range removeWords {
		t.Logf("removing the word: %q\n", word)
		removed := tt.Delete(word)
		if !removed {
			t.Errorf("ecptected to remove %q, but couldn't\n", word)
		}
	}
	t.Logf("[SEARCHING] our words...")
	for _, word := range words {
		contains := tt.Search(word)
		if !contains && !InSet(word, removeWords) {
			t.Errorf("expected: %v, got: %v\n", true, contains)
		}
		t.Logf("contains the word %q: %v\n", word, contains)
	}
}

func TestTrie_String(t *testing.T) {
	words := []string{"going", "and", "go", "a", "golang", "angler", "mango", "angle", "man"}
	tt := NewTrie()
	for _, word := range words {
		tt.Insert(word)
	}
	fmt.Println(tt)
}
