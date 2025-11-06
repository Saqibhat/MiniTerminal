package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// bygger en enkel shell som aksepterer
// kommandoer som kjører visse OS-funksjoner eller programmer. For OS
// funksjoner, se golang's innebygde OS og ioutil pakker.
//
// Shellen skal implementeres gjennom en kommandolinje
// applikasjon; som lar brukeren utføre alle funksjonene
// spesifisert i oppgaven. Info som [sti] er kommandoargumenter
//
// Viktig: Promptet til shellen skal skrive ut gjeldende mappe.
// For eksempel, noe som dette:
//   /Users/meling/Dropbox/work/opsys/2020/meling-stud-labs/lab3>
//
// Vi foreslår å bruke et mellomrom etter > symbolet.

// Terminal inneholder
type Terminal struct {
	gjeldendeMappe string
	historikk      []string
}

// NyTerminal oppretter en ny terminal-instans
func NyTerminal() *Terminal {
	gjeldendeMappe, _ := os.Getwd()
	return &Terminal{
		gjeldendeMappe: gjeldendeMappe,
		historikk:      make([]string, 0),
	}
}

// Utfør utfører en gitt kommando
func (t *Terminal) Utfør(kommando string) {
	kommando = strings.TrimSpace(kommando)
	if kommando == "" {
		return
	}

	// Legg til i historikk
	t.historikk = append(t.historikk, kommando)

	// Parse kommando og argumenter
	deler := strings.Fields(kommando)
	cmd := deler[0]
	args := deler[1:]

	switch cmd {
	case "avslutt":
		t.håndterAvslutt()
	case "cd":
		t.håndterCd(args)
	case "ls":
		t.håndterLs(args)
	case "mkdir":
		t.håndterMkdir(args)
	case "rm":
		t.håndterRm(args)
	case "opprett":
		t.håndterOpprett(args)
	case "cat":
		t.håndterCat(args)
	case "hjelp":
		t.håndterHjelp(args)
	case "historikk":
		t.håndterHistorikk()
	case "head":
		t.håndterHead(args)
	case "tail":
		t.håndterTail(args)
	default:
		fmt.Printf("Kommando ikke funnet: %s\n", cmd)
		fmt.Println("Skriv 'hjelp' for tilgjengelige kommandoer.")
	}
}

func (t *Terminal) håndterAvslutt() {
	fmt.Println("Ha det!")
	os.Exit(0)
}

func (t *Terminal) håndterCd(args []string) {
	if len(args) == 0 {
		fmt.Println("cd: mangler sti-argument")
		return
	}

	sti := args[0]
	if sti == ".." {
		t.gjeldendeMappe = filepath.Dir(t.gjeldendeMappe)
	} else if filepath.IsAbs(sti) {
		t.gjeldendeMappe = sti
	} else {
		t.gjeldendeMappe = filepath.Join(t.gjeldendeMappe, sti)
	}

	// Verifiser at mappen eksisterer og oppdater hvis gyldig
	if info, err := os.Stat(t.gjeldendeMappe); err != nil {
		fmt.Printf("cd: %v\n", err)
		// Gå tilbake til forrige mappe
		if wd, err := os.Getwd(); err == nil {
			t.gjeldendeMappe = wd
		}
	} else if !info.IsDir() {
		fmt.Printf("cd: %s er ikke en mappe\n", t.gjeldendeMappe)
		if wd, err := os.Getwd(); err == nil {
			t.gjeldendeMappe = wd
		}
	} else {
		os.Chdir(t.gjeldendeMappe)
	}
}

func (t *Terminal) håndterLs(args []string) {
	sti := t.gjeldendeMappe
	if len(args) > 0 {
		sti = args[0]
		if !filepath.IsAbs(sti) {
			sti = filepath.Join(t.gjeldendeMappe, sti)
		}
	}

	filer, err := ioutil.ReadDir(sti)
	if err != nil {
		fmt.Printf("ls: %v\n", err)
		return
	}

	for _, fil := range filer {
		if fil.IsDir() {
			fmt.Printf("%-20s <MAPPE>\n", fil.Name())
		} else {
			fmt.Printf("%-20s %d bytes\n", fil.Name(), fil.Size())
		}
	}
}

func (t *Terminal) håndterMkdir(args []string) {
	if len(args) == 0 {
		fmt.Println("mkdir: mangler mappenavn")
		return
	}

	sti := args[0]
	if !filepath.IsAbs(sti) {
		sti = filepath.Join(t.gjeldendeMappe, sti)
	}

	err := os.MkdirAll(sti, 0755)
	if err != nil {
		fmt.Printf("mkdir: %v\n", err)
	} else {
		fmt.Printf("Mappe opprettet: %s\n", sti)
	}
}

func (t *Terminal) håndterRm(args []string) {
	if len(args) == 0 {
		fmt.Println("rm: mangler fil/mappenavn")
		return
	}

	rekursiv := false
	startIdx := 0

	// Sjekk for -r flagg
	if len(args) > 0 && args[0] == "-r" {
		rekursiv = true
		startIdx = 1
	}

	if startIdx >= len(args) {
		fmt.Println("rm: mangler fil/mappenavn")
		return
	}

	sti := args[startIdx]
	if !filepath.IsAbs(sti) {
		sti = filepath.Join(t.gjeldendeMappe, sti)
	}

	var err error
	if rekursiv {
		err = os.RemoveAll(sti)
	} else {
		err = os.Remove(sti)
	}

	if err != nil {
		fmt.Printf("rm: %v\n", err)
	} else {
		fmt.Printf("Fjernet: %s\n", sti)
	}
}

