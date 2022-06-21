package cli

import "os/exec"

func (c AkashCommand) AsCmd() *exec.Cmd {
	return exec.Command(
		AkashCli()[0],
		c.Headless()...,
	)
}
