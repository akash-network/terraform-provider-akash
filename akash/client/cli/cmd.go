package cli

import (
	"bytes"
	"encoding/json"
	"errors"
	"os/exec"
	"strings"
)

func (c AkashCommand) AsCmd() *exec.Cmd {
	return exec.Command(
		AkashCli()[0],
		c.Headless()...,
	)
}

func (c AkashCommand) Raw() ([]byte, error) {
	cmd := c.AsCmd()

	var errb bytes.Buffer
	cmd.Stderr = &errb
	out, err := cmd.Output()
	if err != nil {
		return nil, errors.New(errb.String())
	}

	return out, nil
}

func (c AkashCommand) DecodeJson(v any) error {
	cmd := c.AsCmd()

	var errb bytes.Buffer
	cmd.Stderr = &errb
	out, err := cmd.Output()
	if err != nil {
		return errors.New(errb.String())
	}

	err = json.NewDecoder(strings.NewReader(string(out))).Decode(v)
	if err != nil {
		return err
	}

	return nil
}