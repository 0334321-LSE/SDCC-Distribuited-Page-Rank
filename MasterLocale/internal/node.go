package internal

import (
	"Master/constants"
	"fmt"
	"git.sr.ht/~sbinet/gg"
	"image/color"
	"math"
	"math/rand"
	"strconv"
)

// Node -> Structure of nodes
type Node struct {
	ID       int
	OutLinks []int
	PageRank float64
}

// Convergence -> check for each node if the difference between previous and
// next value of pageRank is less than epsilon constant, otherwise isn't converged
func Convergence(previous []float64, actual []float64) bool {
	for i := 0; i < len(previous); i++ {
		if !CheckConvergence(previous[i], actual[i]) {
			return false
		}
	}
	return true
}

// CheckConvergence -> check for a node if the difference between previous and
// next value of pageRank is less than epsilon constant or not
func CheckConvergence(previous float64, actual float64) bool {
	var config constants.Config
	constants.ReadJsonConfig(&config)
	if float32(math.Abs(previous-actual)) < config.Epsilon {
		return true
	}
	return false
}

// ListOfPageRank -> return a list containing page rank values
func ListOfPageRank(list []*Node) []float64 {
	var pageRankList []float64
	for _, n := range list {
		pageRankList = append(pageRankList, n.PageRank)
	}
	return pageRankList
}

// GetOutLinks -> return out-links of a node in int32 format
func GetOutLinks(node *Node) []int32 {
	var outLinks []int32
	for _, link := range node.OutLinks {
		outLinks = append(outLinks, int32(link))
	}
	return outLinks
}

// PlotGraphByPageRank -> plot-graph representation where nodes with greater PageRank are bigger
func PlotGraphByPageRank(nodes []*Node) {
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

	// Calcola il PageRank massimo tra tutti i nodi.
	var maxPageRank float64
	for _, node := range nodes {
		if node.PageRank > maxPageRank {
			maxPageRank = node.PageRank
		}
	}

	for _, node := range nodes {
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
				nodePositions[node.ID] = struct{ x, y float64 }{x, y}
				break
			}
		}

		// Calcola il raggio proporzionale al PageRank del nodo.
		radius := float64(5 + int(node.PageRank*40/maxPageRank)) // Scala il raggio in base al PageRank massimo.

		// Genera un colore chiaro per il nodo.
		nodeColor := color.RGBA{
			R: uint8(rand.Intn(150) + 100),
			G: uint8(rand.Intn(150) + 100),
			B: uint8(rand.Intn(150) + 100),
			A: 255,
		}

		nodeColors[node.ID] = nodeColor

		// Disegna un cerchio rappresentante il nodo.
		dc.SetColor(nodeColor)
		dc.DrawCircle(nodePositions[node.ID].x, nodePositions[node.ID].y, radius)
		dc.Fill()

		// Disegna il numero del nodo all'interno del cerchio.
		dc.SetRGB(0, 0, 0)
		dc.LoadFontFace("luxisr.ttf", 14) // Imposta la dimensione del testo a 14.
		dc.DrawStringAnchored(strconv.Itoa(node.ID), nodePositions[node.ID].x, nodePositions[node.ID].y, 0.5, 0.5)
	}

	// Disegna gli archi tra i nodi.
	dc.SetRGB(0, 0, 0)
	for _, node := range nodes {
		x1 := nodePositions[node.ID].x
		y1 := nodePositions[node.ID].y
		for _, outLink := range node.OutLinks {
			// Trova il nodo collegato tramite l'OutLink.
			neighbor := findNodeByID(nodes, outLink)
			if neighbor == nil {
				continue
			}

			x2 := nodePositions[neighbor.ID].x
			y2 := nodePositions[neighbor.ID].y

			// Calcola la distanza tra i due punti.
			dx := x2 - x1
			dy := y2 - y1
			distance := math.Sqrt(dx*dx + dy*dy)

			// Calcola i punti di intersezione tra l'arco e il bordo dei cerchi dei nodi.
			radius1 := float64(5 + int(node.PageRank*40/maxPageRank))
			radius2 := float64(5 + int(neighbor.PageRank*40/maxPageRank))
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
	if err := dc.SavePNG("./output/PR-graph.png"); err != nil {
		fmt.Println("Errore durante il salvataggio del grafico:", err)
		return
	}

	fmt.Println("Graph with page ranks saved as PR-Graph.png")
}

// findNodeByID -> as the name says ...
func findNodeByID(nodes []*Node, id int) *Node {
	for _, node := range nodes {
		if node.ID == id {
			return node
		}
	}
	return nil
}
