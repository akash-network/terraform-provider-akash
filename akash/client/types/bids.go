package types

type Bids []Bid

type Bid struct {
	Id BidId `json:"bid_id"`
}

type BidId struct {
	Provider string `json:"provider"`
}
