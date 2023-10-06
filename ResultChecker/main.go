package main

import (
	"ResultChecker/constants"
	"ResultChecker/internal"
	"ResultChecker/mapper"
	"ResultChecker/reducer"
	"errors"
	"fmt"
	"log"
	"os"
	"sync"
	"time"
)

func main() {
	time.Sleep(2 * time.Second)
	start := time.Now()
	var config constants.Config
	constants.ReadJsonConfig(&config)
	logMessage := fmt.Sprintf("\n%s \n- Starting pagerank algorithm -", time.Now().Format("2006-01-02 15:04:05"))
	internal.WriteOnLog(logMessage)

	internal.CreateRandomGraph(config.NumNodes, config.EdgesToAttach, config.Seed, config.PrintGraph)

	graph := internal.Convert(config.GraphPath)
	if graph == nil {
		log.Fatalf("Something went wrong during file opening, aborting")
		return
	}

	var oldPageRankList []float64
	var newPageRankList []float64

	// output results
	mapperOutputArrayList := make([]*mapper.MapperOutput, config.NumNodes)
	reducerOutputArrayList := make([]*reducer.ReducerOutput, config.NumNodes)
	mapperCleanUpOutputArrayList := make([]*mapper.CleanUpOutput, config.NumNodes)
	reducerCleanUpOutputArrayList := make([]*reducer.ReducerOutput, config.NumNodes)

	numNodes := len(graph)
	convergence := false
	iteration := 0

	logMessage = fmt.Sprintf("\n%s \n- Starting pagerank algorithm -", time.Now().Format("2006-01-02 15:04:05"))
	wol(logMessage)
	var wg sync.WaitGroup
	for !convergence && iteration < config.MaxIteration {
		func() {
			iteration++
			aggregatePageRankShares := make(map[int][]float32)
			sinkMass := 0.0
			log.Printf("\nIteration number: %d", iteration)
			oldPageRankList = internal.ListOfPageRank(graph)

			//----- MAPPER -> MAP -----
			logMessage = fmt.Sprintf("\n\nITERATION -> %d", iteration)
			wol(logMessage)
			for _, node := range graph {
				wg.Add(1)
				node := node
				go func() {
					defer wg.Done()
					mapperInput := mapper.MapperInput{
						PageRank:      float32(node.PageRank),
						AdjacencyList: internal.GetOutLinks(node),
					}
					mapperOutput, err := Map(&mapperInput)
					mapperOutputArrayList[node.ID] = mapperOutput
					if err != nil {
						log.Fatalf("Error when calling Map function: %s", err)
					}
				}()
			}
			wg.Wait()
			//Shuffle parts of map-reduce paradigm
			// For each node in the graph save its contributes to other nodes
			for _, graphNode := range graph {
				for _, node := range mapperOutputArrayList[graphNode.ID].GetAdjacencyList() {
					//Update page rank aggregate table by appending for each node in the mapper output corresponding share
					aggregatePageRankShares[int(node)] = append(aggregatePageRankShares[int(node)], mapperOutputArrayList[graphNode.ID].PageRankShare)
				}
			}

			//----- REDUCER -> REDUCE -----
			for _, node := range graph {
				wg.Add(1)

				node := node
				go func() {
					defer wg.Done()
					reducerInput := reducer.ReducerInput{
						NodeId:         int32(node.ID),
						PageRankShares: aggregatePageRankShares[node.ID],
						GraphSize:      int32(numNodes),
					}
					reducerOutput, err := Reduce(&reducerInput)
					reducerOutputArrayList[node.ID] = reducerOutput
					if err != nil {
						log.Fatalf("Error when calling Reduce function: %s", err)
					}
				}()
			}
			wg.Wait()
			//Update node page rank value
			for _, node := range graph {
				//Update node page rank value
				node.PageRank = float64(reducerOutputArrayList[node.ID].NewRankValue)
			}

			//----- CLEAN UP PHASE -----
			//----- MAPPER-> CLEAN UP -----
			for _, node := range graph {
				node := node
				wg.Add(1)
				go func() {
					defer wg.Done()

					mapperInput := mapper.CleanUpInput{
						PageRank:      float32(node.PageRank),
						AdjacencyList: internal.GetOutLinks(node),
					}
					//Sums sink's mass
					cleanUpOutput, err := CleanUp(&mapperInput)
					mapperCleanUpOutputArrayList[node.ID] = cleanUpOutput
					if err != nil {
						log.Fatalf("Error when calling Map function: %s", err)
					}
				}()
			}
			wg.Wait()
			for _, node := range graph {
				sinkMass += float64(mapperCleanUpOutputArrayList[node.ID].SinkMass)
			}

			//----- REDUCER -> REDUCE-CLEANUP -----
			for _, node := range graph {
				node := node
				wg.Add(1)
				go func() {
					defer wg.Done()
					reducerCleanUpInput := reducer.ReducerCleanUpInput{
						NodeId:          int32(node.ID),
						CurrentPageRank: float32(node.PageRank),
						GraphSize:       int32(numNodes),
						SinkMass:        float32(sinkMass),
					}
					reducerOutput, err := ReduceCleanUp(&reducerCleanUpInput)
					reducerCleanUpOutputArrayList[node.ID] = reducerOutput

					if err != nil {
						log.Fatalf("Error when calling Reduce function: %s", err)
					}
				}()

			}
			wg.Wait()
			for _, node := range graph {
				//Update node page rank value
				node.PageRank = float64(reducerCleanUpOutputArrayList[node.ID].NewRankValue)
			}

			//Get new page rank to check the differences between the old ones
			newPageRankList = internal.ListOfPageRank(graph)

			i := 0
			//Save on the log intermediate update
			for _, nodePageRank := range newPageRankList {
				wol(fmt.Sprintf("\nNode %d -> PageRank %f", i, nodePageRank))
				i++
			}

			// Check the convergence of the algorithm
			convergence = internal.Convergence(oldPageRankList, newPageRankList)

		}()

	}

	// Print final results
	if convergence {
		log.Printf("\n\nObtained convergence after %d iteration, here final results: ", iteration)
		logMessage = fmt.Sprintf("\n\nObtained convergence after %d iteration, here final results: ", iteration)
		wol(logMessage)

	} else {
		log.Print("\n\nConvergence isn't obtained try to do more iterations")
		wol("\n\nConvergence isn't obtained try to do more iterations")

	}
	pageRankSum := 0.0
	for _, node := range graph {
		log.Printf("\nNodo: %d, PageRank: %f", node.ID, node.PageRank)
		logMessage = fmt.Sprintf("\nNodo: %d, PageRank: %f", node.ID, node.PageRank)
		wol(logMessage)
		pageRankSum += node.PageRank
	}

	log.Print("\n--Consistency check--\nSum of pageRank values: ", pageRankSum)
	wol(fmt.Sprintf("\n\n--Consistency check--\nSum of pageRank values: %f", pageRankSum))
	// models.PlotGraphByPageRank(graph)
	wol("\n\nPage rank algorithm run is done, bye bye\n")
	wol("----------------------------------------------------")

	elapsed := time.Since(start)
	log.Printf("\nPageRank algorithm tooks: %s", elapsed)
}

