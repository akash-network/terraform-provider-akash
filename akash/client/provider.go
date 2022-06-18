package client

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"os"
	"os/exec"
	"strings"
	"terraform-provider-akash/akash/client/types"
)

func SendManifest(ctx context.Context, dseq string, provider string, manifestLocation string) (string, error) {
	cmd := exec.CommandContext(
		ctx,
		AkashBinary,
		"provider",
		"send-manifest",
		manifestLocation,
		"--dseq",
		dseq,
		"--provider",
		provider,
		"--home",
		os.Getenv("AKASH_HOME"),
		"--from",
		os.Getenv("AKASH_KEY_NAME"),
		"-o",
		"json",
	)

	tflog.Debug(ctx, strings.Join(cmd.Args, " "))

	var errb bytes.Buffer
	cmd.Stderr = &errb
	out, err := cmd.Output()
	if err != nil {
		return "", errors.New(errb.String())
	}

	tflog.Debug(ctx, fmt.Sprintf("Response content: %s", out))

	return string(out), nil
}

func GetLeaseStatus(ctx context.Context, dseq string, provider string) (*types.LeaseStatus, error) {
	cmd := exec.CommandContext(
		ctx,
		AkashBinary,
		"provider",
		"lease-status",
		"--home",
		os.Getenv("AKASH_HOME"),
		"--dseq",
		dseq,
		"--provider",
		provider,
		"--from",
		os.Getenv("AKASH_KEY_NAME"),
	)

	var errb bytes.Buffer
	cmd.Stderr = &errb
	out, err := cmd.Output()
	if err != nil {
		return nil, errors.New(errb.String())
	}

	leaseStatus := types.LeaseStatus{}
	err = json.NewDecoder(strings.NewReader(string(out))).Decode(&leaseStatus)
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
