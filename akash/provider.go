package akash

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
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
				DefaultFunc: schema.EnvDefaultFunc("AKASH_NET", ""),
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
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("AKASH_HOME", "~/.akash"),
			},
			"path": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("AKASH_PATH", ""),
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

	keyName := d.Get("key_name").(string)
	keyringBackend := d.Get("keyring_backend").(string)
	accountAddress := d.Get("account_address").(string)
	net := d.Get("net").(string)
	version := d.Get("chain_version").(string)
	chainId := d.Get("chain_id").(string)
	node := d.Get("node").(string)
	home := d.Get("home").(string)
	path := d.Get("path").(string)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	if keyName == "" {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to create Akash client",
			Detail:   "Key name was not provided and is not available on the system",
		})

		return nil, diags
	}

	configuration := client.AkashConfiguration{
		KeyName:        keyName,
		KeyringBackend: keyringBackend,
		AccountAddress: accountAddress,
		Net:            net,
		Version:        version,
		ChainId:        chainId,
		Node:           node,
		Home:           home,
		Path:           path,
	}

	tflog.Debug(ctx, fmt.Sprintf("Starting provider with %+v", configuration))

	akash := client.New(ctx, configuration)

	return akash, diags
}
