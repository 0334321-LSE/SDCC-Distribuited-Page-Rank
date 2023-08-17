package reducer

import (
	"context"
	"errors"
)

type ReducerHeartbeat struct {
}

func (heartbeat *ReducerHeartbeat) mustEmbedUnimplementedReducerHeartbeatServer() {}

func (heartbeat *ReducerHeartbeat) Ping(ctx context.Context, input *ReducerHeartbeatRequest) (*ReducerHeartbeatResponse, error) {
	if input != nil {
		return &ReducerHeartbeatResponse{
			Alive: true,
		}, nil
	} else {
		return nil, errors.New("input isn't valid")
	}
}
