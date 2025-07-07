package main

import (
	"bufio"
	"os"
	"strings"
	"testing"
)

// TestMain viene eseguito prima di tutti i test in questo pacchetto.
// Lo usiamo per creare il file di dizionario necessario per i test di `carica`.
func TestMain(m *testing.M) {
	// Contenuto del file dizionario per i test di 'carica'
	dizionarioContent := `casa
cosa
castoro
topo
anatra
coniglio
gatto
ratto
rotto`

	// Crea il file
	err := os.WriteFile("dizionario_carica.txt", []byte(dizionarioContent), 0644)
	if err != nil {
		// Se non possiamo creare il file, è inutile eseguire i test
		panic("Impossibile creare il file di test 'dizionario_carica.txt'")
	}

	// Esegui tutti i test
	exitCode := m.Run()

	// Rimuovi il file dopo che i test sono finiti
	os.Remove("dizionario_carica.txt")

	// Esci con il codice di uscita dei test
	os.Exit(exitCode)
}

// TestRunner legge i casi di test da suite_test.txt e li esegue.
func TestRunner(t *testing.T) {
	file, err := os.Open("suite_test.txt")
	if err != nil {
		t.Fatalf("Impossibile aprire il file della suite di test 'suite_test.txt': %v", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var currentTestName string
	var inputBuilder, attesoBuilder strings.Builder
	var isReadingInput, isReadingAtteso bool

	runTest := func() {
		if currentTestName != "" {
			// Pulisci e normalizza l'output atteso
			// Rimuove spazi bianchi iniziali/finali e normalizza i newline
			atteso := strings.TrimSpace(attesoBuilder.String()) + "\n"
			if atteso == "\n" { // Caso di output vuoto
				atteso = ""
			} else if strings.HasSuffix(atteso, "]\n") {
				// Non aggiungere un ulteriore newline se finisce con la parentesi
			} else if strings.Count(atteso, "\n") == 1 && len(strings.TrimSpace(atteso)) > 0 && !strings.Contains(atteso, "[") {
				// Output di singola linea come 'distanza' o 'non esiste', già con newline
			} else {
				atteso = strings.TrimSpace(attesoBuilder.String()) + "\n"
			}

			input := strings.TrimSpace(inputBuilder.String())

			// Esegui il sub-test
			t.Run(currentTestName, func(t *testing.T) {
				output := eseguiTest(input)

				// Normalizza l'output per il confronto
				normalizedOutput := strings.TrimSpace(output) + "\n"
				if output == "" {
					normalizedOutput = ""
				} else if strings.HasSuffix(output, "]\n") {
					normalizedOutput = output
				}

				// Confronto
				if normalizedOutput != atteso {
					t.Errorf("\n--- TEST FALLITO: %s ---\nInput:\n%s\n\nOutput Eseguito:\n<<<<<\n%s>>>>>\n\nOutput Atteso:\n<<<<<\n%s>>>>>",
						currentTestName, input, output, strings.TrimSuffix(atteso, "\n"))
				}
			})
		}
		// Reset per il prossimo test
		currentTestName = ""
		inputBuilder.Reset()
		attesoBuilder.Reset()
		isReadingInput = false
		isReadingAtteso = false
	}

	for scanner.Scan() {
		line := scanner.Text()

		switch {
		case strings.HasPrefix(line, "nome:"):
			currentTestName = strings.TrimSpace(strings.TrimPrefix(line, "nome:"))
		case line == "--INPUT--":
			isReadingInput = true
			isReadingAtteso = false
		case line == "--ATTESO--":
			isReadingInput = false
			isReadingAtteso = true
		case line == "--END--":
			runTest()
		case line == "--TEST--":
			// Ignora, è solo un separatore
		default:
			if isReadingInput {
				inputBuilder.WriteString(line + "\n")
			} else if isReadingAtteso {
				attesoBuilder.WriteString(line + "\n")
			}
		}
	}

	// Esegui l'ultimo test nel file
	runTest()

	if err := scanner.Err(); err != nil {
		t.Fatalf("Errore durante la lettura del file di test: %v", err)
	}
}
