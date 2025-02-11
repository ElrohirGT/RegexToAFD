package main

import (
    "fmt"

	l "github.com/ElrohirGT/RegexToAFD/lib"
)

func main() {
    afd := l.AFD{
		InitialState: "q0",
		Transitions: map[l.AFDState]map[l.AlphabetInput]l.AFDState{
			"q0": {"0": "q1", "1": "q0"},
			"q1": {"0": "q0", "1": "q2"},
			"q2": {"0": "q2", "1": "q3"},
            "q3": {"0": "q1", "1": "q2"},
		},
		AcceptanceStates: l.Set[string]{"2": struct{}{}},	
    }
	// Generar el SVG
	svg := afd.ToSVG()

	// Guardar en un HTML
	htmlFile := "afd.html"
	if err := l.GenerateHTML(svg, htmlFile); err != nil {
		fmt.Println("Error al generar el HTML:", err)
	} else {
		fmt.Println("HTML generado:", htmlFile)
	}
}
