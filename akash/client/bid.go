package client

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"terraform-provider-akash/akash/client/cli"
	"terraform-provider-akash/akash/client/types"
	"time"
)

func (ak *AkashClient) GetBids(seqs Seqs, timeout time.Duration) (types.Bids, error) {
	bids := types.Bids{}
	for timeout > 0 && len(bids) <= 0 {
		startTime := time.Now()
		// Check bids on deployments and choose one
		currentBids, err := queryBidList(ak, seqs)
		if err != nil {
			tflog.Error(ak.ctx, "Failed to query bid list")
			tflog.Debug(ak.ctx, err.Error())

			return nil, err
		}
		tflog.Debug(ak.ctx, fmt.Sprintf("Received %d bids", len(bids)))
		bids = currentBids
		time.Sleep(time.Second)
		timeout -= time.Since(startTime)
	}

	return bids, nil
}

func queryBidList(ak *AkashClient, seqs Seqs) (types.Bids, error) {
	cmd := cli.AkashCli(ak).Query().Market().Bid().List().
		SetDseq(seqs.Dseq).SetGseq(seqs.Gseq).SetOseq(seqs.Oseq).
		SetChainId(ak.Config.ChainId).OutputJson()

	bidsSliceWrapper := types.BidsSliceWrapper{}
	if err := cmd.DecodeJson(&bidsSliceWrapper); err != nil {
		return nil, err
	}

	bids := types.Bids{}
	for _, bidWrapper := range bidsSliceWrapper.BidWrappers {
		bids = append(bids, bidWrapper.Bid)
	}

	return bids, nil
}
