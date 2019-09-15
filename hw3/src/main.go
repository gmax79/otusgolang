package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sort"
	"strings"
	"unicode"
)

// Pair struct to sorting slice by value
type Pair struct {
	Key   string
	Value int
}
type PairList []Pair

func (p PairList) Len() int           { return len(p) }
func (p PairList) Less(i, j int) bool { return p[i].Value < p[j].Value }
func (p PairList) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

func main() {
	// Read text from file
	file, err := os.Open("test.txt")
	if err != nil {
		log.Fatal(err)
	}

	// Create map with counters per word
	wordsMap := make(map[string]int)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		words := strings.Split(scanner.Text(), " ")
		for _, word := range words {
			if word = getWord(word); word == "" {
				continue
			}
			counter, _ := wordsMap[word]
			wordsMap[word] = counter + 1
		}
	}
	file.Close()

	// Find first 10 words
	// Translate map into slice first
	wordsList := make(PairList, 0, len(wordsMap))
	for k, v := range wordsMap {
		wordsList = append(wordsList, Pair{k, v})
	}
	// Sort slice by value from max to min
	sort.Sort(sort.Reverse(wordsList))

	// Get first 10 words
	count := len(wordsList)
	if count > 10 {
		count = 10
	}

	fmt.Printf("Found next %d words, as most offten:\n", count)
	for i := 0; i < count; i++ {
		fmt.Println(wordsList[i].Key, " = ", wordsList[i].Value)
	}
}

// getWord - filter function. Returns string with letters only.
// Trim not letters from begin and from end of paramter (points, commas etc)
func getWord(str string) string {
	for i, r := range str {
		if unicode.IsLetter(r) {
			str = str[i:]
			break
		}
	}
	for i, r := range str {
		if !unicode.IsLetter(r) {
			return str[:i]
		}
	}
	return str
}
