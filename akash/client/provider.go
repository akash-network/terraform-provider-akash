package client

import (
	"errors"
	"fmt"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"terraform-provider-akash/akash/client/cli"
	"terraform-provider-akash/akash/client/types"
)

func (ak *AkashClient) SendManifest(dseq string, provider string, manifestLocation string) (string, error) {

	cmd := cli.AkashCli(ak).Provider().SendManifest(manifestLocation).
		SetDseq(dseq).SetProvider(provider).SetHome(ak.Config.Home).
		SetKeyringBackend(ak.Config.KeyringBackend).SetFrom(ak.Config.KeyName).
		SetNode(ak.Config.Node).OutputJson()

	out, err := cmd.Raw()
	if err != nil {
		return "", err
	}

	tflog.Debug(ak.ctx, fmt.Sprintf("Response content: %s", out))

	return string(out), nil
}

func (ak *AkashClient) GetLeaseStatus(seqs Seqs, provider string) (*types.LeaseStatus, error) {

	cmd := cli.AkashCli(ak).Provider().LeaseStatus().
		SetHome(ak.Config.Home).SetDseq(seqs.Dseq).SetGseq(seqs.Gseq).SetOseq(seqs.Oseq).
		SetNode(ak.Config.Node).SetProvider(provider).SetFrom(ak.Config.KeyName)

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
