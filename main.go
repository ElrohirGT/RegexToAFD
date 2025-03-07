package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"

	l "github.com/ElrohirGT/RegexToAFD/lib"
)

func main() {
	// Disable loggin
	log.SetOutput(io.Discard)

	// DEFAULT alphabet from program.
	// You can define a new one using: NewAlphabetFromString
	alph := DEFAULT_ALPHABET
	bst := new(l.BST)
	table := []*l.TableRow{}
	afd := new(l.AFD)

	scanner := bufio.NewScanner(os.Stdin)
	words := []string{} // First element -> regex, Second element -> chain

	for scanner.Scan() {
		fmt.Println(scanner.Text())
		words = append(words, scanner.Text())
	}

	tokens := alph.ToPostfix(words[0])
	// reverseSlice(tokens)
	bst.Insertion(tokens)

	list := bst.List()
	table = l.ConvertTreeToTable(bst, list)

	afd = l.ConvertFromTableToAFD(table)
	afd = MinimizeAFD(afd)

	// (\.|\*)+([0-9]?)
	// .*|+01|2|3|4|5|6|7|8|9|_|.
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
