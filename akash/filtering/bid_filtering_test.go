package filtering

import (
	"errors"
	"terraform-provider-akash/akash/client/types"
	"testing"
)

func TestFilterPipeline_Execute(t *testing.T) {
	pipeline := NewFilterPipeline(types.Bids{
		types.Bid{Id: types.BidId{Provider: "test1"}, Price: types.BidPrice{Amount: 42}},
		types.Bid{Id: types.BidId{Provider: "test2"}, Price: types.BidPrice{Amount: 24}},
		types.Bid{Id: types.BidId{Provider: "test3"}, Price: types.BidPrice{Amount: 37}},
	})

	t.Run("should return the expected bids", func(t *testing.T) {
		result, _ := pipeline.Pipe(func(bids types.Bids) (types.Bids, error) {
			newBids := make(types.Bids, 0)
			for _, bid := range bids {
				if bid.Price.Amount < 40 {
					newBids = append(newBids, bid)
				}
			}
			return newBids, nil
		}).Execute()

		if len(result) != 2 {
			t.Errorf("Expected bids to have size %d after pipeline run, has size %d instead", 2, len(result))
		}

		if result[0] != pipeline.source[1] || result[1] != pipeline.source[2] {
			t.Errorf("Wrong result from pipeline %+v", result)
		}
	})
}

func TestFilterPipeline_Cheapest(t *testing.T) {

	t.Run("should return the cheapest bid with pipes", func(t *testing.T) {
		pipeline := NewFilterPipeline(types.Bids{
			types.Bid{Id: types.BidId{Provider: "test1"}, Price: types.BidPrice{Amount: 42}},
			types.Bid{Id: types.BidId{Provider: "test2"}, Price: types.BidPrice{Amount: 24}},
			types.Bid{Id: types.BidId{Provider: "test3"}, Price: types.BidPrice{Amount: 37}},
		})

		result, _ := pipeline.Pipe(func(bids types.Bids) (types.Bids, error) {
			newBids := make(types.Bids, 0)
			for _, bid := range bids {
				if bid.Price.Amount > 30 {
					newBids = append(newBids, bid)
				}
			}
			return newBids, nil
		}).Reduce(Cheapest)

		if result != pipeline.source[2] {
			t.Errorf("Expected result to be %+v, got %+v instead", pipeline.source[2], result)
		}
	})

	t.Run("should return the cheapest bid without pipes", func(t *testing.T) {
		pipeline := NewFilterPipeline(types.Bids{
			types.Bid{Id: types.BidId{Provider: "test1"}, Price: types.BidPrice{Amount: 42}},
			types.Bid{Id: types.BidId{Provider: "test2"}, Price: types.BidPrice{Amount: 24}},
			types.Bid{Id: types.BidId{Provider: "test3"}, Price: types.BidPrice{Amount: 37}},
		})

		result, _ := pipeline.Reduce(Cheapest)

		t.Logf("Values in pipeline: %+v", pipeline.source)

		if result != pipeline.source[1] {
			t.Errorf("Expected result to be %+v, got %+v instead", pipeline.source[1], result)
		}
	})

	t.Run("should return error if pipe fails", func(t *testing.T) {
		pipeline := NewFilterPipeline(types.Bids{
			types.Bid{Id: types.BidId{Provider: "test1"}, Price: types.BidPrice{Amount: 42}},
			types.Bid{Id: types.BidId{Provider: "test2"}, Price: types.BidPrice{Amount: 24}},
			types.Bid{Id: types.BidId{Provider: "test3"}, Price: types.BidPrice{Amount: 37}},
		})

		_, err := pipeline.Pipe(func(bids types.Bids) (types.Bids, error) {
			return types.Bids{}, errors.New("failed")
		}).Reduce(Cheapest)

		if err == nil {
			t.Errorf("Expected and error")
		}

		if err.Error() != "failed" {
			t.Errorf("Expected error to be \"failed\", got %s instead", err.Error())
		}
	})
}
