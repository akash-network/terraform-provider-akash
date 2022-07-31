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

const KeyName = "key_name"
const KeyringBackend = "keyring_backend"
const AccountAddress = "account_address"
const Net = "net"
const ChainVersion = "chain_version"
const ChainId = "chain_id"
const Node = "node"
const Home = "home"
const Path = "path"

// Provider -
func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			KeyName: &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("AKASH_KEY_NAME", ""),
			},
			KeyringBackend: &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("AKASH_KEYRING_BACKEND", "os"),
			},
			AccountAddress: &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("AKASH_ACCOUNT_ADDRESS", ""),
			},
			Net: &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("AKASH_NET", "akash"),
			},
			ChainVersion: &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("AKASH_VERSION", ""),
			},
			ChainId: &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("AKASH_CHAIN_ID", ""),
			},
			Node: &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("AKASH_NODE", ""),
			},
			Home: &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				DefaultFunc: schema.EnvDefaultFunc("AKASH_HOME", func() string {
					homeDir, _ := os.UserHomeDir()
					return homeDir + "/.akash"
				}()),
			},
			Path: &schema.Schema{
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
		KeyName:        d.Get(KeyName).(string),
		KeyringBackend: d.Get(KeyringBackend).(string),
		AccountAddress: d.Get(AccountAddress).(string),
		Net:            d.Get(Net).(string),
		ChainVersion:   d.Get(ChainVersion).(string),
		ChainId:        d.Get(ChainId).(string),
		Node:           d.Get(Node).(string),
		Home:           d.Get(Home).(string),
		Path:           d.Get(Path).(string),
	}

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	if diags, valid := validateConfiguration(diags, config); !valid {
		return nil, diags
	}

	configuration := client.AkashConfiguration{
		KeyName:        config[KeyName],
		KeyringBackend: config[KeyringBackend],
		AccountAddress: config[AccountAddress],
		Net:            config[Net],
		Version:        config[ChainVersion],
		ChainId:        config[ChainId],
		Node:           config[Node],
		Home:           config[Home],
		Path:           config[Path],
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
