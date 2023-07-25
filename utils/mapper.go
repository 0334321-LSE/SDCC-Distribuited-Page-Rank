package utils

import (
	"PageRank/models"
	"fmt"
	"strings"
)

/*type Mapper struct {
	node Node
}*/

// MapPageRank -> Map function to evaluate each node contributes
func MapPageRank(node *models.Node, result chan<- string) {
	numOutLinks := float64(len(node.OutLinks))
	if numOutLinks > 0 {
		pagerankShare := node.PageRank / numOutLinks
		//Writes on channel nodeID, associated pagerankShare and outNode-IDs
		result <- fmt.Sprintf("%s\t%f\t%s", node.ID, pagerankShare, strings.Join(node.OutLinks, ","))
		fmt.Println(fmt.Sprintf("%s\t%f\t%s", node.ID, pagerankShare, strings.Join(node.OutLinks, ",")))
	}
	fmt.Println("-----------------------------")
	return
}
