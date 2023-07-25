package models

import (
	"PageRank/constants"
	"math"
	"strconv"
	"strings"
)

/*Structure of nodes*/
type Node struct {
	ID       string
	OutLinks []string
	PageRank float64
}

// Sum -> it sums PageRank contributes
func Sum(values <-chan float64) float64 {
	sum := 0.0
	for value := range values {
		sum += value
	}
	return sum
}

func Convergence(previous []float64, actual []float64) bool {
	for i := 0; i < len(previous); i++ {
		if !CheckConvergence(previous[i], actual[i]) {
			return false
		}
	}
	return true
}

// CheckConvergence -> check for each node if the difference between previous and
// next value of pageRank are similar or not
func CheckConvergence(previous float64, actual float64) bool {
	if math.Abs(previous-actual) < constants.Epsilon {
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

// UpdatePageRanks -> take the result from reducer and update PageRank values for each node
func UpdatePageRanks(nodeList []*Node, result <-chan string) {

	//Map contains pageRanks for each node
	var pageRankMap map[string]float64
	pageRankMap = make(map[string]float64)

	// Obtains page rank for each node
	for line := range result {
		//Divides row into nodeID, pageRankValue and save it into pageRankMap
		lineParts := strings.Split(line, "\t")
		node := lineParts[0]
		linePagerankValue, _ := strconv.ParseFloat(lineParts[1], 64)
		pageRankMap[node] = linePagerankValue
	}
	// Then update nodeList
	for _, node := range nodeList {
		node.PageRank = pageRankMap[node.ID]
	}
}
