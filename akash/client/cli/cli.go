package cli

import "fmt"

type AkashCommand []string

func AkashCli() AkashCommand {
	return AkashCommand{"../bin/akash"}
}

func (c AkashCommand) Tx() AkashCommand {
	return c.append("tx")
}

func (c AkashCommand) Deployment() AkashCommand {
	return c.append("deployment")
}

func (c AkashCommand) Create() AkashCommand {
	return c.append("create")
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

func (c AkashCommand) OutputJson() AkashCommand {
	return c.append("-o").append("json")
}

func (c AkashCommand) Dseq(dseq string) AkashCommand {
	return c.append("--dseq").append(dseq)
}

func (c AkashCommand) Oseq(oseq string) AkashCommand {
	return c.append("--oseq").append(oseq)
}

func (c AkashCommand) Gseq(gseq string) AkashCommand {
	return c.append("--gseq").append(gseq)
}

func (c AkashCommand) Provider(provider string) AkashCommand {
	return c.append("--provider").append(provider)
}

func (c AkashCommand) Owner(owner string) AkashCommand {
	return c.append("--owner").append(owner)
}

func (c AkashCommand) Fees(amount int64) AkashCommand {
	return c.append("--fees").append(fmt.Sprintf("%duakt", amount))
}

func (c AkashCommand) AutoAccept() AkashCommand {
	return c.append("-y")
}

func (c AkashCommand) From(key string) AkashCommand {
	return c.append("--from").append(key)
}

func (c AkashCommand) GasAuto() AkashCommand {
	return c.append("--gas=auto")
}
func (c AkashCommand) GasAdjustment() AkashCommand {
	return c.append("--gas-adjustment=1.15")
}

func (c AkashCommand) GasPrices() AkashCommand {
	return c.append("--gas-prices=0.025uakt")
}

func (c AkashCommand) Headless() []string {
	return c[1:]
}

func (c AkashCommand) append(str string) AkashCommand {
	return append(c, str)
}
