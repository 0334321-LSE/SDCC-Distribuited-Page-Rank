package main

import (
	"PageRank/constants"
	"PageRank/mapper"
	"PageRank/models"
	"PageRank/reducer"
	"PageRank/utils"
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"os"
	"time"
)

func main() {
	start := time.Now()
	utils.CreateRandomGraph(10, 5, 3)

	graph := utils.Convert(constants.GraphPath)
	if graph == nil {
		fmt.Println("Something went wrong during file opening, aborting")
		return
	}
	var oldPageRankList []float64
	var newPageRankList []float64

	numNodes := len(graph)
	convergence := false
	iteration := 0

	logMessage := fmt.Sprintf("\n%s \n- Starting pagerank algorithm -", time.Now().Format("2006-01-02 15:04:05"))
	writeOnLog(logMessage)

	for !convergence || iteration == constants.MaxIteration {
		func() {
			iteration++
			aggregatePageRankShares := make(map[int][]float32)
			sinkMass := 0.0
			log.Printf("\nIteration number: %d", iteration)
			oldPageRankList = models.ListOfPageRank(graph)

			//----- MAPPER -> MAP -----
			// Create a grpc client connection with port 9000 localhost
			var conn *grpc.ClientConn
			conn, err := grpc.Dial(":9000", grpc.WithTransportCredentials(insecure.NewCredentials()))
			if err != nil {
				log.Fatalf("Could not connect: %s", err)
			}
			defer func(conn *grpc.ClientConn) {
				err := conn.Close()
				if err != nil {
					log.Fatalf("Something went wrong during connection closing %v", err)
				}
			}(conn)

			// Connection with MapperClient at port 9000, for each node launch MAP job
			logMessage = fmt.Sprintf("\n\nITERATION -> %d", iteration)
			writeOnLog(logMessage)
			mapperConnection := mapper.NewMapperClient(conn)
			for _, node := range graph {
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
					log.Fatalf("Error when calling Map function: %s", err)
				}
			}

			//----- REDUCER -> REDUCE -----
			// Create a grpc client connection with port 9001 localhost
			conn2, err := grpc.Dial(":9001", grpc.WithTransportCredentials(insecure.NewCredentials()))
			if err != nil {
				log.Fatalf("Could not connect: %s", err)
			}
			defer func(conn2 *grpc.ClientConn) {
				err := conn2.Close()
				if err != nil {
					log.Fatalf("Something went wrong during connection closing %v", err)
				}
			}(conn2)

			// Connection with ReducerClient at port 9001, for each node launch REDUCE-job
			reducerConnection := reducer.NewReducerClient(conn2)

			for _, node := range graph {
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
			// Create a grpc client connection with port 9000 localhost
			conn, err = grpc.Dial(":9000", grpc.WithTransportCredentials(insecure.NewCredentials()))
			if err != nil {
				log.Fatalf("Could not connect: %s", err)
			}
			defer func(conn *grpc.ClientConn) {
				err := conn.Close()
				if err != nil {
					log.Fatalf("Something went wrong during connection closing %v", err)
				}
			}(conn)

			mapperConnection = mapper.NewMapperClient(conn)
			for _, node := range graph {
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
			// Create a grpc client connection with port 9001 localhost
			conn2, err = grpc.Dial(":9001", grpc.WithTransportCredentials(insecure.NewCredentials()))
			if err != nil {
				log.Fatalf("Could not connect: %s", err)
			}
			defer func(conn2 *grpc.ClientConn) {
				err := conn2.Close()
				if err != nil {
					log.Fatalf("Something went wrong during connection closing %v", err)
				}
			}(conn2)

			// Connection with ReducerClient at port 9001, for each node launch REDUCE-CLEANUP-job
			reducerConnection = reducer.NewReducerClient(conn2)
			for _, node := range graph {
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
				writeOnLog(fmt.Sprintf("\nNode %d -> PageRank %f", i, nodePageRank))
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
		writeOnLog(logMessage)

	} else {
		log.Print("\n\nConvergence isn't obtained try to do more iterations")
		writeOnLog("\n\nConvergence isn't obtained try to do more iterations")

	}
	pageRankSum := 0.0
	for _, node := range graph {
		log.Printf("\nNodo: %d, PageRank: %f", node.ID, node.PageRank)
		logMessage = fmt.Sprintf("\nNodo: %d, PageRank: %f", node.ID, node.PageRank)
		writeOnLog(logMessage)
		pageRankSum += node.PageRank
	}

	log.Print("\n--Consistency check--\nSum of pageRank values: ", pageRankSum)
	writeOnLog(fmt.Sprintf("\n\n--Consistency check--\nSum of pageRank values: %f", pageRankSum))
	models.PlotGraphByPageRank(graph)
	writeOnLog("\n\nPage rank algorithm run is done, bye bye\n")
	writeOnLog("----------------------------------------------------")

	elapsed := time.Since(start)
	log.Printf("\nPageRank algorithm tooks: %s", elapsed)
}

func writeOnLog(logMessage string) {
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
