package cli_test

import (
	"context"
	"terraform-provider-akash/akash/client"
	"terraform-provider-akash/akash/client/cli"
	"testing"
)

func TestAkashCliAppendsSubcommands(t *testing.T) {
	cmd := cli.AkashCli(client.New(context.TODO(), client.AkashConfiguration{})).Query().Market().Bid().List()
	expectedSize := 5

	if len(cmd.Content) != expectedSize {
		t.Logf("Expected command to have %d subcommands, found %d (%+v)", expectedSize, len(cmd.Content), cmd)
		t.Fail()
	}
}

func TestAkashCliHeadlessSizeIsCorrect(t *testing.T) {
	cmd := cli.AkashCli(client.New(context.TODO(), client.AkashConfiguration{})).Query().Market().Bid().List().Headless()
	expectedSize := 4

	if len(cmd) != expectedSize {
		t.Logf("Expected command to have %d subcommands, found %d (%+v)", expectedSize, len(cmd), cmd)
		t.Fail()
	}
}
