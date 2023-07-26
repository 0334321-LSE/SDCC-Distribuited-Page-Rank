package utils

import (
	"fmt"
	"git.sr.ht/~sbinet/gg"
	"image/color"
	"math"
	"math/rand"
	"os"
	"strconv"
	"strings"
)

// BarabasiAlbertGraph genera un grafo utilizzando l'algoritmo Barabasi-Albert.
// nodes è il numero di nodi nel grafo, edgesToAttach rappresenta quanti archi collegare ad ogni nuovo nodo.
func BarabasiAlbertGraph(nodes, edgesToAttach int) map[int][]int {
	if nodes < 1 {
		return nil
	}

	graph := make(map[int][]int)
	degreeSum := make([]int, nodes)

	// Aggiungi il primo nodo.
	graph[0] = []int{}
	degreeSum[0] = 0

	for i := 1; i < nodes; i++ {
		graph[i] = []int{}

		for j := 0; j < edgesToAttach; j++ {
			// Scegli un nodo esistente in modo casuale, considerando il grado dei nodi.
			attachNode := rand.Intn(i)
			graph[i] = append(graph[i], attachNode)
			graph[attachNode] = append(graph[attachNode], i)

			// Aggiorna la somma dei gradi per i nodi coinvolti.
			degreeSum[i]++
			degreeSum[attachNode]++
		}
	}

	return graph
}

// plotGraph genera un grafico con nodi colorati con i numeri dei nodi all'interno dei cerchi e archi senza frecce. Lo salva come immagine PNG.
func plotGraph(graph map[int][]int) {
	rand.Seed(42) // Imposta un seed fisso per avere la stessa disposizione dei nodi ad ogni esecuzione.

	dc := gg.NewContext(1240, 1754) // Dimensioni in pixel per un foglio A4 verticale (210 x 297 mm a 300 dpi).

	dc.SetRGB(1, 1, 1)
	dc.Clear()

	// Distanza minima tra i nodi (puoi regolare questo valore per aumentare o diminuire la distanza).
	minNodeDistance := 150.0

	// Mappa per tenere traccia dei colori associati ai nodi.
	nodeColors := make(map[int]color.Color)

	// Mappa per tenere traccia delle posizioni dei nodi.
	nodePositions := make(map[int]struct{ x, y float64 })

	// Calcola il grado massimo tra tutti i nodi.
	var maxDegree int
	for _, neighbors := range graph {
		if degree := len(neighbors); degree > maxDegree {
			maxDegree = degree
		}
	}

	for i := 0; i < len(graph); i++ {
		node := i
		neighbors := graph[i]
		// Genera una posizione casuale per il nodo.
		for {
			x := rand.Float64() * 1240 // Valore casuale tra 0 e 1240 (larghezza del canvas).
			y := rand.Float64() * 1754 // Valore casuale tra 0 e 1754 (altezza del canvas).

			// Controlla che la posizione generata non sia troppo vicina ad altri nodi o non esca dai bordi del foglio A4.
			overlap := false
			if x < 40 || x > 1240-40 || y < 40 || y > 1754-40 {
				overlap = true
			}
			for _, pos := range nodePositions {
				dx := x - pos.x
				dy := y - pos.y
				distance := math.Sqrt(dx*dx + dy*dy)
				if distance < minNodeDistance {
					overlap = true
					break
				}
			}

			if !overlap {
				nodePositions[node] = struct{ x, y float64 }{x, y}
				break
			}
		}

		// Calcola il raggio proporzionale al grado del nodo.
		nodeDegree := len(neighbors)
		radius := float64(5 + nodeDegree*40/maxDegree) // Scala il raggio in base al grado massimo.

		// Genera un colore chiaro per il nodo.
		nodeColor := color.RGBA{
			R: uint8(rand.Intn(150) + 100),
			G: uint8(rand.Intn(150) + 100),
			B: uint8(rand.Intn(150) + 100),
			A: 255,
		}

		nodeColors[node] = nodeColor

		// Disegna un cerchio rappresentante il nodo.
		dc.SetColor(nodeColor)
		dc.DrawCircle(nodePositions[node].x, nodePositions[node].y, radius)
		dc.Fill()

		// Disegna il numero del nodo all'interno del cerchio.
		dc.SetRGB(0, 0, 0)
		dc.LoadFontFace("luxisr.ttf", 14) // Imposta la dimensione del testo a 14.
		dc.DrawStringAnchored(fmt.Sprintf("%d", node), nodePositions[node].x, nodePositions[node].y, 0.5, 0.5)
	}

	// Disegna gli archi tra i nodi.
	dc.SetRGB(0, 0, 0)
	for node, neighbors := range graph {
		x1 := nodePositions[node].x
		y1 := nodePositions[node].y
		for _, neighbor := range neighbors {
			x2 := nodePositions[neighbor].x
			y2 := nodePositions[neighbor].y

			// Calcola la distanza tra i due punti.
			dx := x2 - x1
			dy := y2 - y1
			distance := math.Sqrt(dx*dx + dy*dy)

			// Calcola i punti di intersezione tra l'arco e il bordo dei cerchi dei nodi.
			radius1 := float64(5 + len(graph[node])*40/maxDegree)
			radius2 := float64(5 + len(graph[neighbor])*40/maxDegree)
			intersectX1 := x1 + (dx/distance)*radius1
			intersectY1 := y1 + (dy/distance)*radius1
			intersectX2 := x2 - (dx/distance)*radius2
			intersectY2 := y2 - (dy/distance)*radius2

			// Disegna l'arco tra i nodi utilizzando i punti di intersezione.
			dc.DrawLine(intersectX1, intersectY1, intersectX2, intersectY2)
			dc.Stroke()
		}
	}

	// Salva l'immagine come file PNG.
	if err := dc.SavePNG("graph.png"); err != nil {
		fmt.Println("Errore durante il salvataggio del grafico:", err)
		return
	}

	fmt.Println("Grafico salvato come graph.png")
}

// writeAdjacencyListToFile scrive la lista di adiacenza del grafo su un file di testo.
func writeAdjacencyListToFile(graph map[int][]int, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	for i := 0; i < len(graph); i++ {
		node := i
		neighbors := graph[i]
		neighborStr := make([]string, len(neighbors))
		for i, neighbor := range neighbors {
			neighborStr[i] = strconv.Itoa(neighbor)
		}
		line := fmt.Sprintf("%d -> %s\n", node, strings.Join(neighborStr, ", "))
		_, err = file.WriteString(line)
		if err != nil {
			return err
		}
	}

	return nil
}

func CreateRandomGraph(numNodes int, edgesToAttach int) {

	// Genera il grafo utilizzando l'algoritmo Barabasi-Albert.
	graph := BarabasiAlbertGraph(numNodes, edgesToAttach)

	// Stampa la lista di adiacenza del grafo.
	fmt.Println("Lista di adiacenza:")
	for node, neighbors := range graph {
		fmt.Printf("%d -> %v\n", node, neighbors)
	}

	// Salva la lista di adiacenza su un file di testo.
	if err := writeAdjacencyListToFile(graph, "graph1.txt"); err != nil {
		fmt.Println("Errore durante il salvataggio della lista di adiacenza:", err)
		return
	}
	fmt.Println("Lista di adiacenza salvata come graph.txt")

	// Plotta e salva il grafico.
	plotGraph(graph)
}
