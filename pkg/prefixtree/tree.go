package prefixtree

type TrieNode struct {
	children map[rune]*TrieNode
	isEnd bool
}


type PrefixTrie struct {
	root *TrieNode
}

func NewTrie() *PrefixTrie {
	return &PrefixTrie{
		root: &TrieNode{children: make(map[rune]*TrieNode)},
	}
}

func (t *PrefixTrie) Insert(word string) {
	node := t.root
	for _, symb := range word {
		if _, ok := node.children[symb]; !ok {
			node.children[symb] = &TrieNode{children: make(map[rune]*TrieNode)}
		}
		node = node.children[symb]
	}
	node.isEnd = true
}

func (t *PrefixTrie) Search(word string) bool {
	node := t.root
	for _, symb := range word {
		next, ok := node.children[symb];
		if !ok {
			return false
		}
		node = next
	}
	return true
}
