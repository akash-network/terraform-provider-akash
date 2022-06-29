package cli

func (c AkashCommand) DefaultGas() AkashCommand {
	return c.GasAuto().GasAdjustment().GasPrices()
}

func (c AkashCommand) DefaultSeqs(dseq string) AkashCommand {
	return c.Dseq(dseq).Gseq("1").Oseq("1")
}
