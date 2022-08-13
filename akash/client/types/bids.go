package types

type BidsSliceWrapper struct {
	BidWrappers []BidWrapper `json:"bids"`
}

type BidWrapper struct {
	Bid Bid `json:"bid"`
}

type Bids []Bid

type Bid struct {
	Id    BidId    `json:"bid_id"`
	Price BidPrice `json:"price"`
}

type BidId struct {
	Provider string `json:"provider"`
}

type BidPrice struct {
	Amount float32 `json:"amount,string"`
}

func (b Bids) GetProviderAddresses() []string {
	addresses := make([]string, 0, len(b))

	for _, bid := range b {
		addresses = append(addresses, bid.Id.Provider)
	}

	return addresses
}

func (b Bids) FindByProvider(provider string) Bid {
	for _, bid := range b {
		if bid.Id.Provider == provider {
			return bid
		}
	}

	return Bid{}
}

// FindAllByProviders finds all the Bid structures that have any of the given providers.
// It returns a slice of all the Bid structures where the providers were found.
func (b Bids) FindAllByProviders(providers []string) Bids {
	bids := make(Bids, 0)

	for _, provider := range providers {
		if bid := b.FindByProvider(provider); bid != (Bid{}) {
			bids = append(bids, bid)
		}
	}

	return bids
}
