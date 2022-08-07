package types

type BidsSliceWrapper struct {
	BidWrappers []BidWrapper `json:"bids"`
}

type BidWrapper struct {
	Bid Bid `json:"bid"`
}

type Bids []Bid

type Bid struct {
	Id     BidId `json:"bid_id"`
	Amount int32 `json:"amount"`
}

type BidId struct {
	Provider string `json:"provider"`
}

func (b Bids) GetProviderAddresses() []string {
	addresses := make([]string, 0, len(b))

	for _, bid := range b {
		addresses = append(addresses, bid.Id.Provider)
	}

	return addresses
}
