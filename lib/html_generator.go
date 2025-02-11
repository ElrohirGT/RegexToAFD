package lib

import (
    "fmt"
    "os"
    "math"
    "strings"
)

var statePositions = map[string][2]int{
    "q0": {100, 200},
	"q1": {300, 100},
	"q2": {500, 200},
}

func (afd *AFD) ToSVG() string {
	svg := `<svg width="600" height="400" viewBox="0 0 600 400" xmlns="http://www.w3.org/2000/svg">`
	svg += `<rect width="100%" height="100%" fill="white"/>`

	radius := 30 // Radio de los nodos

	// Mapa para agrupar etiquetas por transición
	transitionLabels := make(map[[2]string][]string)

	// Recopilar transiciones
	for from, transitions := range afd.Transitions {
		for input, to := range transitions {
			key := [2]string{from, to}
			transitionLabels[key] = append(transitionLabels[key], input)
		}
	}

	// Dibujar los estados
	for state, pos := range statePositions {
		x, y := pos[0], pos[1]
		fill := "white"
		if afd.AcceptanceStates.Contains(state) {
			fill = "lightgreen"
		}
		svg += fmt.Sprintf(`<circle cx="%d" cy="%d" r="%d" stroke="black" stroke-width="2" fill="%s"/>`, x, y, radius, fill)
		svg += fmt.Sprintf(`<text x="%d" y="%d" font-size="16" text-anchor="middle" fill="black">%s</text>`, x, y+5, state)
	}

	// Dibujar las transiciones sin sobrescribir etiquetas
	for key, inputs := range transitionLabels {
		from, to := key[0], key[1]
		x1, y1 := statePositions[from][0], statePositions[from][1]
		x2, y2 := statePositions[to][0], statePositions[to][1]

		labels := strings.Join(inputs, ", ") // Combinar etiquetas

		if from == to {
			// 🔄 Dibujar loop más abajo
			loopRadius := 40
			offset := 0

			svg += fmt.Sprintf(
				`<path d="M %d %d C %d %d, %d %d, %d %d" stroke="black" stroke-width="2" fill="none" marker-end="url(#arrow)"/>`,
				x1, y1-radius-offset,   
				x1+loopRadius+10, y1-radius-loopRadius-offset-10, 
				x1-loopRadius-10, y1-radius-loopRadius-offset-10, 
				x1, y1-radius-offset,   
			)

			// Etiqueta en el loop
			svg += fmt.Sprintf(
				`<text x="%d" y="%d" font-size="14" fill="black">%s</text>`,
				x1, y1-radius-loopRadius-15, labels)

		} else {
			// 🔀 Dibujar transiciones normales
			dx := float64(x2 - x1)
			dy := float64(y2 - y1)
			dist := math.Sqrt(dx*dx + dy*dy)

			if dist > 0 {
				unitDx := dx / dist
				unitDy := dy / dist
				x1 += int(unitDx * float64(radius))
				y1 += int(unitDy * float64(radius))
				x2 -= int(unitDx * float64(radius))
				y2 -= int(unitDy * float64(radius))
			}

			// Línea de transición
			svg += fmt.Sprintf(
				`<line x1="%d" y1="%d" x2="%d" y2="%d" stroke="black" stroke-width="2" marker-end="url(#arrow)"/>`,
				x1, y1, x2, y2)

			// Etiqueta de transición combinada
			svg += fmt.Sprintf(
				`<text x="%d" y="%d" font-size="14" fill="black">%s</text>`,
				(x1+x2)/2, (y1+y2)/2 - 5, labels)
		}
	}

	// Definir flechas más grandes
	svg += `<defs>
		<marker id="arrow" markerWidth="15" markerHeight="15" refX="12" refY="6" orient="auto">
			<path d="M 0 0 L 12 6 L 0 12 z" fill="black"/>
		</marker>
	</defs>`

	svg += `</svg>`
	return svg
}



func GenerateHTML(svgContent, outputHTML string) error {
    htmlContent := fmt.Sprintf(`
	<!DOCTYPE html>
	<html lang="en">
	<head>
		<meta charset="UTF-8">
		<meta name="viewport" content="width=device-width, initial-scale=1.0">
		<title>Visualización AFD</title>
	</head>
	<body>
		<h2>Autómata Finito Determinista</h2>
		%s
	</body>
	</html>`, svgContent)

	return os.WriteFile(outputHTML, []byte(htmlContent), 0644)
}
