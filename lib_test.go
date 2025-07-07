/*
	"libreria" di test per gli esami, attenzione a modificare questo file!
*/

package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"reflect"
	"strings"
	"testing"
)

var diffwidth int = 80

var HEADER string = "\n\n\n" + strings.Repeat(" ", (diffwidth-64)/2) + "___   ---   ===   ^^^   ***   TEST   ***   ^^^   ===   ---   ___"

/*func TestCompila(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		fmt.Println("*** errore nella lettura della directory corrente ***")
		t.Fail()
		return
	}

	nomeexe := path.Base(wd) // strippato diventa nome eseguibile
	nomego := nomeexe + ".go"
	//nometest := nomeexe + "_test.go"

	fexe, err := os.Stat(nomeexe)
	if fexe == nil {
		fmt.Println(HEADER)
		fmt.Println()
		fmt.Print("Verifico compilazione... ")
		fmt.Println("*** c'è qualche problema sul nome della directory o del file .go (non corrispondenza con le specifiche?) ***")
		fmt.Println(err)
		t.Fail()
		return
	}

	tExe := fexe.ModTime()
	//fmt.Println(nomeexe, tExe)

	fgo, _ := os.Stat(nomego)
	tGo := fgo.ModTime()

	if tGo.After(tExe) {
		fmt.Println("**************************************************************************")
		fmt.Println("*** ATTENZIONE! il sorgente non è stato compilato dopo le modifiche!!! ***")
		fmt.Println("**************************************************************************")
		t.Fail()
	}
}
*/

/* a tendere supportare anche mac (darwin)
 */
/*
func TestLinux(t *testing.T) {
	if runtime.GOOS != "linux" {
		fmt.Println(HEADER)
		fmt.Println()
		fmt.Print("Controllo sistema operativo...", runtime.GOOS)
		fmt.Println()
		fmt.Println("*************************************************")
		fmt.Println("* ATTENZIONE! sistema operativo NON supportato! *")
		fmt.Println("*************************************************")
	}
	//fmt.Println("--------------------------------------")
}
*/

/*
NUOVA base?
lancia stud e oracolo (il nome è hardcoded 'oracolo') e confronta

- BISOGNA lasciare nelle dir del tema il file `oracolo` eseguibile compilato dal nostro sorgente e impacchettarglielo nel tar
- si possono scrivere SIA i test classici che alcuni con questa `confronta`
- per ora vedi cancellaParole per un esempio d'uso
*/
func ConfrontaConOracolo(
	t *testing.T,
	progname string,
	filestdinput string, // se nome vuoto viene creato un contenuto a "NIENTE"
	args ...string) {

	fmt.Println(HEADER)
	fmt.Println()
	fmt.Println("[ Questo test confronta l'output studente con l'output atteso ]")
	fmt.Println()
	fmt.Println(">>> L'eseguibile da testare (", progname, ") deve essere stato compilato! <<<")
	fmt.Println()

	fileoracolo := "./oracolo"

	// fail fast se manca eseguibile studente
	if _, err := os.Stat(progname); err != nil {
		fmt.Println(">>> Manca eseguibile studente!")
		t.Fail()
		return
	}

	// fail fast se manca oracolo
	if _, err := os.Stat(fileoracolo); err != nil {
		fmt.Println(">>> Manca oracolo!")
		t.Fail()
		return
	}

	// prep i due stdin
	stdin1, filestdinputGlob, err1 := wrapStdin(filestdinput)
	stdin2, filestdinputGlob, err2 := wrapStdin(filestdinput)

	if err1 != nil || err2 != nil {
		t.Fail()
		return
	}

	//fmt.Println(stdin)

	// lancia oracolo e cattura out
	// chiama LanciaGenericaConFileInOutAtteso - NO, meglio rifarla qui così depreco le altre e questa diventa la base

	oracolo := exec.Command(fileoracolo, args...)
	oracolo.Stdin = stdin1
	oracoloout, err := oracolo.CombinedOutput()
	if err != nil {
		fmt.Printf("Attenzione! ORACOLO uscito con codice: %s\n\t>>> (non è un test fallito se si termina il programma con un esplicito os.Exit)\n", err)
	}
	//fmt.Println(">>> oracolo:", string(oracoloout))

	studente := exec.Command(progname, args...)
	studente.Stdin = stdin2
	studenteout, err := studente.CombinedOutput()
	if err != nil {
		fmt.Printf("Attenzione! STUDENTE uscito con codice: %s\n\t>>> (non è un test fallito se si termina il programma con un esplicito os.Exit)\n", err)
	}
	//fmt.Println(">>> studente:", string(studenteout))

	fmt.Printf("/// Argomenti a linea di comando:\t%s\n", args)
	fmt.Println()
	fmt.Printf("/// File per StdInput:\t%s\n", filestdinputGlob)
	fmt.Println()
	fmt.Println("### eseguo diff...")
	fmt.Println()
	out := Diff2strings(string(studenteout), "studente", string(oracoloout), "oracolo")

	//fmt.Printf("\n/// Output:\n%s\n", string(stdout))
	//fmt.Printf("\n/// Output atteso:\n%s\n", expectedOutString)

	//if string(stdout) != oracolo {
	if out != "" {

		fmt.Println(strings.Repeat("-", diffwidth))
		fmt.Println("[ studente ]", strings.Repeat(" ", diffwidth-12-8-10), "[ atteso ]")
		fmt.Println(strings.Repeat("-", diffwidth))
		fmt.Println(out)
		fmt.Println(strings.Repeat("-", diffwidth))

		fmt.Println(">>> FAIL! differisce da output atteso")
		t.Fail()
	}

	oracolo.Process.Kill()
	studente.Process.Kill()
	fmt.Println()
}

