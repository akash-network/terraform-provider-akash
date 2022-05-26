package client

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"os/exec"
	"strings"
	"terraform-provider-hashicups/akash/client/types"
	"time"
)

func GetBids(ctx context.Context, dseq string, timeout time.Duration) (types.Bids, error) {
	bids := types.Bids{}
	for timeout > 0 && len(bids) <= 0 {
		startTime := time.Now()
		// Check bids on deployments and choose one
		currentBids, err := queryBidList(dseq)
		if err != nil {
			tflog.Error(ctx, "Failed to query bid list")
			tflog.Debug(ctx, fmt.Sprintf("%s", err))
			return nil, err
		}
		tflog.Debug(ctx, fmt.Sprintf("Received %d bids", len(bids)))
		bids = currentBids
		time.Sleep(time.Second)
		timeout -= time.Since(startTime)
	}

	return bids, nil
}

func queryBidList(dseq string) (types.Bids, error) {
	cmd := exec.Command(
		AKASH_BINARY,
		"query",
		"market",
		"bid",
		"list",
		"--dseq",
		dseq,
		"-o",
		"json",
	)

	var errb bytes.Buffer
	cmd.Stderr = &errb
	out, err := cmd.Output()
	if err != nil {
		return nil, errors.New(errb.String())
	}

	bidsSliceWrapper := types.BidsSliceWrapper{}
	err = json.NewDecoder(strings.NewReader(string(out))).Decode(&bidsSliceWrapper)
	if err != nil {
		return nil, err
	}

	bids := types.Bids{}
	for _, bidWrapper := range bidsSliceWrapper.BidWrappers {
		bids = append(bids, bidWrapper.Bid)
	}

	return bids, nil
}
