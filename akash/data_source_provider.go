package akash

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"strconv"
	"terraform-provider-akash/akash/client/praetor"
	"terraform-provider-akash/akash/client/types"
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
						},
					},
				},
			},
		},
	}
}

func dataSourceProvidersRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	var providers []types.Provider
	if d.Get("all_providers").(bool) {
		tflog.Info(ctx, "All providers requested")
		providers = praetor.GetAllProviders()
	} else {
		tflog.Info(ctx, "Active providers requested")
		providers = praetor.GetActiveProviders()
	}
	tflog.Info(ctx, fmt.Sprintf("Got %d providers", len(providers)))
	akashProviders := make([]map[string]interface{}, 0, len(providers))
	for _, p := range providers {
		if p.Uptime < float32(d.Get("minimum_uptime").(float64)) {
			tflog.Debug(ctx, fmt.Sprintf("Provider %s uptime (%.2f) is below minimum", p.Address, p.Uptime))
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
