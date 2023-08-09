package mapper

import (
	"context"
	"log"
)

type Mapper struct {
}

func (mapper *Mapper) mustEmbedUnimplementedMapperServer() {
	//TODO implement me
	panic("implement me")
}

// Map ->  function to evaluate each node contributes
func (mapper *Mapper) Map(ctx context.Context, input *MapperInput) (*MapperOutput, error) {
	numOutLinks := float32(len(input.AdjacencyList))
	if numOutLinks > 0 {
		pagerankShare := input.PageRank / numOutLinks

		log.Printf("Share: %f for nodes: %v\n", pagerankShare, input.GetAdjacencyList())

		return &MapperOutput{
			PageRankShare: pagerankShare,
			AdjacencyList: input.GetAdjacencyList(),
		}, nil

	} else {
		//If here, node has zero out-links
		zero := 0.0
		return &MapperOutput{
			PageRankShare: float32(zero),
			AdjacencyList: input.GetAdjacencyList(),
		}, nil
	}
}

// CleanUp -> to manage sinks and random jump factor
func (mapper *Mapper) CleanUp(ctx context.Context, input *CleanUpInput) (*CleanUpOutput, error) {
	numOutLinks := float32(len(input.AdjacencyList))
	if numOutLinks == 0 {
		log.Printf("Sink finded")
		return &CleanUpOutput{
			SinkMass: input.PageRank,
		}, nil
	} else {
		//If here, node has out-links, not interesting in clean-up phase
		zero := 0.0
		return &CleanUpOutput{
			SinkMass: float32(zero),
		}, nil
	}
}
