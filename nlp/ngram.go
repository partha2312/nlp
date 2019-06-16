package nlp

import (
	"fmt"
	"regexp"
	"sort"
	"strings"
	"time"

	lru "github.com/partha2312/nlp/datastructures/lru"
	trie "github.com/partha2312/nlp/datastructures/trie"
)

type NGram interface {
	ConstructNGrams(text string)
	Fetch(word string) []string
}

type nGram struct {
	trieBiGram  trie.Trie
	trieTriGram trie.Trie

	biOccurance  map[string]int
	triOccurance map[string]int

	biDone  chan (bool)
	triDone chan (bool)

	biFetch  chan (map[string]float32)
	triFetch chan (map[string]float32)

	lru lru.LRU
}

type result struct {
	word        string
	probability float32
}

type resultSorter []result

func (r resultSorter) Len() int           { return len(r) }
func (r resultSorter) Swap(i, j int)      { r[i], r[j] = r[j], r[i] }
func (r resultSorter) Less(i, j int) bool { return r[i].probability > r[j].probability }

func NewNGram() NGram {
	trieBiGram := trie.NewTrie()
	trieTriGram := trie.NewTrie()
	biOccurance := make(map[string]int)
	triOccurance := make(map[string]int)
	biDone := make(chan (bool))
	triDone := make(chan (bool))
	biFetch := make(chan (map[string]float32))
	triFetch := make(chan (map[string]float32))
	lru := lru.NewLRU(100)
	return &nGram{
		trieBiGram:   trieBiGram,
		trieTriGram:  trieTriGram,
		biOccurance:  biOccurance,
		triOccurance: triOccurance,
		biDone:       biDone,
		triDone:      triDone,
		biFetch:      biFetch,
		triFetch:     triFetch,
		lru:          lru,
	}
}

func (n *nGram) ConstructNGrams(text string) {
	go n.constructBiGramTrie(text)
	go n.constructTriGramTrie(text)
	<-n.biDone
	<-n.triDone
}

func (n *nGram) Fetch(words string) []string {
	words = sanitizeString(words)
	if len(words) == 0 {
		return nil
	}
	wordsArr := strings.Split(words, " ")

	last := wordsArr[len(wordsArr)-1]
	lastButOne := ""
	if len(wordsArr) > 1 {
		lastButOne = wordsArr[len(wordsArr)-2]
	}

	if val := n.lru.Get(last + lastButOne); val != nil {
		if result, ok := val.([]string); ok {
			return result
		}
	}

	start := time.Now()
	go n.biGramFetch(last)
	var triGramsProcessed map[string]float32
	if len(wordsArr) > 1 {
		go n.triGramFetch(last, lastButOne)
		triGramsProcessed = <-n.triFetch
	}
	biGramsProcessed := <-n.biFetch
	fmt.Println(fmt.Sprintf("fetch completed in %v", time.Since(start)))

	combined := postProcess(biGramsProcessed, triGramsProcessed)
	sort.Sort(resultSorter(combined))

	result := make([]string, 3)
	for idx, c := range combined {
		result[idx] = c.word
		if idx == 2 {
			break
		}
	}

	n.lru.Put(last+lastButOne, result)

	return result
}

func (n *nGram) biGramFetch(last string) {
	start := time.Now()
	defer fmt.Println(fmt.Sprintf("bi gram fetch completed in %v", time.Since(start)))
	biGrams := n.trieBiGram.Search(last + "$")
	n.biFetch <- process(biGrams, n.biOccurance[last])
}

func (n *nGram) triGramFetch(last, lastButOne string) {
	start := time.Now()
	defer fmt.Println(fmt.Sprintf("tri gram fetch completed in %v", time.Since(start)))
	triGrams := n.trieTriGram.Search(last + lastButOne + "$")
	n.triFetch <- process(triGrams, n.triOccurance[last+lastButOne])
}

func (n *nGram) constructBiGramTrie(text string) {
	fmt.Println("constructing bigram")
	start := time.Now()
	lines := strings.Split(text, ".")
	for _, line := range lines {
		line = sanitizeString(line)
		words := strings.Split(line, " ")
		for i := 0; i < len(words)-1; i++ {
			word1 := sanitizeString(words[i])
			word2 := sanitizeString(words[i+1])
			if len(word1) == 0 || len(word2) == 0 {
				continue
			}
			o1 := 0
			if cnt, ok := n.biOccurance[word1]; ok {
				o1 = cnt
			}
			o1++
			o2 := 0
			if cnt, ok := n.biOccurance[word1]; ok {
				o2 = cnt
			}
			o2++
			n.biOccurance[word1] = o1
			n.biOccurance[word2] = o2
			key := strings.Trim(fmt.Sprintf("%s$%s", word1, word2), " ")
			n.trieBiGram.Insert(key)
		}
	}
	fmt.Println(fmt.Sprintf("bigram constructed. took %v", time.Since(start)))
	n.biDone <- true
}

func (n *nGram) constructTriGramTrie(text string) {
	fmt.Println("constructing trigram")
	start := time.Now()
	lines := strings.Split(text, ".")
	for _, line := range lines {
		line = sanitizeString(line)
		words := strings.Split(line, " ")
		for i := 0; i < len(words)-2; i++ {
			word1 := sanitizeString(words[i])
			word2 := sanitizeString(words[i+1])
			word3 := sanitizeString(words[i+2])
			if len(word1) == 0 || len(word2) == 0 || len(word3) == 0 {
				continue
			}
			o1 := 0
			if cnt, ok := n.triOccurance[word1+word2]; ok {
				o1 = cnt
			}
			o1++
			o2 := 0
			if cnt, ok := n.triOccurance[word2+word3]; ok {
				o2 = cnt
			}
			o2++
			n.triOccurance[word1+word2] = o1
			n.triOccurance[word2+word3] = o2
			key := strings.Trim(fmt.Sprintf("%s%s$%s", word1, word2, word3), " ")
			n.trieTriGram.Insert(key)
		}
	}
	fmt.Println(fmt.Sprintf("trigram constructed. took %v", time.Since(start)))
	n.triDone <- true
}

func process(matrix map[string]int, count int) map[string]float32 {
	result := make(map[string]float32)
	for key, value := range matrix {
		result[key] = float32(value) / float32(count)
	}
	return result
}

func postProcess(bi map[string]float32, tri map[string]float32) []result {
	combined := make([]result, 0)
	for k1, v1 := range bi {
		if v2, ok := tri[k1]; ok {
			if v1 >= v2 {
				combined = append(combined, result{k1, v1})
			} else {
				combined = append(combined, result{k1, v2})
			}
			delete(tri, k1)
		} else {
			combined = append(combined, result{k1, v1})
		}
	}
	for k1, v1 := range tri {
		combined = append(combined, result{k1, v1})
	}
	return combined
}

func sanitizeString(word string) string {
	reg, err := regexp.Compile("[^a-zA-Z ]+")
	if err != nil {
		fmt.Println(err)
	}
	return reg.ReplaceAllString(strings.ToLower(strings.Trim(word, " ")), "")
}
