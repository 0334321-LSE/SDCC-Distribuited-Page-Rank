package utils

import (
	"PageRank/constants"
	"fmt"
	"strconv"
	"strings"
)

// ReducePageRank -> reduce function, aggregate results and send the update on channel
func ReducePageRank(pagerankShares <-chan string, result chan<- string, graphSize int) {

	//Map contains pageRanks for each node
	var aggregatePageRanks map[string]float64
	aggregatePageRanks = make(map[string]float64)
	var nodesList []string

	// Populate the map by summing marginal pr-value of associated node
	for line := range pagerankShares {
		//Divides row into nodeID, pageRankValue and outgoingLink
		lineParts := strings.Split(line, "\t")
		nodesList = append(nodesList, lineParts[0])
		linePagerankValue, _ := strconv.ParseFloat(lineParts[1], 64)
		outgoingLinks := strings.Split(lineParts[2], ",")
		//Add to cumulate pageRank the contribute of actual node
		for _, link := range outgoingLinks {
			aggregatePageRanks[link] += linePagerankValue
		}
	}
	//Now pageRank values are computed, print it to chanel
	for _, node := range nodesList {
		newPageRank := ((1.0 - constants.DampingFactor) / float64(graphSize)) + constants.DampingFactor*aggregatePageRanks[node]
		result <- fmt.Sprintf("%s\t%f", node, newPageRank)
		fmt.Println(fmt.Sprintf("%s\t%f", node, newPageRank))
	}

	fmt.Println("-----------------------------")

	return
}
