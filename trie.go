package golite

import "strings"

const (
	WildKey string = ":WILD"
)

type Trie struct {
	root *Node
}

type Node struct {
	children     map[string]*Node
	controller   Controller
	hasWildChild bool
}

func NewTrie() *Trie {
	return &Trie{
		root: newNode(false),
	}
}

func newNode(hasWildChild bool) *Node {
	return &Node{
		children:     make(map[string]*Node),
		hasWildChild: hasWildChild,
	}
}

func isWildWord(word string) bool {
	return word == WildKey || strings.HasPrefix(word, ":")
}

func (t *Trie) Add(path string, controller Controller) {
	words := strings.Split(path, "/")
	node := t.root
	for _, w := range words {
		if isWildWord(w) {
			if node.hasWildChild {
				node = node.children[WildKey]
				continue
			}
			child := newNode(true)
			node.children[WildKey] = child
			node = child
		} else {
			if child, ok := node.children[w]; ok {
				node = child
				continue
			}
			child := newNode(false)
			node.children[w] = child
			node = child
		}
	}
	node.controller = controller
}

func (t *Trie) Get(path string) Controller {
	words := strings.Split(path, "/")
	node := t.root
	for i, w := range words {
		isLast := i == len(words)-1
		if isWildWord(w) {
			node = node.children[WildKey]
		} else {
			node = node.children[w]
		}
		if node == nil {
			return nil
		}
		if isLast && node.controller != nil {
			return node.controller
		}
	}
	return nil
}