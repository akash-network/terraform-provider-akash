package client

import (
	"bytes"
	"context"
	"errors"
	"os"
	"os/exec"
)

func CreateLease(ctx context.Context, dseq string, provider string) error {
	return transactionCreateLease(dseq, provider)
}

func transactionCreateLease(dseq string, provider string) error {
	cmd := exec.Command(
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
		"-o",
		"json",
	)

	var errb bytes.Buffer
	cmd.Stderr = &errb
	_, err := cmd.Output()
	if err != nil {
		return errors.New(errb.String())
	}

	return nil
}
