package datastructures

const (
	trieLength = 27
)

type Trie interface {
	Insert(word string)
	Search(word string) map[string]int
}

type trieNode struct {
	children []*trieNode
	isWord   bool
	count    int
}

func newTrieNode() *trieNode {
	return &trieNode{
		children: make([]*trieNode, trieLength),
		isWord:   false,
		count:    0,
	}
}

type trie struct {
	root *trieNode
}

func NewTrie() Trie {
	n := newTrieNode()
	return &trie{n}
}

func (t *trie) Insert(word string) {
	temp := t.root
	for i := 0; i < len(word); i++ {
		idx := int(word[i] - 'a')
		if word[i] == '$' {
			idx = 26
		}
		if temp.children[idx] == nil {
			temp.children[idx] = newTrieNode()
		}
		temp = temp.children[idx]
		if i == len(word)-1 {
			temp.isWord = true
			temp.count++
		}
	}
}

func helpSearch(words map[string]int, current string, temp *trieNode) {
	if temp == nil {
		return
	}
	for i := 0; i < len(temp.children); i++ {
		if temp.children[i] == nil {
			continue
		}
		if i != 26 {
			current += string(i + 'a')
		}
		if temp.children[i].isWord {
			words[current] = temp.children[i].count
		}
		helpSearch(words, current, temp.children[i])
		if len(current) > 0 {
			current = current[:len(current)-1]
		}
	}
}

func (t *trie) Search(word string) map[string]int {
	words := make(map[string]int)
	temp := t.root

	for i := 0; i < len(word); i++ {
		idx := int(word[i] - 'a')
		if word[i] == '$' {
			idx = 26
		}
		if temp.children[idx] == nil {
			return nil
		}
		temp = temp.children[idx]
	}
	helpSearch(words, "", temp)

	return words
}
