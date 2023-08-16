package reducer

import (
	"Reducer/constants"
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
	var config constants.Config
	constants.ReadJsonConfig(&config)
	if input != nil {
		accumulator := float32(0.0)
		for _, rank := range input.PageRankShares {
			accumulator += rank
		}

		newPageRank := (1.0-config.DampingFactor)/float32(input.GraphSize) + config.DampingFactor*accumulator
		log.Printf("\nNodeID: %d -> evaluated page rank: %f", input.NodeId, newPageRank)

		return &ReducerOutput{
			NodeId:       input.NodeId,
			NewRankValue: newPageRank,
		}, nil
	} else {

		return nil, errors.New("input isn't valid")
	}

}

// ReduceCleanUp -> use the cleanUp formula to fix page rank value
func (reducer *Reducer) ReduceCleanUp(ctx context.Context, input *ReducerCleanUpInput) (*ReducerOutput, error) {
	var config constants.Config
	constants.ReadJsonConfig(&config)
	if input != nil {
		massShares := input.SinkMass / float32(input.GraphSize)
		newPageRank := (1.0-config.DampingFactor)/float32(input.GraphSize) + config.DampingFactor*(input.CurrentPageRank+massShares)
		return &ReducerOutput{
			NodeId:       input.NodeId,
			NewRankValue: newPageRank,
		}, nil
	} else {

		return nil, errors.New("input isn't valid")
	}
}
