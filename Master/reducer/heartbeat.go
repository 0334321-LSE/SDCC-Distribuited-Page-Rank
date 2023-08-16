package reducer

import (
	"context"
	"errors"
	"log"
)

type ReducerHeartbeat struct {
}

func (heartbeat *ReducerHeartbeat) mustEmbedUnimplementedReducerHeartbeatServer() {}

func (heartbeat *ReducerHeartbeat) Ping(ctx context.Context, input *ReducerHeartbeatRequest) (*ReducerHeartbeatResponse, error) {
	if input != nil {
		log.Printf("Pong")
		return &ReducerHeartbeatResponse{
			Alive: true,
		}, nil
	} else {
		return nil, errors.New("input isn't valid")
	}
}
