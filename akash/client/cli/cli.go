package cli

type AkashCommand []string

func AkashCli() AkashCommand {
	return AkashCommand{"../bin/akash"}
}

func (c AkashCommand) Tx() AkashCommand {
	return c.append("tx")
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

func (c AkashCommand) OutputJson() AkashCommand {
	return c.append("-o").append("json")
}

func (c AkashCommand) Dseq(dseq string) AkashCommand {
	return c.append("--dseq").append(dseq)
}

func (c AkashCommand) Headless() []string {
	return c[1:]
}

func (c AkashCommand) append(str string) AkashCommand {
	return append(c, str)
}
