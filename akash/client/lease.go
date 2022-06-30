package client

import (
	"context"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"os"
	"strings"
	"terraform-provider-akash/akash/client/cli"
)

func CreateLease(ctx context.Context, dseq string, provider string) (string, error) {
	cmd := cli.AkashCli().Tx().Market().Lease().Create().DefaultSeqs(dseq).
		SetProvider(provider).SetOwner(os.Getenv("AKASH_ACCOUNT_ADDRESS")).SetFrom(os.Getenv("AKASH_KEY_NAME")).
		DefaultGas().AutoAccept().OutputJson()

	tflog.Info(ctx, strings.Join(cmd.AsCmd().Args, " "))

	out, err := cmd.Raw()
	if err != nil {
		if strings.Contains(err.Error(), "error unmarshalling") {
			return CreateLease(ctx, dseq, provider)
		}

		return "", err
	}

	return string(out), nil
}
