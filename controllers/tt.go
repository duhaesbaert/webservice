package main

import (
	"fmt"
)

const (
	ALPHABETSIZE = 26
)

type trieNode struct {
	children [ALPHABETSIZE]*trieNode
	isEnd    bool
}
type trie struct {
	root *trieNode
}

func main() {
	trie := startTrie()

	fmt.Println("INSERIR PALAVRAS")
	wordsToAdd := []string{
		"aluno",
		"aluna",
		"professor",
		"cachorro",
		"correto",
		"teste",
	}
	fmt.Printf("Palavras")
	for _, word := range wordsToAdd {
		trie.insert(word)
		fmt.Print(" \"" + word + "\"")
	}
	fmt.Println(" inseridas.")

	wordsToSearch := []string{
		"aluno",
		"aluna",
		"alunos",
		"professor",
		"professores",
		"cachorro",
		"correto",
		"teste",
		"testes",
	}
	fmt.Println("----------------")

	fmt.Println("BUSCAR POR PALAVRAS")
	wordFound := make([]string, 0)
	wordNotFound := make([]string, 0)
	for _, word := range wordsToSearch {
		if found := trie.search(word); !found {
			wordNotFound = append(wordNotFound, word)
			continue
		}
		wordFound = append(wordFound, word)
	}

	fmt.Printf("Palavras")
	for _, wf := range wordFound {
		fmt.Printf(" \"" + wf + "\"")
	}
	fmt.Println(" encontradas na busca")

	fmt.Printf("Palavras")
	for _, wnf := range wordNotFound {
		fmt.Printf(" \"" + wnf + "\"")
	}
	fmt.Println(" não encontradas na busca")
	fmt.Println("----------------")

	fmt.Println("REMOVER PALAVRAS")
	wordsToRemove := []string{
		"aluno",
		"professor",
		"correto",
	}
	for _, word := range wordsToRemove {
		trie.remove(word)
	}

	fmt.Printf("Palavras")
	for _, word := range wordsToRemove {
		fmt.Print(" \"" + word + "\"")
	}
	fmt.Println(" removidas.")
	fmt.Println("----------------")

	fmt.Println("BUSCAR POR PALAVRAS APÓS REMOÇÃO")
	wordFound = make([]string, 0)
	wordNotFound = make([]string, 0)
	for _, word := range wordsToAdd {
		if found := trie.search(word); !found {
			wordNotFound = append(wordNotFound, word)
			continue
		}
		wordFound = append(wordFound, word)
	}

	fmt.Printf("Palavras")
	for _, wf := range wordFound {
		fmt.Printf(" \"" + wf + "\"")
	}
	fmt.Println(" encontradas na busca")

	fmt.Printf("Palavras")
	for _, wnf := range wordNotFound {
		fmt.Printf(" \"" + wnf + "\"")
	}
	fmt.Println(" não encontradas na busca")
	fmt.Println("----------------")
}

func startTrie() *trie {
	return &trie{root: &trieNode{}}
}

func (t *trie) insert(s string) {
	tn := t.root
	for i := 0; i < len(s); i++ {
		navIndex := s[i] - 'a'
		if tn.children[navIndex] == nil {
			tn.children[navIndex] = &trieNode{}
		}
		tn = tn.children[navIndex]
	}
	tn.isEnd = true
}

func (t *trie) search(s string) bool {
	tn := t.root
	for i := 0; i < len(s); i++ {
		navIndex := s[i] - 'a'
		if tn.children[navIndex] == nil {
			return false
		}
		tn = tn.children[navIndex]
	}

	return tn.isEnd
}

func (t *trie) remove(s string) bool {
	removeRec(t.root, s, 0)

	return !(t.search(s))
}

func removeRec(t *trieNode, s string, depth int) *trieNode {
	if t == nil {
		return nil
	}

	if depth == len(s) {
		if t.isEnd {
			t.isEnd = false
		}

		if isEmpty(*t) {
			t = nil
		}
		return t
	}
	navIndex := int(s[depth] - 'a')
	t.children[navIndex] = removeRec(t.children[navIndex], s, depth+1)

	if isEmpty(*t) && t.isEnd == false {
		t = nil
	}
	return t
}

func isEmpty(t trieNode) bool {
	for i := 0; i < ALPHABETSIZE; i++ {
		if t.children[i] != nil {
			return false
		}
	}
	return true
}
