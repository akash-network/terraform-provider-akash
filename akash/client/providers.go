package client

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"os"
	"os/exec"
)

func SendManifest(ctx context.Context, dseq string, provider string, manifestLocation string) (string, error) {
	cmd := exec.CommandContext(
		ctx,
		AKASH_BINARY,
		"provider",
		"send-manifest",
		manifestLocation,
		"--dseq",
		dseq,
		"--provider",
		provider,
		"--home",
		os.Getenv("AKASH_HOME"),
		"--from",
		os.Getenv("AKASH_KEY_NAME"),
		"-o",
		"json",
	)

	var errb bytes.Buffer
	cmd.Stderr = &errb
	out, err := cmd.Output()
	if err != nil {
		return "", errors.New(errb.String())
	}

	tflog.Debug(ctx, fmt.Sprintf("Response contect: %s", out))

	return string(out), nil
}
