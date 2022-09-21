package cli

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

type AkashCommand struct {
	ctx     context.Context
	Content []string
}

type AkashCliClient interface {
	GetContext() context.Context
	GetPath() string
}

func AkashCli(client AkashCliClient) AkashCommand {
	path := client.GetPath()
	if path == "" {
		path = "akash"
	}

	return AkashCommand{
		ctx:     client.GetContext(),
		Content: []string{path},
	}
}

func (c AkashCommand) Tx() AkashCommand {
	return c.append("tx")
}

func (c AkashCommand) Deployment() AkashCommand {
	return c.append("deployment")
}

func (c AkashCommand) Get() AkashCommand {
	return c.append("get")
}

func (c AkashCommand) Create() AkashCommand {
	return c.append("create")
}

func (c AkashCommand) Update() AkashCommand {
	return c.append("update")
}

func (c AkashCommand) LeaseStatus() AkashCommand {
	return c.append("lease-status")
}

func (c AkashCommand) SendManifest(path string) AkashCommand {
	return c.append("send-manifest").append(path)
}

func (c AkashCommand) Close() AkashCommand {
	return c.append("close")
}

func (c AkashCommand) Query() AkashCommand {
	return c.append("query")
}

func (c AkashCommand) Market() AkashCommand {
	return c.append("market")
}

func (c AkashCommand) Provider() AkashCommand {
	return c.append("provider")
}

func (c AkashCommand) Bid() AkashCommand {
	return c.append("bid")
}

func (c AkashCommand) List() AkashCommand {
	return c.append("list")
}

func (c AkashCommand) Lease() AkashCommand {
	return c.append("lease")
}

func (c AkashCommand) Manifest(path string) AkashCommand {
	return c.append(path)
}

/** OPTIONS **/

func (c AkashCommand) SetDseq(dseq string) AkashCommand {
	return c.append("--dseq").append(dseq)
}

func (c AkashCommand) SetOseq(oseq string) AkashCommand {
	return c.append("--oseq").append(oseq)
}

func (c AkashCommand) SetGseq(gseq string) AkashCommand {
	return c.append("--gseq").append(gseq)
}

func (c AkashCommand) SetProvider(provider string) AkashCommand {
	return c.append("--provider").append(provider)
}

func (c AkashCommand) SetHome(home string) AkashCommand {
	return c.append("--home").append(home)
}

func (c AkashCommand) SetOwner(owner string) AkashCommand {
	return c.append("--owner").append(owner)
}

func (c AkashCommand) SetFees(amount int64) AkashCommand {
	return c.append("--fees").append(fmt.Sprintf("%duakt", amount))
}

func (c AkashCommand) SetFrom(key string) AkashCommand {
	return c.append("--from").append(key)
}

func (c AkashCommand) GasAuto() AkashCommand {
	return c.append("--gas=auto")
}
func (c AkashCommand) SetGasAdjustment(adjustment float32) AkashCommand {
	return c.append(fmt.Sprintf("--gas-adjustment=%2f", adjustment))
}

func (c AkashCommand) SetGasPrices() AkashCommand {
	return c.append("--gas-prices=0.025uakt")
}

func (c AkashCommand) SetChainId(chainId string) AkashCommand {
	return c.append("--chain-id").append(chainId)
}

func (c AkashCommand) SetNode(node string) AkashCommand {
	return c.append("--node").append(node)
}

func (c AkashCommand) SetKeyringBackend(keyringBackend string) AkashCommand {
	return c.append("--keyring-backend").append(keyringBackend)
}

func (c AkashCommand) SetSignMode(mode string) AkashCommand {
	supportedModes := map[string]bool{
		"default":    true,
		"amino-json": true,
	}

	if _, ok := supportedModes[mode]; !ok {
		tflog.Error(c.ctx, fmt.Sprintf("Mode '%s' not supported", mode))
		return c
	}

	return c.append("--sign-mode").append(mode)
}

func (c AkashCommand) AutoAccept() AkashCommand {
	return c.append("-y")
}

func (c AkashCommand) OutputJson() AkashCommand {
	return c.append("-o").append("json")
}

func (c AkashCommand) Headless() []string {
	return c.Content[1:]
}

func (c AkashCommand) append(str string) AkashCommand {
	c.Content = append(c.Content, str)
	return c
}