func wrapStdin(filestdinput string) (stdin io.Reader, filestdin string, err error) {
	filestdin = filestdinput

	if len(filestdinput) == 0 {
		stdin = strings.NewReader("NIENTE") // dummy
		filestdin = "<non era previsto input da stdin>"
	} else {
		stdin, err = os.Open(filestdinput)
		if err != nil {
			fmt.Println(">>> Non posso aprire file stdin:", filestdin)
		}
	}
	return
}

/*
TODO FATTORIZZARE LE VARIE LANCIA... ?
- in ingresso c'è stdin (potenzialmente vuoto) stringa/nomefile, args (potenzialmente vuoto)
- in output catturo stdout
- per testare confronto stdout con oracolo
- stdin posso usare direttamente i filename!!!
*/

/*
la base è tutto in forma di stringa (in e oracolo)
*/
func LanciaGenerica(
	t *testing.T,
	progname string,
	strinput string,
	oracolo string,
	verbose bool,
	args ...string) {
	subproc := exec.Command(progname, args...)
	subproc.Stdin = strings.NewReader(strinput)
	stdout, err := subproc.CombinedOutput() // invece di Run()

	if err != nil {
		fmt.Printf("Attenzione! Uscito con codice: %s\n\t>>> (non è un test fallito se si termina il programma con un esplicito os.Exit)\n", err)
	}

	out := Diff2strings(string(stdout), "studente", oracolo, "atteso")

	if out != "" {
		fmt.Printf("TEST %s -> FAIL\n", t.Name())
		if verbose {
			fmt.Printf("\n/// Argomenti a linea di comando:\n\t%s\n", args)
			fmt.Printf("\n/// StdInput:\n%s\n", strinput)

			fmt.Println(strings.Repeat("-", diffwidth))
			fmt.Println("[ studente ]", strings.Repeat(" ", diffwidth-12-8-10), "[ atteso ]")
			fmt.Println(strings.Repeat("-", diffwidth))
			fmt.Println(out)
			fmt.Println(strings.Repeat("-", diffwidth))
		}
		t.Fail()
	} else {
		fmt.Printf("TEST %s -> OK\n", t.Name())
	}

	subproc.Process.Kill()
	fmt.Println()
}

