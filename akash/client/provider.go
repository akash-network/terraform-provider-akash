package client

import (
	"errors"
	"fmt"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"os"
	"strings"
	"terraform-provider-akash/akash/client/cli"
	"terraform-provider-akash/akash/client/types"
)

func (ak *AkashClient) SendManifest(dseq string, provider string, manifestLocation string) (string, error) {

	cmd := cli.AkashCli(ak.ctx).Provider().SendManifest(manifestLocation).
		SetDseq(dseq).SetProvider(provider).SetHome(os.Getenv("AKASH_HOME")).
		SetFrom(os.Getenv("AKASH_KEY_NAME")).OutputJson()

	tflog.Debug(ak.ctx, strings.Join(cmd.AsCmd().Args, " "))

	out, err := cmd.Raw()
	if err != nil {
		return "", err
	}

	tflog.Debug(ak.ctx, fmt.Sprintf("Response content: %s", out))

	return string(out), nil
}

func (ak *AkashClient) GetLeaseStatus(dseq string, provider string) (*types.LeaseStatus, error) {

	cmd := cli.AkashCli(ak.ctx).Provider().LeaseStatus().
		SetHome(os.Getenv("AKASH_HOME")).SetDseq(dseq).SetProvider(provider).
		SetFrom(os.Getenv("AKASH_KEY_NAME"))

	leaseStatus := types.LeaseStatus{}
	err := cmd.DecodeJson(&leaseStatus)
	if err != nil {
		return nil, err
	}

	return &leaseStatus, nil
}

func (ak *AkashClient) FindCheapest(bids types.Bids) (string, error) {
	if len(bids) == 0 {
		tflog.Error(ak.ctx, "Empty bid slice")
		return "", errors.New("empty bid slice")
	}

	var cheapestBid *types.Bid = nil

	tflog.Info(ak.ctx, "Finding cheapest bid")

	for _, bid := range bids {
		if cheapestBid == nil || cheapestBid != nil && bid.Amount < cheapestBid.Amount {
			cheapestBid = &bid
		}
	}

	return cheapestBid.Id.Provider, nil
}
