package cli

func (c AkashCommand) DefaultGas() AkashCommand {
	return c.GasAuto().SetGasAdjustment().SetGasPrices().SetSignMode("amino-json")
}

func (c AkashCommand) DefaultSeqs(dseq string) AkashCommand {
	return c.SetDseq(dseq).SetGseq("1").SetOseq("1")
}
