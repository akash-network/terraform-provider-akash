package akash

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"strconv"
	"terraform-provider-akash/akash/client"
	"terraform-provider-akash/akash/client/providers-api"
	"terraform-provider-akash/akash/client/types"
	"terraform-provider-akash/akash/extensions"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceProviders() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceProvidersRead,
		Schema: map[string]*schema.Schema{
			"all_providers": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"minimum_uptime": &schema.Schema{
				Type:     schema.TypeFloat,
				Optional: true,
				Default:  0,
			},
			"providers": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"address": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"active": &schema.Schema{
							Type:     schema.TypeBool,
							Computed: true,
						},
						"uptime": &schema.Schema{
							Type:     schema.TypeFloat,
							Computed: true,
						},
						"attributes": &schema.Schema{
							Type:     schema.TypeMap,
							Computed: true,
							Optional: true,
						},
					},
				},
			},
			"required_attributes": &schema.Schema{
				Type:     schema.TypeMap,
				Optional: true,
			},
		},
	}
}

func dataSourceProvidersRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics
	akash := m.(*client.AkashClient)

	providersClient := providers_api.New(akash.Config.ProvidersApi)

	var providers []types.Provider
	if d.Get("all_providers").(bool) {
		tflog.Info(ctx, "All providers requested")
		providers, _ = providersClient.GetAllProviders()
	} else {
		tflog.Info(ctx, "Active providers requested")
		providers, _ = providersClient.GetActiveProviders()
	}

	if attributes, ok := d.GetOk("required_attributes"); ok {
		attrs := attributes.(map[string]interface{})
		if safeAttributes, err := extensions.SafeCastMapValues[string, string](attrs); err == nil {
			providers = getProvidersWithAttributes(providers, safeAttributes)
		} else {
			tflog.Error(ctx, fmt.Sprintf("Could not cast required_attributes: %s", err))
			return diag.FromErr(err)
		}
	} else {
		tflog.Info(ctx, "No required_attributes provided")
	}

	tflog.Info(ctx, fmt.Sprintf("Got %d providers", len(providers)))
	akashProviders := make([]map[string]interface{}, 0, len(providers))
	for _, p := range providers {
		if p.Uptime < float32(d.Get("minimum_uptime").(float64)) {
			tflog.Debug(ctx, fmt.Sprintf("provider %s uptime (%.2f) is below minimum", p.Address, p.Uptime))
		} else {
			akashProviders = append(akashProviders, map[string]interface{}{
				"address":    p.Address,
				"active":     p.Active,
				"uptime":     p.Uptime,
				"attributes": p.Attributes,
			})
		}
	}
	if err := d.Set("providers", akashProviders); err != nil {
		return diag.FromErr(err)
	}

	// always run
	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))

	return diags
}

func getProvidersWithAttributes(providers []types.Provider, requiredAttributes map[string]string) []types.Provider {
	newProviders := make([]types.Provider, 0, len(providers))

	for _, p := range providers {
		if extensions.IsSubset(p.Attributes, requiredAttributes) {
			newProviders = append(newProviders, p)
		}
	}

	return newProviders
}
