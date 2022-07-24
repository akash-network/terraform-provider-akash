package akash

import (
	"context"
	"errors"
	"fmt"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"os"
	"strings"
	"terraform-provider-akash/akash/client"
	"terraform-provider-akash/akash/client/types"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const IdSeparator = ":"
const DeploymentIdDseq = 0
const DeploymentIdOwner = 1
const DeploymentIdProvider = 2

func resourceDeployment() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDeploymentCreate,
		ReadContext:   resourceDeploymentRead,
		UpdateContext: resourceDeploymentUpdate,
		DeleteContext: resourceDeploymentDelete,
		Schema: map[string]*schema.Schema{
			"sdl": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"last_updated": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
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
			"escrow_account_owner": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"escrow_account_balance_denom": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"escrow_account_balance_amount": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"provider_address": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"services": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"service_name": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"service_uri": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceDeploymentCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	akash := m.(*client.AkashClient)

	manifestLocation, err := CreateTemporaryDeploymentFile(ctx, d.Get("sdl").(string))

	if err != nil {
		return diag.FromErr(err)
	}

	dseq, err := akash.CreateDeployment(manifestLocation)
	if err != nil {
		return diag.FromErr(err)
	}

	bids, diagnostics := queryBids(ctx, akash, dseq)
	if diagnostics != nil {
		return diagnostics
	}

	provider := selectProvider(ctx, akash, bids)

	if diagnostics := createLease(ctx, akash, dseq, provider); diagnostics != nil {
		err := akash.DeleteDeployment(dseq, os.Getenv("AKASH_ACCOUNT_ADDRESS"))
		if err != nil {
			return diag.FromErr(err)
		}
		return diagnostics
	}
	if diagnostics := sendManifest(ctx, akash, dseq, provider, manifestLocation); diagnostics != nil {
		err := akash.DeleteDeployment(dseq, os.Getenv("AKASH_ACCOUNT_ADDRESS"))
		if err != nil {
			return diag.FromErr(err)
		}
		return diagnostics
	}
	if diagnostics := setCreatedState(d, dseq, provider); diagnostics != nil {
		err := akash.DeleteDeployment(dseq, os.Getenv("AKASH_ACCOUNT_ADDRESS"))
		if err != nil {
			return diag.FromErr(err)
		}
		return diagnostics
	}

	d.SetId(dseq + IdSeparator + os.Getenv("AKASH_ACCOUNT_ADDRESS") + IdSeparator + provider)

	return resourceDeploymentRead(ctx, d, m)
}

func queryBids(ctx context.Context, akash *client.AkashClient, dseq string) (types.Bids, diag.Diagnostics) {
	tflog.Debug(ctx, "Querying available bids")
	bids, err := akash.GetBids(dseq, time.Minute)
	if err != nil {
		return nil, diag.FromErr(err)
	}
	if len(bids) == 0 {
		return nil, diag.FromErr(errors.New("no bids on deployment"))
	}
	tflog.Info(ctx, fmt.Sprintf("Received %d bids in the deployment", len(bids)))
	return bids, nil
}

func sendManifest(ctx context.Context, akash *client.AkashClient, dseq string, provider string, manifestLocation string) diag.Diagnostics {
	tflog.Info(ctx, "Sending the manifest")
	// Send the manifest
	res, err := akash.SendManifest(dseq, provider, manifestLocation)
	if err != nil {
		return diag.FromErr(err)
	}
	tflog.Debug(ctx, fmt.Sprintf("Result: %s", res))
	return nil
}

func setCreatedState(d *schema.ResourceData, dseq string, provider string) diag.Diagnostics {
	if err := d.Set("deployment_dseq", dseq); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("deployment_owner", os.Getenv("AKASH_ACCOUNT_ADDRESS")); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("deployment_state", "active"); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("provider_address", provider); err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func selectProvider(ctx context.Context, akash *client.AkashClient, bids types.Bids) string {
	// Select the provider
	provider, err := akash.FindCheapest(bids)
	if err != nil {
		diag.FromErr(err)
		return ""
	}

	tflog.Debug(ctx, fmt.Sprintf("Selected provider %s", provider))
	return provider
}

func createLease(ctx context.Context, akash *client.AkashClient, dseq string, provider string) diag.Diagnostics {
	tflog.Info(ctx, "Creating lease")
	// Create a lease
	lease, err := akash.CreateLease(dseq, provider)
	if err != nil {
		return diag.FromErr(err)
	}
	tflog.Debug(ctx, fmt.Sprintf("Lease return: %s", lease))
	return nil
}

func resourceDeploymentRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	akash := m.(*client.AkashClient)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	deploymentId := strings.Split(d.Id(), IdSeparator)

	deployment, err := akash.GetDeployment(deploymentId[DeploymentIdDseq], deploymentId[DeploymentIdOwner])
	if err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("deployment_dseq", deployment["deployment_dseq"]); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("deployment_owner", deployment["deployment_owner"]); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("deployment_state", deployment["deployment_state"]); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("escrow_account_owner", deployment["escrow_account_owner"]); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("escrow_account_balance_denom", deployment["escrow_account_balance_denom"]); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("escrow_account_balance_amount", deployment["escrow_account_balance_amount"]); err != nil {
		return diag.FromErr(err)
	}

	leaseStatus, err := akash.GetLeaseStatus(deploymentId[DeploymentIdDseq], deploymentId[DeploymentIdProvider])
	if err != nil {
		return diag.FromErr(err)
	}

	var services []interface{}

	for key, value := range leaseStatus.Services {
		service := make(map[string]interface{})
		service["service_name"] = key
		service["service_uri"] = strings.Join(value.URIs, " | ")

		services = append(services, service)
	}

	if err := d.Set("services", services); err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func resourceDeploymentUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	akash := m.(*client.AkashClient)

	deploymentId := strings.Split(d.Id(), IdSeparator)

	dseq := deploymentId[DeploymentIdDseq]
	provider := deploymentId[DeploymentIdProvider]

	if d.HasChange("sdl") {
		manifestLocation, err := CreateTemporaryDeploymentFile(ctx, d.Get("sdl").(string))

		// Update the deployment
		if err := akash.UpdateDeployment(dseq, manifestLocation); err != nil {
			return diag.FromErr(err)
		}

		if diagnostics := sendManifest(ctx, akash, dseq, provider, manifestLocation); diagnostics != nil {
			return diagnostics
		}

		err = d.Set("last_updated", time.Now().Format(time.RFC850))
		if err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceDeploymentRead(ctx, d, m)
}

func resourceDeploymentDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	akash := m.(*client.AkashClient)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	deploymentId := strings.Split(d.Id(), IdSeparator)

	err := akash.DeleteDeployment(deploymentId[DeploymentIdDseq], deploymentId[DeploymentIdOwner])
	if err != nil {
		return diag.FromErr(err)
	}

	// d.SetId("") is automatically called assuming delete returns no errors, but
	// it is added here for explicitness.
	d.SetId("")

	return diags
}
