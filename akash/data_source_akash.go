package akash

import (
	"context"
	"strconv"
	"terraform-provider-akash/akash/client"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceDeployments() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDeploymentsRead,
		Schema: map[string]*schema.Schema{
			"deployments": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"deployment_state": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"deployment_dseq": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"deployment_owner": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceDeploymentsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	akash := m.(*client.AkashClient)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	deployments, err := akash.GetDeployments()
	if err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("deployments", deployments); err != nil {
		return diag.FromErr(err)
	}

	// always run
	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))

	return diags
}