/*
si carica tutto in memoria... :(
*/
func LanciaGenericaConFileOutAtteso(
	t *testing.T,
	nomeProg string,
	strinput string,
	oracoloFilename string,
	verbose bool,
	args ...string) {

	content, err := ioutil.ReadFile(oracoloFilename)
	if err != nil {
		log.Fatal(err)
	}
	text := string(content)
	//fmt.Println(text)

	LanciaGenerica(t, nomeProg, strinput, text, verbose, args...)
}

func LanciaGenericaConFileInOutAtteso(
	t *testing.T,
	nomeProg string,
	inputFilename string,
	oracoloFilename string,
	verbose bool,
	args ...string) {

	input, err := ioutil.ReadFile(inputFilename)
	if err != nil {
		log.Fatal(err)
	}
	in := string(input)

	exout, err := ioutil.ReadFile(oracoloFilename)
	if err != nil {
		log.Fatal(err)
	}
	out := string(exout)
	//fmt.Println(text)

	LanciaGenerica(t, nomeProg, in, out, verbose, args...)
}

func Diff2files(fn1, fn2 string) (out string) {
	//cmd := exec.Command("diff", "-y", fn1, fn2)
	cmd := exec.Command("diff", "-y", "--suppress-common-lines", "-W", fmt.Sprint(diffwidth), fn1, fn2)
	//cmd := exec.Command("diff", "-y","-W 200", "--color=always", fn1, fn2) // verificare se c'è opzione color dappertutto

	var outbuf, errbuf strings.Builder // or bytes.Buffer
	cmd.Stdout = &outbuf
	cmd.Stderr = &errbuf
	cmd.Run()
	out = outbuf.String() + errbuf.String()
	return
}

/*
(inefficiente, lo so) crea due file temp, ci rovescia le due stringhe e chiama altro diff
*/
func Diff2strings(str1, l1, str2, l2 string) string {
	//TODO val la pena fattorizzare?
	tmpfile1, err1 := ioutil.TempFile("", l1+".*")
	if err1 != nil {
		log.Fatal(err1)
	}
	defer os.Remove(tmpfile1.Name()) // clean up
	if _, err1 := tmpfile1.Write([]byte(str1)); err1 != nil {
		log.Fatal(err1)
	}
	if err1 := tmpfile1.Close(); err1 != nil {
		log.Fatal(err1)
	}

	tmpfile2, err2 := ioutil.TempFile("", l2+".*")
	if err2 != nil {
		log.Fatal(err2)
	}
	defer os.Remove(tmpfile2.Name()) // clean up
	if _, err2 := tmpfile2.Write([]byte(str2)); err2 != nil {
		log.Fatal(err2)
	}
	if err2 := tmpfile2.Close(); err2 != nil {
		log.Fatal(err2)
	}

	return Diff2files(tmpfile1.Name(), tmpfile2.Name())
}

func CaptureOutput(fn interface{}, args ...interface{}) string {
	rescueStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	f := reflect.ValueOf(fn)
	if f.Type().NumIn() != len(args) {
		panic("incorrect number of parameters!")
	}
	inputs := make([]reflect.Value, len(args))
	for k, in := range args {
		inputs[k] = reflect.ValueOf(in)
	}
	f.Call(inputs)

	w.Close()
	out, _ := ioutil.ReadAll(r)
	os.Stdout = rescueStdout
	return string(out)
}

func checkOutput(output, expected string) bool {
	if output != expected {
		fmt.Printf("\n/// eseguo diff... \n")
		fmt.Println("\nESECUZIONE:\n<<<<<")
		fmt.Printf("%s", output)
		fmt.Println(">>>>")

		fmt.Println("\nATTESO:\n<<<<<")
		fmt.Printf("%s", expected)
		fmt.Println(">>>>")
		return false
	}
	return true

}
