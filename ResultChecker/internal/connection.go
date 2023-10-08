package internal

import (
	"ResultChecker/mapper"
	"ResultChecker/reducer"
	"container/ring"
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"log"
	"sort"
	"time"
)

// SetGrpcConnection -> return map that contains connections with grpcClient
func SetGrpcConnection(r *ring.Ring) map[int][]*grpc.ClientConn {
	var err error
	var connection *grpc.ClientConn
	newMap := make(map[int][]*grpc.ClientConn)
	// Initialize connection with each container for Mapper and Heartbeat service
	for i := 0; i < r.Len(); i++ {
		// Mapper service
		connection, err = grpc.Dial(r.Value.(string), grpc.WithTransportCredentials(insecure.NewCredentials()))
		newMap[i] = append(newMap[i], connection)
		//Next container
		r = r.Next()
		if err != nil {
			log.Printf("Could not connect with %s\n Error: %v", r.Value.(string), err)
		}
	}
	return newMap
}

// FixMapKeys -> Take connection map with some missing connection, return a map without "holes" in key set
func FixMapKeys(originalMap map[int][]*grpc.ClientConn, mapType string) map[int][]*grpc.ClientConn {
	if len(originalMap) == 0 {
		log.Fatalf("\n%s connection map is empty, no more connection are available.\nTry to re-launch program.", mapType)
	}
	newMap := make(map[int][]*grpc.ClientConn)
	// Obtains key from originalMap
	var keys []int
	for key := range originalMap {
		keys = append(keys, key)
	}
	sort.Ints(keys)
	// Now the key are sorted, fix keys of connection map
	i := 0
	for _, key := range keys {
		newMap[i] = append(newMap[i], originalMap[key][0])
		i++
	}
	return newMap
}

// FixMapsKeys -> Take all the connections maps and fixes possible holes in key set
func FixMapsKeys(connWithMapper *map[int][]*grpc.ClientConn, connWithMapperHb *map[int][]*grpc.ClientConn, connWithReducer *map[int][]*grpc.ClientConn, connWithReducerHb *map[int][]*grpc.ClientConn) {
	*connWithMapper = FixMapKeys(*connWithMapper, "Mapper")
	*connWithMapperHb = FixMapKeys(*connWithMapperHb, "MapperHeartbeat")
	*connWithReducer = FixMapKeys(*connWithReducer, "Reducer")
	*connWithReducerHb = FixMapKeys(*connWithReducerHb, "ReducerHeartbeat")
}

// CheckIfMapperIsAlive -> as the name says, check by doing a ping if there is a mapper alive, return the id of the worker that must be called
func CheckIfMapperIsAlive(m int, connWithMapper *map[int][]*grpc.ClientConn, mapperRing **ring.Ring, connWithMapperHb *map[int][]*grpc.ClientConn, mapperHbRing **ring.Ring) int {

	var chosen int
	index := 0
	for alive := false; alive == false && index != (*mapperRing).Len(); {
		//M % N.Container to establish which one must be chosen (round-robin)
		chosen = (m + index) % (*mapperRing).Len()

		// Connection with heartbeat service of chosen Mapper on port 5000X
		// Set 5 second timeout
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		mapperHeartbeatConnection := mapper.NewMapperHeartbeatClient((*connWithMapperHb)[chosen][0])
		request := mapper.MapperHeartbeatRequest{
			Alive: false,
		}
		response, err := mapperHeartbeatConnection.Ping(ctx, &request)
		if err != nil {
			// If occurs a timeout
			if status.Code(err) == codes.DeadlineExceeded {
				log.Printf("\nError when calling mapper-%d, try with another one", chosen+1)
			} else {
				// Otherwise if occurs some-other problem
				log.Printf("\nError when calling mapper-%d, try with another one", chosen+1)
			}
		} else {
			// If there isn't an error, container is alive
			alive = response.GetAlive()
			return chosen
		}
		index++
	}
	if index == (*mapperRing).Len() {
		log.Fatalf("No more Mapper are available, try to re-start program")
	}
	return chosen
}

// CheckIfReducerIsAlive -> as the name says, check by doing a ping if there is a reducer alive, return the id of the worker that must be called
func CheckIfReducerIsAlive(m int, connWithReducer *map[int][]*grpc.ClientConn, reducerRing **ring.Ring, connWithReducerHb *map[int][]*grpc.ClientConn, reducerHbRing **ring.Ring) int {
	var chosen int
	index := 0
	for alive := false; alive == false && (*reducerRing).Len() != index; {
		//M % N.Container to establish which one must be chosen (round-robin)
		chosen = (m + index) % (*reducerRing).Len()

		// Connection with heartbeat service of chosen Mapper on port 5000X
		// Set 5 second timeout
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		reducerHeartbeatConnection := reducer.NewReducerHeartbeatClient((*connWithReducerHb)[chosen][0])
		request := reducer.ReducerHeartbeatRequest{
			Alive: false,
		}
		response, err := reducerHeartbeatConnection.Ping(ctx, &request)
		if err != nil {
			// If occurs a timeout
			if status.Code(err) == codes.DeadlineExceeded {
				log.Printf("\nError when calling reducer-%d, try with another one", chosen+1)
			} else {
				// Otherwise if occurs some-other problem
				log.Printf("\nError when calling reducer-%d, try with another one", chosen+1)
			}
		} else {
			// If there isn't an error, container is alive
			alive = response.Alive
			return chosen
		}
		index++
	}
	if (*reducerRing).Len() == index {
		log.Fatalf("No more Reducer are available, try to re-start program")
	}
	return chosen
}

// CloseClientConn -> close all the opened client connections
func CloseClientConn(m map[int][]*grpc.ClientConn) {
	for i := 0; i < len(m); i++ {
		func(conn *grpc.ClientConn) {
			err := conn.Close()
			if err != nil {
				log.Fatalf("Something went wrong during connection closing %v", err)
			}
		}(m[i][0])
	}
}
