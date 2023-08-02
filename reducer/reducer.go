package reducer

import (
	"PageRank/constants"
	"context"
	"errors"
	"log"
)

type Reducer struct {
}

func (reducer *Reducer) mustEmbedUnimplementedReducerServer() {
	//TODO implement me
	panic("implement me")
}

// Reduce ->  sum aggregated page rank values and return the value
func (reducer *Reducer) Reduce(ctx context.Context, input *ReducerInput) (*ReducerOutput, error) {

	if input != nil {
		accumulator := float32(0.0)
		for _, rank := range input.PageRankShares {
			accumulator += rank
		}

		newPageRank := float32((1.0-constants.DampingFactor)/float64(input.GraphSize)) + constants.DampingFactor*accumulator
		log.Printf("NodeID: %s evaluated page rank: %f\n", input.NodeId, newPageRank)

		return &ReducerOutput{
			NodeId:       input.NodeId,
			NewRankValue: newPageRank,
		}, nil
	} else {

		return nil, errors.New("Input isn't valid")
	}

}
