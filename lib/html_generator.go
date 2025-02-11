package lib

import (
    "fmt"
    "os"
    "math"
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

	// Dibujar las transiciones
	for from, transitions := range afd.Transitions {
		for input, to := range transitions {
			x1, y1 := statePositions[from][0], statePositions[from][1]
			x2, y2 := statePositions[to][0], statePositions[to][1]

			// Calcular la dirección del vector
			dx := float64(x2 - x1)
			dy := float64(y2 - y1)
			dist := math.Sqrt(dx*dx + dy*dy)

			if dist > 0 {
				// Normalizar el vector de dirección
				unitDx := dx / dist
				unitDy := dy / dist

				// Ajustar los puntos para que la línea comience y termine en los bordes del nodo
				x1 = x1 + int(unitDx*float64(radius))
				y1 = y1 + int(unitDy*float64(radius))
				x2 = x2 - int(unitDx*float64(radius))
				y2 = y2 - int(unitDy*float64(radius))
			}

			// Dibujar línea de transición con flecha
			svg += fmt.Sprintf(
				`<line x1="%d" y1="%d" x2="%d" y2="%d" stroke="black" stroke-width="2" marker-end="url(#arrow)"/>`,
				x1, y1, x2, y2)

			// Etiqueta de transición
			svg += fmt.Sprintf(
				`<text x="%d" y="%d" font-size="14" fill="black">%s</text>`,
				(x1+x2)/2, (y1+y2)/2, input)
		}
	}

	// Definir flechas
	svg += `<defs>
		<marker id="arrow" markerWidth="10" markerHeight="10" refX="10" refY="5" orient="auto">
			<path d="M 0 0 L 10 5 L 0 10 z" fill="black"/>
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