func (t *Terminal) håndterOpprett(args []string) {
	if len(args) == 0 {
		fmt.Println("opprett: mangler filnavn")
		return
	}

	sti := args[0]
	if !filepath.IsAbs(sti) {
		sti = filepath.Join(t.gjeldendeMappe, sti)
	}

	fil, err := os.Create(sti)
	if err != nil {
		fmt.Printf("opprett: %v\n", err)
		return
	}
	defer fil.Close()

	fmt.Printf("Fil opprettet: %s\n", sti)
}

func (t *Terminal) håndterCat(args []string) {
	if len(args) == 0 {
		fmt.Println("cat: mangler filnavn")
		return
	}

	sti := args[0]
	if !filepath.IsAbs(sti) {
		sti = filepath.Join(t.gjeldendeMappe, sti)
	}

	innhold, err := ioutil.ReadFile(sti)
	if err != nil {
		fmt.Printf("cat: %v\n", err)
		return
	}

	fmt.Print(string(innhold))
}

func (t *Terminal) håndterHead(args []string) {
	linjer := 10 // standard
	var filnavn string

	if len(args) == 0 {
		fmt.Println("head: mangler filnavn")
		return
	}

	if len(args) == 1 {
		filnavn = args[0]
	} else if len(args) == 2 {
		if n, err := strconv.Atoi(args[0]); err == nil {
			linjer = n
			filnavn = args[1]
		} else {
			filnavn = args[0]
		}
	}

	if !filepath.IsAbs(filnavn) {
		filnavn = filepath.Join(t.gjeldendeMappe, filnavn)
	}

	innhold, err := ioutil.ReadFile(filnavn)
	if err != nil {
		fmt.Printf("head: %v\n", err)
		return
	}

	filLinjer := strings.Split(string(innhold), "\n")
	slutt := len(filLinjer)
	if linjer < slutt {
		slutt = linjer
	}

	for i := 0; i < slutt; i++ {
		fmt.Println(filLinjer[i])
	}
}

func (t *Terminal) håndterTail(args []string) {
	linjer := 10 // standard
	var filnavn string

	if len(args) == 0 {
		fmt.Println("tail: mangler filnavn")
		return
	}

	if len(args) == 1 {
		filnavn = args[0]
	} else if len(args) == 2 {
		if n, err := strconv.Atoi(args[0]); err == nil {
			linjer = n
			filnavn = args[1]
		} else {
			filnavn = args[0]
		}
	}

	if !filepath.IsAbs(filnavn) {
		filnavn = filepath.Join(t.gjeldendeMappe, filnavn)
	}

	innhold, err := ioutil.ReadFile(filnavn)
	if err != nil {
		fmt.Printf("tail: %v\n", err)
		return
	}

	filLinjer := strings.Split(string(innhold), "\n")
	start := len(filLinjer) - linjer
	if start < 0 {
		start = 0
	}

	for i := start; i < len(filLinjer); i++ {
		if i < len(filLinjer)-1 || filLinjer[i] != "" {
			fmt.Println(filLinjer[i])
		}
	}
}

func (t *Terminal) håndterHistorikk() {
	for i, cmd := range t.historikk {
		fmt.Printf("%d: %s\n", i+1, cmd)
	}
}

func (t *Terminal) håndterHjelp(args []string) {
	if len(args) == 0 {
		fmt.Println("Tilgjengelige kommandoer:")
		fmt.Println("  avslutt        - avslutt programmet")
		fmt.Println("  cd [sti]       - bytt mappe til angitt sti")
		fmt.Println("  ls [sti]       - list elementer i gjeldende eller angitt mappe")
		fmt.Println("  mkdir [sti]    - opprett en mappe")
		fmt.Println("  rm [sti]       - fjern en fil eller mappe")
		fmt.Println("  rm -r [sti]    - fjern mappe og innhold rekursivt")
		fmt.Println("  opprett [sti]  - opprett en fil")
		fmt.Println("  cat [fil]      - vis innholdet i en fil")
		fmt.Println("  head [n] fil   - vis første n linjer i fil (standard 10)")
		fmt.Println("  tail [n] fil   - vis siste n linjer i fil (standard 10)")
		fmt.Println("  historikk      - vis kommandohistorikk")
		fmt.Println("  hjelp [kommando] - vis hjelp for spesifikk kommando")
	} else {
		cmd := args[0]
		switch cmd {
		case "cd":
			fmt.Println("cd [sti] - Bytt gjeldende mappe til angitt sti")
			fmt.Println("  Bruk 'cd ..' for å gå opp ett nivå")
		case "ls":
			fmt.Println("ls [sti] - List filer og mapper")
			fmt.Println("  Viser størrelse for filer og <MAPPE> for mapper")
		case "rm":
			fmt.Println("rm [sti] - Fjern fil eller tom mappe")
			fmt.Println("rm -r [sti] - Fjern mappe og alt innhold rekursivt")
		case "head":
			fmt.Println("head [n] fil - Vis første n linjer i fil (standard 10)")
		case "tail":
			fmt.Println("tail [n] fil - Vis siste n linjer i fil (standard 10)")
		default:
			fmt.Printf("Ingen detaljert hjelp tilgjengelig for '%s'\n", cmd)
		}
	}
}

func (t *Terminal) hentPrompt() string {
	return fmt.Sprintf("%s> ", t.gjeldendeMappe)
}

func main() {
	fmt.Println("Velkommen til terminalen!")

	terminal := NyTerminal()
	leser := bufio.NewReader(os.Stdin)

	for {
		fmt.Print(terminal.hentPrompt())
		inndata, err := leser.ReadString('\n')
		if err != nil {
			fmt.Printf("Feil ved lesing av inndata: %v\n", err)
			continue
		}

		terminal.Utfør(strings.TrimSpace(inndata))
	}
}