func Map(input *mapper.MapperInput) (*mapper.MapperOutput, error) {
	numOutLinks := float32(len(input.AdjacencyList))
	if numOutLinks > 0 {
		pagerankShare := input.PageRank / numOutLinks

		log.Printf("Share: %f for nodes: %v\n", pagerankShare, input.GetAdjacencyList())

		return &mapper.MapperOutput{
			PageRankShare: pagerankShare,
			AdjacencyList: input.GetAdjacencyList(),
		}, nil

	} else {
		//If here, node has zero out-links
		zero := 0.0
		return &mapper.MapperOutput{
			PageRankShare: float32(zero),
			AdjacencyList: input.GetAdjacencyList(),
		}, nil
	}
}

// CleanUp -> to manage sinks and random jump factor
func CleanUp(input *mapper.CleanUpInput) (*mapper.CleanUpOutput, error) {
	numOutLinks := float32(len(input.AdjacencyList))
	if numOutLinks == 0 {
		log.Printf("Sink finded")
		return &mapper.CleanUpOutput{
			SinkMass: input.PageRank,
		}, nil
	} else {
		//If here, node has out-links, not interesting in clean-up phase
		zero := 0.0
		return &mapper.CleanUpOutput{
			SinkMass: float32(zero),
		}, nil
	}
}

// Reduce ->  sum aggregated page rank values and return the value
func Reduce(input *reducer.ReducerInput) (*reducer.ReducerOutput, error) {
	var config constants.Config
	constants.ReadJsonConfig(&config)
	if input != nil {
		accumulator := float32(0.0)
		for _, rank := range input.PageRankShares {
			accumulator += rank
		}

		newPageRank := (1.0-config.DampingFactor)/float32(input.GraphSize) + config.DampingFactor*accumulator
		log.Printf("\nNodeID: %d -> evaluated page rank: %f", input.NodeId, newPageRank)

		return &reducer.ReducerOutput{
			NodeId:       input.NodeId,
			NewRankValue: newPageRank,
		}, nil
	} else {

		return nil, errors.New("input isn't valid")
	}

}

// ReduceCleanUp -> use the cleanUp formula to fix page rank value
func ReduceCleanUp(input *reducer.ReducerCleanUpInput) (*reducer.ReducerOutput, error) {
	var config constants.Config
	constants.ReadJsonConfig(&config)
	if input != nil {
		massShares := input.SinkMass / float32(input.GraphSize)
		newPageRank := (1.0-config.DampingFactor)/float32(input.GraphSize) + config.DampingFactor*(input.CurrentPageRank+massShares)
		return &reducer.ReducerOutput{
			NodeId:       input.NodeId,
			NewRankValue: newPageRank,
		}, nil
	} else {

		return nil, errors.New("input isn't valid")
	}
}

func wol(logMessage string) {
	// Open the file in append mod
	file, err := os.OpenFile("log.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("Impossible to open log file: %v", err)
		return
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Fatalf("Impossibile to close log file: %v", err)
		}
	}(file)

	// Write into the log
	_, err = file.WriteString(logMessage)
	if err != nil {
		log.Fatalf("Impossible to write on log: %v", err)
		return
	}

}
