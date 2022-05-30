package client

import (
	"bytes"
	"context"
	"errors"
	"os"
	"os/exec"
)

func CreateLease(ctx context.Context, dseq string, provider string) (string, error) {
	cmd := exec.CommandContext(
		ctx,
		AKASH_BINARY,
		"tx",
		"market",
		"lease",
		"create",
		"--dseq",
		dseq,
		"--gseq",
		"1",
		"--oseq",
		"1",
		"--provider",
		provider,
		"--owner",
		os.Getenv("AKASH_ACCOUNT_ADDRESS"),
		"--from",
		os.Getenv("AKASH_KEY_NAME"),
		"--gas=\"auto\"",
		"--gas-adjustment=1.15",
		"--gas-prices=\"0.025uakt\"",
		"-y",
		"-o",
		"json",
	)

	var errb bytes.Buffer
	cmd.Stderr = &errb
	out, err := cmd.Output()
	if err != nil {
		return "", errors.New(errb.String())
	}

	return string(out), nil
}
