package akash

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"os"
	"terraform-provider-akash/akash/client"
)

// Provider -
func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"key_name": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("AKASH_KEY_NAME", ""),
			},
			"keyring_backend": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("AKASH_KEYRING_BACKEND", "os"),
			},
			"account_address": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("AKASH_ACCOUNT_ADDRESS", ""),
			},
			"net": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("AKASH_NET", "akash"),
			},
			"chain_version": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("AKASH_VERSION", ""),
			},
			"chain_id": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("AKASH_CHAIN_ID", ""),
			},
			"node": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("AKASH_NODE", ""),
			},
			"home": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				DefaultFunc: schema.EnvDefaultFunc("AKASH_HOME", func() string {
					homeDir, _ := os.UserHomeDir()
					return homeDir + "/.akash"
				}()),
			},
			"path": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("AKASH_PATH", "akash"),
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"akash_deployment": resourceDeployment(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"akash_deployments": dataSourceDeployments(),
		},
		ConfigureContextFunc: providerConfigure,
	}
}

func providerConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	tflog.Info(ctx, "Configuring the provider")

	config := map[string]string{
		"key_name":        d.Get("key_name").(string),
		"keyring_backend": d.Get("keyring_backend").(string),
		"account_address": d.Get("account_address").(string),
		"net":             d.Get("net").(string),
		"chain_version":   d.Get("chain_version").(string),
		"chain_id":        d.Get("chain_id").(string),
		"node":            d.Get("node").(string),
		"home":            d.Get("home").(string),
		"path":            d.Get("path").(string),
	}

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	if diags, valid := validateConfiguration(diags, config); !valid {
		return nil, diags
	}

	configuration := client.AkashConfiguration{
		KeyName:        config["key_name"],
		KeyringBackend: config["keyring_backend"],
		AccountAddress: config["accountAddress"],
		Net:            config["net"],
		Version:        config["version"],
		ChainId:        config["chainId"],
		Node:           config["node"],
		Home:           config["home"],
		Path:           config["path"],
	}

	tflog.Debug(ctx, fmt.Sprintf("Starting provider with %+v", configuration))

	akash := client.New(ctx, configuration)

	return akash, diags
}

func validateConfiguration(diags diag.Diagnostics, config map[string]string) (diag.Diagnostics, bool) {
	for k, v := range config {
		if v == "" {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Unable to create Akash client",
				Detail:   fmt.Sprintf("Parameter '%s' was not provided and is not available on the system", k),
			})

			return diags, false
		}
	}

	return nil, true
}
