package main

import (
	"fmt"
	"sort"
	"strings"
)

type CasoTest struct {
	nome   string
	input  string
	atteso string
}

var prog = "./solution"
var verbose = true

// si assume che solo il penultimo comando nella stringa di input del test stampi su stdin
func eseguiTest(input string) string {
	lines := strings.Split(input, "\n")
	var d dizionario
	// si eseguono comandi che precedono la stampa su stdin
	for _, line := range lines[:len(lines)-1] {
		if line[0] == 'c' && len(line) == 1 {
			d = newDizionario()
		} else {
			esegui(d, line)
		}
	}

	// cattura stringa stampata su stdin
	comando := lines[len(lines)-1]
	cOper := strings.Fields(comando)[0]
	output := CaptureOutput(esegui, d, lines[len(lines)-1])

	// ordina output nel caso di comandi "p" e "s"
	if cOper == "p" || cOper == "s" {
		output = ordinaOutput(output[:len(output)-1]) // rimuove il carattere di newline finale
	}
	return output
}

func ordinaLineeTraDelimitatori(input string, delSx string, delDx string) string {
	lines := strings.Split(input, "\n")
	lines = lines[1 : len(lines)-1]
	if len(lines) == 0 {
		return input
	}

	sort.Strings(lines)
	sortedInput := strings.Join(lines, "\n")
	return fmt.Sprintf("%s\n%s\n%s", delSx, sortedInput, delDx)

}
func ordinaOutput(input string) string {
	return fmt.Sprintf("%s\n", ordinaLineeTraDelimitatori(input, "[", "]"))
}
