package client

import (
	"context"
	"terraform-provider-akash/akash/client/types"
	"testing"
)

func TestFindCheapestReturnsErrorOnEmptyBidsList(t *testing.T) {
	ctx := context.TODO()
	provider, err := FindCheapest(ctx, types.Bids{})
	expectedError := "empty bid slice"

	if provider != "" {
		t.Logf("Provider should be empty string, is \"%s\" instead", provider)
		t.Fail()
	}

	if err == nil {
		t.Logf("Should have returned an error, returned nil instead")
		t.Fail()
	}

	if err.Error() != expectedError {
		t.Logf("Error should be \"%s\", is \"%s\" instead", expectedError, err.Error())
		t.Fail()
	}
}
