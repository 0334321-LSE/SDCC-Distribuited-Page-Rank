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
