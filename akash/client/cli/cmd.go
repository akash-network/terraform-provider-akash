package cli

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"os/exec"
	"strings"
)

func (c AkashCommand) AsCmd() *exec.Cmd {
	return exec.Command(
		c.Content[0],
		c.Headless()...,
	)
}

func (c AkashCommand) Raw() ([]byte, error) {
	cmd := c.AsCmd()

	tflog.Debug(c.ctx, strings.Join(cmd.Args, " "))

	var errb bytes.Buffer
	cmd.Stderr = &errb
	out, err := cmd.Output()
	if err != nil {
		tflog.Warn(c.ctx, fmt.Sprintf("Could not execute command: %s", err.Error()))
		if strings.Contains(errb.String(), "error unmarshalling") {
			return c.Raw()
		}

		return nil, errors.New(errb.String())
	}

	return out, nil
}

func (c AkashCommand) DecodeJson(v any) error {
	cmd := c.AsCmd()

	tflog.Debug(c.ctx, strings.Join(cmd.Args, " "))

	var errb bytes.Buffer
	cmd.Stderr = &errb
	out, err := cmd.Output()
	if err != nil {
		tflog.Warn(c.ctx, fmt.Sprintf("Could not execute command: %s", err.Error()))
		if strings.Contains(errb.String(), "error unmarshalling") {
			return c.DecodeJson(v)
		}

		return errors.New(errb.String())
	}

	err = json.NewDecoder(strings.NewReader(string(out))).Decode(v)
	if err != nil {
		tflog.Debug(c.ctx, fmt.Sprintf("Error while unmarshalling command output"))
		return err
	}

	return nil
}
