package main

import (
	"bufio"
	"fmt"
	"os"

	l "github.com/ElrohirGT/RegexToAFD/lib"
)

func main() {

	// DEFAULT alphabet from program.
	// You can define a new one using: NewAlphabetFromString
	alph := DEFAULT_ALPHABET
	bst := new(l.BST)
	table := []*l.TableRow{}
	afd := new(l.AFD)

	file, err := os.Open("input.txt")

	if err != nil {
		fmt.Println("Error while opening file:", err)
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	words := []string{} // First element -> regex, Second element -> chain

	for scanner.Scan() {
		fmt.Println(scanner.Text())
		words = append(words, scanner.Text())
	}

	tokens := alph.ToPostfix(words[0])
	reverseSlice(tokens)
	bst.Insertion(tokens)

	table = l.ConvertTreeToTable(bst.List())

	afd = l.ConvertFromTableToAFD(table)
	afd = MinimizeAFD(afd)

	// Generar el SVG
	svg := afd.ToSVG()

	if afd.Derivation(words[1]) {
		fmt.Println("Cadena aceptada")
	} else {
		fmt.Println("Cadena rechazada")
	}

	// Guardar en un HTML
	htmlFile := "afd.html"
	if err := l.GenerateHTML(svg, htmlFile); err != nil {
		fmt.Println("Error al generar el HTML:", err)
	} else {
		fmt.Println("HTML generado:", htmlFile)
	}
}

func reverseSlice(arr []l.RX_Token) {
	n := len(arr)
	for i := 0; i < n/2; i++ {
		arr[i], arr[n-1-i] = arr[n-1-i], arr[i]
	}
}
