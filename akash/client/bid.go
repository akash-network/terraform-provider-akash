package client

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"strings"
	"terraform-provider-akash/akash/client/cli"
	"terraform-provider-akash/akash/client/types"
	"time"
)

func GetBids(ctx context.Context, dseq string, timeout time.Duration) (types.Bids, error) {
	bids := types.Bids{}
	for timeout > 0 && len(bids) <= 0 {
		startTime := time.Now()
		// Check bids on deployments and choose one
		currentBids, err := queryBidList(ctx, dseq)
		if err != nil {
			tflog.Error(ctx, "Failed to query bid list")
			tflog.Debug(ctx, fmt.Sprintf("Error: %s", err))

			if strings.Contains(err.Error(), "error unmarshalling") {
				continue
			}
			return nil, err
		}
		tflog.Debug(ctx, fmt.Sprintf("Received %d bids", len(bids)))
		bids = currentBids
		time.Sleep(time.Second)
		timeout -= time.Since(startTime)
	}

	return bids, nil
}

func queryBidList(ctx context.Context, dseq string) (types.Bids, error) {
	cmd := cli.AkashCli().Query().Market().Bid().List().Dseq(dseq).OutputJson()

	tflog.Debug(ctx, strings.Join(cmd.AsCmd().Args, " "))

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
