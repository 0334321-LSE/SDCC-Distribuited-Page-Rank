package main

import (
	"PageRank/constants"
	"PageRank/models"
	"PageRank/utils"
	"fmt"
	"log"
	"sync"
	"time"
)

func main() {
	start := time.Now()
	utils.CreateRandomGraph(50, 2)

	graph := utils.Convert(constants.GraphPath)
	if graph == nil {
		fmt.Println("Something went wrong during file opening, aborting")
		return
	}
	var oldPageRank []float64
	var newPageRankList []float64

	numNodes := len(graph)
	convergence := false
	iteration := 0

	for !convergence || iteration == constants.MaxIteration {

		iteration++
		fmt.Println("\nIteration number: ", iteration)
		oldPageRank = models.ListOfPageRank(graph)

		// Chanel for intermediate result of mappers
		mapOutput := make(chan string, numNodes)
		// Chanel for intermediate result of reducers
		reduceOutput := make(chan string, numNodes)

		/* Without go routine
		// Starting mapper phase, for each node launch a mapper
		for _, node := range graph {
			utils.MapPageRank(node, mapOutput)
		}

		// Closing Map channel
		close(mapOutput)

		// Launch reduce phase
		utils.ReducePageRank(mapOutput, reduceOutput, numNodes)

		// Closing Reduce channel
		close(reduceOutput)
		*/

		/* with go routine */

		wg := new(sync.WaitGroup)
		// Starting mapper nodes, for each node launch a mapper
		for _, node := range graph {
			wg.Add(1)
			go func(node *models.Node) {
				defer wg.Done()
				utils.MapPageRank(node, mapOutput)

			}(node)
		}

		// Closing Map chanel
		wg.Wait()
		close(mapOutput)

		// Starting reducer nodes for each node launch a mapper
		wg.Add(1)
		go func() {
			defer wg.Done()
			utils.ReducePageRank(mapOutput, reduceOutput, numNodes)
		}()

		// Closing Reduce chanel
		wg.Wait()
		close(reduceOutput)

		// Update page rank value
		models.UpdatePageRanks(graph, reduceOutput)

		newPageRankList = models.ListOfPageRank(graph)

		// Check the convergence of the algorithm
		convergence = models.Convergence(oldPageRank, newPageRankList)
	}

	// Print final results
	if convergence {
		fmt.Println("Obtained convergence, here final results:")
	} else {
		fmt.Println("Convergence isn't obtained try to do more iterations")
	}
	for _, node := range graph {
		fmt.Printf("Nodo: %s, PageRank: %f\n", node.ID, node.PageRank)
	}

	models.PlotGraphByPageRank(graph)
	elapsed := time.Since(start)
	log.Printf("PageRank algorithm tooks: %s", elapsed)
}
