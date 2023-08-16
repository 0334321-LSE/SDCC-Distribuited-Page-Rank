package internal

import (
	"Master/mapper"
	"Master/reducer"
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
			log.Fatalf("Could not connect: %s", err)
		}
		defer func(conn *grpc.ClientConn) {
			err := conn.Close()
			if err != nil {
				log.Fatalf("Something went wrong during connection closing %v", err)
			}
		}(newMap[i][0])

	}
	return newMap
}

// FixMapKeys -> Take a map with some missing connection, return a map without "holes" in key set
func FixMapKeys(originalMap map[int][]*grpc.ClientConn) map[int][]*grpc.ClientConn {
	if len(originalMap) == 0 {
		log.Fatal("\nMap is empty, no more connection are available")

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
	for key := range keys {
		newMap[i] = append(newMap[i], originalMap[key][0])
		i++
	}
	return newMap
}

func CheckIfMapperIsAlive(m int, connWithMapper map[int][]*grpc.ClientConn, mapperRing *ring.Ring, connWithMapperHb map[int][]*grpc.ClientConn, mapperHbRing *ring.Ring) int {
	var chosen int
	for alive := false; alive == false; {
		//M % N.Container to establish which one must be chosen (round-robin)
		chosen = m % mapperRing.Len()

		// Connection with heartbeat service of chosen Mapper on port 5000X
		// Set 5 second timeout
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		mapperHeartbeatConnection := mapper.NewMapperHeartbeatClient(connWithMapperHb[chosen][0])
		request := mapper.MapperHeartbeatRequest{
			Alive: false,
		}
		response, err := mapperHeartbeatConnection.Ping(ctx, &request)
		if err != nil {
			// If occurs a timeout
			if status.Code(err) == codes.DeadlineExceeded {
				log.Printf("\nTimeout expired, removed connection with container")
				// Remove container from rings
				RemoveFromRing(mapperRing, chosen)
				RemoveFromRing(mapperHbRing, chosen)

				// Discard down-connections
				delete(connWithMapper, chosen)
				delete(connWithMapperHb, chosen)

				// Fix keys of connection map
				connWithMapper = FixMapKeys(connWithMapper)
				connWithMapperHb = FixMapKeys(connWithMapperHb)

			} else {
				log.Printf("Error when calling Map function: %v", err)
			}
		}
		// If there isn't an error, container is alive
		alive = response.Alive
	}
	return chosen
}

func CheckIfReducerIsAlive(m int, connWithReducer map[int][]*grpc.ClientConn, reducerRing *ring.Ring, connWithReducerHb map[int][]*grpc.ClientConn, reducerHbRing *ring.Ring) int {
	var chosen int
	for alive := false; alive == false; {
		//M % N.Container to establish which one must be chosen (round-robin)
		chosen = m % reducerRing.Len()

		// Connection with heartbeat service of chosen Mapper on port 5000X
		// Set 5 second timeout
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		reducerHeartbeatConnection := reducer.NewReducerHeartbeatClient(connWithReducerHb[chosen][0])
		request := reducer.ReducerHeartbeatRequest{
			Alive: false,
		}
		response, err := reducerHeartbeatConnection.Ping(ctx, &request)
		if err != nil {
			// If occurs a timeout
			if status.Code(err) == codes.DeadlineExceeded {
				log.Printf("\nTimeout expired, removed connection with container")
				// Remove container from rings
				RemoveFromRing(reducerRing, chosen)
				RemoveFromRing(reducerHbRing, chosen)

				// Discard down-connections
				delete(connWithReducer, chosen)
				delete(connWithReducerHb, chosen)

				// Fix keys of connection map
				connWithReducer = FixMapKeys(connWithReducer)
				connWithReducerHb = FixMapKeys(connWithReducerHb)

			} else {
				log.Printf("Error when calling Map function: %v", err)
			}
		}
		// If there isn't an error, container is alive
		alive = response.Alive
	}
	return chosen
}
