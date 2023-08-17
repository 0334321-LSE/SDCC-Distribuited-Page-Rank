package main

import (
	"Master/constants"
	"Master/internal"
	"Master/mapper"
	"Master/reducer"
	"container/ring"
	"context"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"time"
)

func main() {
	time.Sleep(2 * time.Second)
	start := time.Now()
	var config constants.Config
	constants.ReadJsonConfig(&config)
	logMessage := fmt.Sprintf("\n%s \n- Starting pagerank algorithm -", time.Now().Format("2006-01-02 15:04:05"))
	internal.WriteOnLog(logMessage)

	internal.CreateRandomGraph(config.NumNodes, config.EdgesToAttach, config.Seed)

	graph := internal.Convert(config.GraphPath)
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
	var mapperHbRing = ring.New(config.NumMapper)
	var reducerRing = ring.New(config.NumReducer)
	var reducerHbRing = ring.New(config.NumReducer)

	//Initializing mapper ring
	for i := 1; i <= config.NumMapper; i++ {
		mapperRing.Value = fmt.Sprintf("app-mapper-%d:%d", i, config.MapperPn+i)
		mapperRing = mapperRing.Next()
		mapperHbRing.Value = fmt.Sprintf("app-mapper-%d:%d", i, config.MapperHbPn+i)
		mapperHbRing = mapperHbRing.Next()
	}

	//Initializing reducer ring
	for i := 1; i <= config.NumReducer; i++ {
		reducerRing.Value = fmt.Sprintf("app-reducer-%d:%d", i, config.ReducerPn+i)
		reducerRing = reducerRing.Next()
		reducerHbRing.Value = fmt.Sprintf("app-reducer-%d:%d", i, config.ReducerHbPn+i)
		reducerHbRing = reducerHbRing.Next()
	}

	connWithMapper := make(map[int][]*grpc.ClientConn)
	connWithMapperHb := make(map[int][]*grpc.ClientConn)
	connWithReducer := make(map[int][]*grpc.ClientConn)
	connWithReducerHb := make(map[int][]*grpc.ClientConn)

	//----- CONNECTIONS INITIALIZING -----

	// Initialize connection with each container for Mapper and Heartbeat service
	// Mapper service
	connWithMapper = internal.SetGrpcConnection(mapperRing)
	// Heartbeat service
	connWithMapperHb = internal.SetGrpcConnection(mapperHbRing)

	// Initialize connection with each container for Reducer and Heartbeat service
	// Reducer service
	connWithReducer = internal.SetGrpcConnection(reducerRing)
	// Heartbeat service
	connWithReducerHb = internal.SetGrpcConnection(reducerHbRing)

	// Fix potential holes in map
	internal.FixMapsKeys(&connWithMapper, &connWithMapperHb, &connWithReducer, &connWithReducerHb)

	for !convergence || iteration == config.MaxIteration {
		func() {
			iteration++
			aggregatePageRankShares := make(map[int][]float32)
			sinkMass := 0.0
			log.Printf("\nIteration number: %d", iteration)
			oldPageRankList = internal.ListOfPageRank(graph)

			// Launch a job for each node of the graph, with a round robing scheduling among containers
			//----- MAPPER -> MAP -----
			logMessage = fmt.Sprintf("\n\nITERATION -> %d", iteration)
			internal.WriteOnLog(logMessage)

			for m, node := range graph {
				if mapperRing.Len() == 0 {
					log.Fatalf("No more container are avaible, try to re-run program")
				}
				chosen := internal.CheckIfMapperIsAlive(m, &connWithMapper, &mapperRing, &connWithMapperHb, &mapperHbRing)
				// If at least one container is alive, launch map task
				// Connection with MapperClient on ports 900X, for each node launch MAP job
				mapperConnection := mapper.NewMapperClient(connWithMapper[chosen][0])
				// Now is connected with I-TH container, launch map task
				mapperInput := mapper.MapperInput{
					PageRank:      float32(node.PageRank),
					AdjacencyList: internal.GetOutLinks(node),
				}
				mapperOutput, err := mapperConnection.Map(context.Background(), &mapperInput)
				if err != nil {
					log.Fatalf("Error when calling Map function: %v", err)
				}
				//Shuffle parts of map-reduce paradigm
				for _, node := range mapperOutput.GetAdjacencyList() {
					//Update page rank aggregate table by appending for each node in the mapper output corresponding share
					aggregatePageRankShares[int(node)] = append(aggregatePageRankShares[int(node)], mapperOutput.PageRankShare)
				}
			}

			//----- REDUCER -> REDUCE -----
			for m, node := range graph {
				if reducerRing.Len() == 0 {
					log.Fatalf("No more container are avaible, try to re-run program")
				}
				chosen := internal.CheckIfReducerIsAlive(m, &connWithReducer, &reducerRing, &connWithReducerHb, &reducerHbRing)
				// If at least one container is alive, launch reduce task
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
				if mapperRing.Len() == 0 {
					log.Fatalf("No more container are avaible, try to re-run program")
				}
				chosen := internal.CheckIfMapperIsAlive(m, &connWithMapper, &mapperRing, &connWithMapperHb, &mapperHbRing)
				// If at least one container is alive, launch map clean up task
				// Connection with MapperClient on ports 900X, for each node launch MAP-CLEAN job
				mapperConnection := mapper.NewMapperClient(connWithMapper[chosen][0])
				mapperInput := mapper.CleanUpInput{
					PageRank:      float32(node.PageRank),
					AdjacencyList: internal.GetOutLinks(node),
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
				if reducerRing.Len() == 0 {
					log.Fatalf("No more container are avaible, try to re-run program")
				}
				chosen := internal.CheckIfReducerIsAlive(m, &connWithReducer, &reducerRing, &connWithReducerHb, &reducerHbRing)
				// If at least one container is alive, launch reducer clean up task
				// Connection with ReducerClient on ports 10000X, for each node launch REDUCE-CLEAN job
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
			newPageRankList = internal.ListOfPageRank(graph)

			i := 0
			//Save on the log intermediate update
			for _, nodePageRank := range newPageRankList {
				internal.WriteOnLog(fmt.Sprintf("\nNode %d -> PageRank %f", i, nodePageRank))
				i++
			}

			// Check the convergence of the algorithm
			convergence = internal.Convergence(oldPageRankList, newPageRankList)

		}()
	}

	// Close client connections
	internal.CloseClientConn(connWithMapper)
	internal.CloseClientConn(connWithMapperHb)
	internal.CloseClientConn(connWithReducer)
	internal.CloseClientConn(connWithReducerHb)

	// Print final results
	if convergence {
		log.Printf("\n\nObtained convergence after %d iteration, here final results: ", iteration)
		logMessage = fmt.Sprintf("\n\nObtained convergence after %d iteration, here final results: ", iteration)
		internal.WriteOnLog(logMessage)

	} else {
		log.Print("\n\nConvergence isn't obtained try to do more iterations")
		internal.WriteOnLog("\n\nConvergence isn't obtained try to do more iterations")

	}
	pageRankSum := 0.0
	for _, node := range graph {
		log.Printf("\nNodo: %d, PageRank: %f", node.ID, node.PageRank)
		logMessage = fmt.Sprintf("\nNodo: %d, PageRank: %f", node.ID, node.PageRank)
		internal.WriteOnLog(logMessage)
		pageRankSum += node.PageRank
	}

	log.Print("\n--Consistency check--\nSum of pageRank values: ", pageRankSum)
	internal.WriteOnLog(fmt.Sprintf("\n\n--Consistency check--\nSum of pageRank values: %f", pageRankSum))
	internal.PlotGraphByPageRank(graph)
	internal.WriteOnLog("\n\nPage rank algorithm run is done, bye bye\n")
	internal.WriteOnLog("----------------------------------------------------")

	elapsed := time.Since(start)
	log.Printf("\nPageRank algorithm tooks: %s", elapsed)
}
