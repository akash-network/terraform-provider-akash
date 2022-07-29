package client_test

import (
	"context"
	"terraform-provider-akash/akash/client"
	"terraform-provider-akash/akash/client/types"
	"testing"
)

func TestFindCheapestReturnsErrorOnEmptyBidsList(t *testing.T) {
	akash := client.New(context.TODO(), client.AkashConfiguration{})
	provider, err := akash.FindCheapest(types.Bids{})
	expectedError := "empty bid slice"

	if provider != "" {
		t.Logf("SetProvider should be empty string, is \"%s\" instead", provider)
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
