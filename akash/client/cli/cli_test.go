package cli

import "testing"

func TestAkashCliAppendsSubcommands(t *testing.T) {
	cmd := AkashCli().Query().Market().Bid().List()
	expectedSize := 5

	if len(cmd) != expectedSize {
		t.Logf("Expected command to have %d subcommands, found %d (%+v)", expectedSize, len(cmd), cmd)
		t.Fail()
	}
}

func TestAkashCliHeadlessSizeIsCorrect(t *testing.T) {
	cmd := AkashCli().Query().Market().Bid().List().Headless()
	expectedSize := 4

	if len(cmd) != expectedSize {
		t.Logf("Expected command to have %d subcommands, found %d (%+v)", expectedSize, len(cmd), cmd)
		t.Fail()
	}
}
