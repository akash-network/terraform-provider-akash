package cli

func (c AkashCommand) DefaultGas() AkashCommand {
	return c.GasAuto().SetGasAdjustment().SetGasPrices().SetSignMode("amino-json")
}

func (c AkashCommand) SetSeqs(dseq string, gseq string, oseq string) AkashCommand {
	return c.SetDseq(dseq).SetGseq(gseq).SetOseq(oseq)
}
