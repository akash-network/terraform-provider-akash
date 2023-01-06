package filtering

import (
	"errors"
	"terraform-provider-akash/akash/client/types"
)

// FilterPipeline represents a pipeline of filters to be applied
type FilterPipeline struct {
	source   types.Bids
	pipeline []FilterPipe
}

// FilterPipe represents the unit of execution for a filtering logic.
// A FilterPipe takes a types.Bids and returns a types.Bids or error.
type FilterPipe func(types.Bids) (types.Bids, error)

// NewFilterPipeline creates a new FilterPipeline
func NewFilterPipeline(source types.Bids) *FilterPipeline {
	return &FilterPipeline{source: source, pipeline: make([]FilterPipe, 0)}
}

// Pipe adds a FilterPipe to the pipeline.
// It returns a reference to the receiver.
func (fp *FilterPipeline) Pipe(filter FilterPipe) *FilterPipeline {
	fp.pipeline = append(fp.pipeline, filter)
	return fp
}

func (fp *FilterPipeline) Reduce(reducer func(types.Bids) (types.Bid, error)) (types.Bid, error) {
	result, err := fp.Execute()
	if err != nil {
		return types.Bid{}, err
	}

	return reducer(result)
}

// Execute runs the pipeline sequentially and returns the resulting types.Bids or error.
func (fp *FilterPipeline) Execute() (types.Bids, error) {
	buffer := fp.source
	for _, pipe := range fp.pipeline {
		result, err := pipe(buffer)
		if err != nil {
			return nil, err
		}
		buffer = result
	}

	return buffer, nil
}

func Cheapest(bids types.Bids) (types.Bid, error) {
	if len(bids) == 0 {
		return types.Bid{}, errors.New("empty bid slice")
	}

	var cheapestBid types.Bid

	for _, bid := range bids {
		if cheapestBid == (types.Bid{}) || cheapestBid != (types.Bid{}) && bid.Price.Amount < cheapestBid.Price.Amount {
			cheapestBid = bid
		}
	}

	return cheapestBid, nil
}
