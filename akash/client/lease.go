package client

import (
	"terraform-provider-akash/akash/client/cli"
)

func (ak *AkashClient) CreateLease(dseq string, provider string) (string, error) {
	cmd := cli.AkashCli(ak).Tx().Market().Lease().Create().DefaultSeqs(dseq).
		SetProvider(provider).SetOwner(ak.Config.AccountAddress).SetFrom(ak.Config.KeyName).
		DefaultGas().SetChainId(ak.Config.ChainId).SetKeyringBackend(ak.Config.KeyringBackend).
		AutoAccept().OutputJson()

	out, err := cmd.Raw()
	if err != nil {
		return "", err
	}

	return string(out), nil
}
