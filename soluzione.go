package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strings"
	"unicode"
)

const (
	cmdCreaDizionario = "C"
	cmdTermina        = "t"
	cmdCaricaFile     = "c"
	cmdStampaParole   = "P"
	cmdStampaSchemi   = "S"
	cmdInserisci      = "i"
	cmdElimina        = "e"
	cmdRicercaSchema  = "r"
	cmdDistanza       = "d"
)

type dizionario struct {
	words   map[string]struct{}
	schemes map[string]struct{}
}

func newDizionario() dizionario {
	return dizionario{
		words:   make(map[string]struct{}),
		schemes: make(map[string]struct{}),
	}
}

func isSchema(s string) bool {
	for _, r := range s {
		if unicode.IsUpper(r) {
			return true
		}
	}
	return false
}

func (d *dizionario) inserisci(w string) {
	if isSchema(w) {
		d.schemes[w] = struct{}{}
	} else {
		d.words[w] = struct{}{}
	}
}

func (d *dizionario) elimina(w string) {
	if isSchema(w) {
		delete(d.schemes, w)
	} else {
		delete(d.words, w)
	}
}

func (d *dizionario) stampaParole() {
	var wordsList []string
	for word := range d.words {
		wordsList = append(wordsList, word)
	}
	sort.Strings(wordsList)
	fmt.Println("[")
	for _, word := range wordsList {
		fmt.Println(word)
	}
	fmt.Println("]")
}

func (d *dizionario) stampaSchemi() {
	var schemesList []string
	for scheme := range d.schemes {
		schemesList = append(schemesList, scheme)
	}
	sort.Strings(schemesList)
	fmt.Println("[")
	for _, scheme := range schemesList {
		fmt.Println(scheme)
	}
	fmt.Println("]")
}

func (d *dizionario) carica(filename string) {
	content, err := os.ReadFile(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading file %s: %v\n", filename, err)
		return
	}
	fields := strings.Fields(string(content))
	for _, field := range fields {
		d.inserisci(field)
	}
}

func checkCompatibility(word, schema string) bool {
	if len(word) != len(schema) {
		return false
	}
	assignment := make(map[rune]rune)
	for i, schemaChar := range schema {
		wordChar := rune(word[i])
		if unicode.IsUpper(schemaChar) {
			if val, ok := assignment[schemaChar]; ok {
				if val != wordChar {
					return false
				}
			} else {
				assignment[schemaChar] = wordChar
			}
		} else if schemaChar != wordChar {
			return false
		}
	}
	return true
}

func (d *dizionario) ricerca(s string) {
	fmt.Printf("%s:[\n", s)
	var compatibleWords []string
	for word := range d.words {
		if checkCompatibility(word, s) {
			compatibleWords = append(compatibleWords, word)
		}
	}
	sort.Strings(compatibleWords)
	for _, w := range compatibleWords {
		fmt.Println(w)
	}
	fmt.Println("]")
}

func minInt(vals ...int) int {
	min := vals[0]
	for _, v := range vals[1:] {
		if v < min {
			min = v
		}
	}
	return min
}

func calculateEditDistance(s1, s2 string) int {
	m, n := len(s1), len(s2)
	dp := make([][]int, m+1)
	for i := range dp {
		dp[i] = make([]int, n+1)
	}
	for i := 0; i <= m; i++ {
		dp[i][0] = i
	}
	for j := 0; j <= n; j++ {
		dp[0][j] = j
	}
	for i := 1; i <= m; i++ {
		for j := 1; j <= n; j++ {
			cost := 1
			if s1[i-1] == s2[j-1] {
				cost = 0
			}
			dp[i][j] = minInt(
				dp[i-1][j]+1,
				dp[i][j-1]+1,
				dp[i-1][j-1]+cost,
			)
			if i > 1 && j > 1 && s1[i-1] == s2[j-2] && s1[i-2] == s2[j-1] {
				dp[i][j] = minInt(dp[i][j], dp[i-2][j-2]+1)
			}
		}
	}
	return dp[m][n]
}

func (d *dizionario) distanza(x, y string) {
	fmt.Println(calculateEditDistance(x, y))
}

func (d *dizionario) catena(x, y string) {
	if _, ok := d.words[x]; !ok {
		fmt.Println("non esiste")
		return
	}
	if _, ok := d.words[y]; !ok {
		fmt.Println("non esiste")
		return
	}
	if x == y {
		fmt.Println("(\n" + x + "\n)")
		return
	}
	queue := [][]string{{x}}
	visited := map[string]struct{}{x: {}}
	for len(queue) > 0 {
		path := queue[0]
		queue = queue[1:]
		last := path[len(path)-1]
		for word := range d.words {
			if _, seen := visited[word]; seen {
				continue
			}
			if calculateEditDistance(last, word) == 1 {
				newPath := append([]string{}, path...)
				newPath = append(newPath, word)
				if word == y {
					fmt.Println("(")
					for _, w := range newPath {
						fmt.Println(w)
					}
					fmt.Println(")")
					return
				}
				queue = append(queue, newPath)
				visited[word] = struct{}{}
			}
		}
	}
	fmt.Println("non esiste")
}

func esegui(d dizionario, s string) dizionario {
	parts := strings.Fields(s)
	if len(parts) == 0 {
		return d
	}
	switch parts[0] {
	case cmdCreaDizionario:
		return newDizionario()
	case cmdCaricaFile:
		if len(parts) == 2 {
			d.carica(parts[1])
		} else if len(parts) == 3 {
			d.catena(parts[1], parts[2])
		}
	case cmdStampaParole, strings.ToLower(cmdStampaParole):
		d.stampaParole()
	case cmdStampaSchemi, strings.ToLower(cmdStampaSchemi):
		d.stampaSchemi()
	case cmdInserisci:
		if len(parts) == 2 {
			d.inserisci(parts[1])
		}
	case cmdElimina:
		if len(parts) == 2 {
			d.elimina(parts[1])
		}
	case cmdRicercaSchema:
		if len(parts) == 2 {
			d.ricerca(parts[1])
		}
	case cmdDistanza:
		if len(parts) == 3 {
			d.distanza(parts[1], parts[2])
		}
	case cmdTermina:
		// handled in main
	default:
		fmt.Fprintf(os.Stderr, "Comando non riconosciuto: %s\n", parts[0])
	}
	return d
}

func main() {
	d := newDizionario()
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if len(line) == 0 {
			continue
		}
		if strings.HasPrefix(line, cmdTermina) {
			break
		}
		d = esegui(d, line)
	}
}
