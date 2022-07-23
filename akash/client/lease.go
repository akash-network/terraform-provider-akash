package client

import (
	"os"
	"terraform-provider-akash/akash/client/cli"
)

func (ak *AkashClient) CreateLease(dseq string, provider string) (string, error) {
	cmd := cli.AkashCli(ak.ctx).Tx().Market().Lease().Create().DefaultSeqs(dseq).
		SetProvider(provider).SetOwner(os.Getenv("AKASH_ACCOUNT_ADDRESS")).SetFrom(os.Getenv("AKASH_KEY_NAME")).
		DefaultGas().AutoAccept().OutputJson()

	out, err := cmd.Raw()
	if err != nil {
		return "", err
	}

	return string(out), nil
}
