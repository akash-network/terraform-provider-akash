package client

import (
	"context"
	"errors"
	"fmt"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"os"
	"strings"
	"terraform-provider-akash/akash/client/cli"
	"terraform-provider-akash/akash/client/types"
)

func SendManifest(ctx context.Context, dseq string, provider string, manifestLocation string) (string, error) {

	cmd := cli.AkashCli().Provider().SendManifest(manifestLocation).
		SetDseq(dseq).SetProvider(provider).SetHome(os.Getenv("AKASH_HOME")).
		SetFrom(os.Getenv("AKASH_KEY_NAME")).OutputJson()

	tflog.Debug(ctx, strings.Join(cmd.AsCmd().Args, " "))

	out, err := cmd.Raw()
	if err != nil {
		return "", err
	}

	tflog.Debug(ctx, fmt.Sprintf("Response content: %s", out))

	return string(out), nil
}

func GetLeaseStatus(ctx context.Context, dseq string, provider string) (*types.LeaseStatus, error) {

	cmd := cli.AkashCli().Provider().LeaseStatus().
		SetHome(os.Getenv("AKASH_HOME")).SetDseq(dseq).SetProvider(provider).
		SetFrom(os.Getenv("AKASH_KEY_NAME"))

	tflog.Info(ctx, strings.Join(cmd.AsCmd().Args, " "))

	leaseStatus := types.LeaseStatus{}
	err := cmd.DecodeJson(&leaseStatus)
	if err != nil {
		return nil, err
	}

	return &leaseStatus, nil
}

func FindCheapest(ctx context.Context, bids types.Bids) (string, error) {
	if len(bids) == 0 {
		tflog.Error(ctx, "Empty bid slice")
		return "", errors.New("empty bid slice")
	}

	var cheapestBid *types.Bid = nil

	tflog.Info(ctx, "Finding cheapest bid")

	for _, bid := range bids {
		if cheapestBid == nil || cheapestBid != nil && bid.Amount < cheapestBid.Amount {
			cheapestBid = &bid
		}
	}

	return cheapestBid.Id.Provider, nil
}
