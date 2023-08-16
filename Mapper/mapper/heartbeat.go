package mapper

import (
	"context"
	"errors"
	"log"
)

type MapperHeartbeat struct {
}

func (heartbeat *MapperHeartbeat) mustEmbedUnimplementedMapperHeartbeatServer() {}

func (heartbeat *MapperHeartbeat) Ping(ctx context.Context, input *MapperHeartbeatRequest) (*MapperHeartbeatResponse, error) {
	if input != nil {
		log.Printf("Pong")
		return &MapperHeartbeatResponse{
			Alive: true,
		}, nil
	} else {
		return nil, errors.New("input isn't valid")
	}
}
