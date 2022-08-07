package akash

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"testing"
)

var testAccProviders map[string]*schema.Provider
var testAccProvider *schema.Provider

func init() {
	testAccProvider = Provider()
	testAccProviders = map[string]*schema.Provider{
		"hashicups": testAccProvider,
	}
}

func TestProvider(t *testing.T) {
	if err := Provider().InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestProviderValidate(t *testing.T) {
	provider := Provider()

	diags := provider.Validate(&terraform.ResourceConfig{Config: map[string]interface{}{
		KeyName:        "test",
		KeyringBackend: "test",
		AccountAddress: "test",
		ChainId:        "test",
		ChainVersion:   "1",
		Net:            "test",
		Node:           "test",
		Home:           "test",
		Path:           "test",
	}})

	if diags.HasError() {
		t.Errorf("Did not expect error. Got: %+v", diags)
	}
}

func TestProviderConfigure(t *testing.T) {
	provider := Provider()
	config := terraform.ResourceConfig{Config: map[string]interface{}{
		KeyName:        "test",
		KeyringBackend: "test",
		AccountAddress: "test",
		ChainId:        "test",
		ChainVersion:   "1",
		Net:            "test",
		Node:           "test",
		Home:           "test",
		Path:           "test",
	}}

	t.Run("should configure with valid configuration", func(t *testing.T) {
		testConfig := config
		diags := provider.Configure(context.TODO(), &testConfig)

		if diags.HasError() {
			t.Errorf("Did not expect error. Got: %+v", diags)
		}
	})

	testCases := []struct {
		requiredField string
	}{
		{KeyName},
		{AccountAddress},
		{ChainVersion},
		{ChainId},
		{Node},
	}

	for _, testCase := range testCases {
		t.Run(fmt.Sprintf("should fail with empty %s", testCase.requiredField), func(t *testing.T) {
			testConfig := config

			testConfig.Config[testCase.requiredField] = ""
			diags := provider.Configure(context.TODO(), &testConfig)

			if !diags.HasError() {
				t.Errorf("Expected an error but got nothing")
			}

			if len(diags) != 1 {
				t.Logf("Errors: %+v", diags)
				t.Errorf("Expected 1 error, got %d", len(diags))
			}
		})

		t.Run(fmt.Sprintf("should fail with unset %s", testCase.requiredField), func(t *testing.T) {
			testConfig := config

			delete(testConfig.Config, testCase.requiredField)
			diags := provider.Configure(context.TODO(), &testConfig)

			if !diags.HasError() {
				t.Errorf("Expected an error but got nothing")
			}

			if len(diags) != 1 {
				t.Logf("Errors: %+v", diags)
				t.Errorf("Expected 1 error, got %d", len(diags))
			}
		})
	}
}
