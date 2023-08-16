package main

import (
	"Master/constants"
	"Master/mapper"
	"Master/models"
	"Master/reducer"
	"Master/utils"
	"container/ring"
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"time"
)

func main() {
	time.Sleep(2 * time.Second)
	start := time.Now()
	var config constants.Config
	constants.ReadJsonConfig(&config)
	logMessage := fmt.Sprintf("\n%s \n- Starting pagerank algorithm -", time.Now().Format("2006-01-02 15:04:05"))
	utils.WriteOnLog(logMessage)

	utils.CreateRandomGraph(config.NumNodes, config.EdgesToAttach, config.Seed)

	graph := utils.Convert(config.GraphPath)
	if graph == nil {
		fmt.Println("Something went wrong during file opening, aborting")
		return
	}
	var oldPageRankList []float64
	var newPageRankList []float64

	numNodes := len(graph)
	convergence := false
	iteration := 0

	//Obtain configuration parameter
	var mapperRing = ring.New(config.NumMapper)
	var reducerRing = ring.New(config.NumReducer)

	for i := 0; i < config.NumMapper; i++ {
		mapperRing.Value = fmt.Sprintf("app-mapper-%d:%d", i+1, config.MapperPN+i)
		mapperRing = mapperRing.Next()
	}

	for i := 0; i < config.NumReducer; i++ {
		reducerRing.Value = fmt.Sprintf("app-reducer-%d:%d", i+1, config.ReducerPN+i)
		reducerRing = reducerRing.Next()
	}

	for !convergence || iteration == config.MaxIteration {
		func() {
			iteration++
			aggregatePageRankShares := make(map[int][]float32)
			sinkMass := 0.0
			log.Printf("\nIteration number: %d", iteration)
			oldPageRankList = models.ListOfPageRank(graph)

			//----- MAPPER -> MAP -----
			// With round-robin policies call each container

			connWithMapper := make(map[int][]*grpc.ClientConn)
			connWithReducer := make(map[int][]*grpc.ClientConn)
			var err error
			var connection *grpc.ClientConn

			// Initialize connection with each container for Mapper task
			for i := 0; i < mapperRing.Len(); i++ {
				connection, err = grpc.Dial(mapperRing.Value.(string), grpc.WithTransportCredentials(insecure.NewCredentials()))
				connWithMapper[i] = append(connWithMapper[i], connection)
				//Next container
				mapperRing = mapperRing.Next()
				if err != nil {
					log.Fatalf("Could not connect: %s", err)
				}
				defer func(conn *grpc.ClientConn) {
					err := conn.Close()
					if err != nil {
						log.Fatalf("Something went wrong during connection closing %v", err)
					}
				}(connWithMapper[i][0])
			}

			logMessage = fmt.Sprintf("\n\nITERATION -> %d", iteration)
			utils.WriteOnLog(logMessage)

			for m, node := range graph {
				//M % N.Container to establish which one must be chosen (round-robin)
				chosen := m % mapperRing.Len()
				// Connection with MapperClient on ports 900X, for each node launch MAP job
				mapperConnection := mapper.NewMapperClient(connWithMapper[chosen][0])
				// Now is connected with I-TH container, launch map task
				mapperInput := mapper.MapperInput{
					PageRank:      float32(node.PageRank),
					AdjacencyList: models.GetOutLinks(node),
				}
				mapperOutput, err := mapperConnection.Map(context.Background(), &mapperInput)
				//Shuffle parts of map-reduce paradigm
				for _, node := range mapperOutput.GetAdjacencyList() {
					//Update page rank aggregate table by appending for each node in the mapper output corresponding share
					aggregatePageRankShares[int(node)] = append(aggregatePageRankShares[int(node)], mapperOutput.PageRankShare)
				}
				if err != nil {
					log.Fatalf("Error when calling Map function: %v", err)
				}
			}

			//----- REDUCER -> REDUCE -----

			// Initialize connection with each container for Reducer task
			for i := 0; i < reducerRing.Len(); i++ {
				connection, err = grpc.Dial(reducerRing.Value.(string), grpc.WithTransportCredentials(insecure.NewCredentials()))
				connWithReducer[i] = append(connWithReducer[i], connection)
				//Next container
				reducerRing = reducerRing.Next()
				if err != nil {
					log.Fatalf("Could not connect: %s", err)
				}
				defer func(conn *grpc.ClientConn) {
					err := conn.Close()
					if err != nil {
						log.Fatalf("Something went wrong during connection closing %v", err)
					}
				}(connWithReducer[i][0])
			}

			for m, node := range graph {
				//M % N.Container to establish which one must be chosen (round-robin)
				chosen := m % reducerRing.Len()
				// Connection with ReducerClient on ports 10000X, for each node launch REDUCE-job
				reducerConnection := reducer.NewReducerClient(connWithReducer[chosen][0])
				// Now is connected with I-TH container, launch map task
				reducerInput := reducer.ReducerInput{
					NodeId:         int32(node.ID),
					PageRankShares: aggregatePageRankShares[node.ID],
					GraphSize:      int32(numNodes),
				}
				reducerOutput, err := reducerConnection.Reduce(context.Background(), &reducerInput)
				if err != nil {
					log.Fatalf("Error when calling Reduce function: %s", err)
				}
				//Update node page rank value
				node.PageRank = float64(reducerOutput.NewRankValue)
			}

			//----- CLEAN UP PHASE -----

			//----- MAPPER-> CLEAN UP -----
			for m, node := range graph {
				//M % N.Container to establish which one must be chosen (round-robin)
				chosen := m % mapperRing.Len()
				mapperConnection := mapper.NewMapperClient(connWithMapper[chosen][0])
				mapperInput := mapper.CleanUpInput{
					PageRank:      float32(node.PageRank),
					AdjacencyList: models.GetOutLinks(node),
				}
				//Sums sink's mass
				cleanUpOutput, err := mapperConnection.CleanUp(context.Background(), &mapperInput)
				sinkMass += float64(cleanUpOutput.SinkMass)
				if err != nil {
					log.Fatalf("Error when calling Map function: %s", err)
				}
			}

			//----- REDUCER -> REDUCE-CLEANUP -----
			for m, node := range graph {
				//M % N.Container to establish which one must be chosen (round-robin)
				chosen := m % reducerRing.Len()
				reducerConnection := reducer.NewReducerClient(connWithReducer[chosen][0])
				reducerCleanUpInput := reducer.ReducerCleanUpInput{
					NodeId:          int32(node.ID),
					CurrentPageRank: float32(node.PageRank),
					GraphSize:       int32(numNodes),
					SinkMass:        float32(sinkMass),
				}
				reducerOutput, err := reducerConnection.ReduceCleanUp(context.Background(), &reducerCleanUpInput)
				if err != nil {
					log.Fatalf("Error when calling Reduce function: %s", err)
				}
				//Update node page rank value
				node.PageRank = float64(reducerOutput.NewRankValue)
			}

			//Get new page rank to check the differences between the old ones
			newPageRankList = models.ListOfPageRank(graph)

			i := 0
			//Save on the log intermediate update
			for _, nodePageRank := range newPageRankList {
				utils.WriteOnLog(fmt.Sprintf("\nNode %d -> PageRank %f", i, nodePageRank))
				i++
			}

			// Check the convergence of the algorithm
			convergence = models.Convergence(oldPageRankList, newPageRankList)

		}()

	}

	// Print final results
	if convergence {
		log.Printf("\n\nObtained convergence after %d iteration, here final results: ", iteration)
		logMessage = fmt.Sprintf("\n\nObtained convergence after %d iteration, here final results: ", iteration)
		utils.WriteOnLog(logMessage)

	} else {
		log.Print("\n\nConvergence isn't obtained try to do more iterations")
		utils.WriteOnLog("\n\nConvergence isn't obtained try to do more iterations")

	}
	pageRankSum := 0.0
	for _, node := range graph {
		log.Printf("\nNodo: %d, PageRank: %f", node.ID, node.PageRank)
		logMessage = fmt.Sprintf("\nNodo: %d, PageRank: %f", node.ID, node.PageRank)
		utils.WriteOnLog(logMessage)
		pageRankSum += node.PageRank
	}

	log.Print("\n--Consistency check--\nSum of pageRank values: ", pageRankSum)
	utils.WriteOnLog(fmt.Sprintf("\n\n--Consistency check--\nSum of pageRank values: %f", pageRankSum))
	models.PlotGraphByPageRank(graph)
	utils.WriteOnLog("\n\nPage rank algorithm run is done, bye bye\n")
	utils.WriteOnLog("----------------------------------------------------")

	elapsed := time.Since(start)
	log.Printf("\nPageRank algorithm tooks: %s", elapsed)
}
