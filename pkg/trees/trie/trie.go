package trie

import (
	"fmt"
	"strings"
)

const (
	trieNodeSize = 26 // math.MaxUint8
)

type trieNode struct {
	children [trieNodeSize]*trieNode
	terminal bool
}

func newTrieNode() *trieNode {
	node := new(trieNode)
	for i := range node.children {
		node.children[i] = nil
	}
	return node
}

type Trie struct {
	root  *trieNode
	depth int
}

func NewTrie() *Trie {
	return &Trie{
		root:  newTrieNode(),
		depth: 0,
	}
}

func (t *Trie) Insert(key string) bool {
	tmp := t.root
	for i := 0; i < len(key); i++ {
		if tmp.children[key[i]-'a'] == nil {
			tmp.children[key[i]-'a'] = newTrieNode()
		}
		tmp = tmp.children[key[i]-'a']
	}
	if tmp.terminal {
		tmp.terminal = false
	} else {
		tmp.terminal = true
	}
	return tmp.terminal
}

func (t *Trie) Search(key string) bool {
	tmp := t.root
	for i := 0; i < len(key); i++ {
		if tmp.children[key[i]-'a'] == nil {
			return false
		}
		tmp = tmp.children[key[i]-'a']
	}
	return tmp != nil && tmp.terminal
}

func (t *Trie) Delete(key string) bool {
	if t.root == nil {
		return false
	}
	tmp := t.root
	var result bool
	t.root = deleteRec(tmp, key, 0, &result)
	return result
}

func deleteRec(node *trieNode, key string, offset int, deleted *bool) *trieNode {
	if offset == len(key)-1 {
		//fmt.Printf(">> [deleting] at end of word: %q\n", key)
		if node.terminal {
			//fmt.Printf(">> [deleting] at terminal node: %q\n", key)
			node.terminal = false
			*deleted = true
			if !nodeHasChildren(node) {
				node = nil
			}
		}
		return node
	}
	prefix := key[offset]
	// fmt.Printf("prefix=%q, offset=%d, len(key)-1=%d\n", prefix, offset, len(key)-1)
	node.children[prefix-'a'] = deleteRec(node.children[prefix-'a'], key, offset+1, deleted)
	if *deleted && !nodeHasChildren(node) && !node.terminal {
		fmt.Printf(">> [deleting] success: %q\n", key)
		node = nil
	}
	return node
}

func nodeHasChildren(node *trieNode) bool {
	if node == nil {
		return false
	}
	for i := range node.children {
		if node.children[i] != nil {
			return true
		}
	}
	return false
}

func (t *Trie) String() string {
	s := new(strings.Builder)
	t.printRec(s, t.root, 0, make([]byte, 0), 0)
	return s.String()
}

func (t *Trie) printRec(s *strings.Builder, node *trieNode, n int, prefix []byte, length int) {
	newprefix := make([]byte, length+2)
	copy(newprefix, prefix[:length])
	newprefix[length+1] = byte(0)

	if node.terminal {
		fmt.Printf("child: %d\n", n)
		s.Write(prefix)
		s.WriteByte('\n')
	}

	for i := range node.children {
		if node.children[i] != nil {
			newprefix[length] = byte(i + 'a')
			// prefix[length] = byte(i + 'a')
			t.printRec(s, node.children[i], i, newprefix, length+1)
		}
	}
}
