package akash

import (
	"context"
	"errors"
	"fmt"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"strings"
	"terraform-provider-akash/akash/client"
	"terraform-provider-akash/akash/client/types"
	"terraform-provider-akash/akash/util"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const IdSeparator = ":"
const DeploymentIdDseq = 0
const DeploymentIdGseq = 1
const DeploymentIdOseq = 2
const DeploymentIdOwner = 3
const DeploymentIdProvider = 4

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
			"deployment_gseq": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"deployment_oseq": &schema.Schema{
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
			"provider_filters": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"provider_preferred": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"enforce": {
							Type:     schema.TypeBool,
							Optional: true,
						},
					},
				},
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

	seqs, err := akash.CreateDeployment(manifestLocation)
	if err != nil {
		return diag.FromErr(err)
	}

	bids, diagnostics := queryBids(ctx, akash, seqs)
	if diagnostics != nil {
		return diagnostics
	}

	var provider string
	bidsProviders := bids.GetProviderAddresses()

	// Handle provider filters
	if f, ok := d.GetOk("provider_filters"); ok {
		tflog.Info(ctx, "Filters provided")

		filters := f.([]interface{})
		filter := filters[0].(map[string]interface{})
		if !ok {
			return diag.FromErr(errors.New("at least one field is expected inside filters"))
		}

		if preferredProvider, ok := filter["provider_preferred"]; ok && util.Contains(bidsProviders, preferredProvider.(string)) {
			tflog.Info(ctx, "Accepting preferred provider's bid")
			provider = preferredProvider.(string)
		} else {
			tflog.Warn(ctx, "Preferred provider did not bid")
			if enforced, ok := filter["enforce"]; ok && enforced.(bool) {
				tflog.Warn(ctx, "Could not find the preferred provider, deleting deployment")
				if err := akash.DeleteDeployment(seqs.Dseq, akash.Config.AccountAddress); err != nil {
					return diag.FromErr(err)
				}
				return diag.FromErr(errors.New("could not find the preferred provider"))
			} else {
				tflog.Warn(ctx, "Not enforcing filters, selecting another provider")
				provider = selectProvider(ctx, akash, bids)
			}
		}
	} else {
		tflog.Info(ctx, "Filters were not provided")
		provider = selectProvider(ctx, akash, bids)
	}

	if diagnostics := createLease(ctx, akash, seqs, provider); diagnostics != nil {
		tflog.Warn(ctx, "Could not create lease, deleting deployment")
		err := akash.DeleteDeployment(seqs.Dseq, akash.Config.AccountAddress)
		if err != nil {
			return diag.FromErr(err)
		}
		return diagnostics
	}

	if diagnostics := sendManifest(ctx, akash, seqs, provider, manifestLocation); diagnostics != nil {
		tflog.Warn(ctx, "Could not send manifest, deleting deployment")
		err := akash.DeleteDeployment(seqs.Dseq, akash.Config.AccountAddress)
		if err != nil {
			return diag.FromErr(err)
		}
		return diagnostics
	}
	tflog.Info(ctx, "Setting created state")
	if diagnostics := setCreatedState(d, akash.Config.AccountAddress, seqs, provider); diagnostics != nil {
		tflog.Warn(ctx, "Could not set state to created, deleting deployment")
		err := akash.DeleteDeployment(seqs.Dseq, akash.Config.AccountAddress)
		if err != nil {
			return diag.FromErr(err)
		}
		return diagnostics
	}

	d.SetId(seqs.Dseq + IdSeparator + seqs.Gseq + IdSeparator + seqs.Oseq + IdSeparator + akash.Config.AccountAddress + IdSeparator + provider)

	return resourceDeploymentRead(ctx, d, m)
}

func queryBids(ctx context.Context, akash *client.AkashClient, seqs client.Seqs) (types.Bids, diag.Diagnostics) {
	tflog.Info(ctx, "Querying available bids")
	bids, err := akash.GetBids(seqs, time.Minute)
	if err != nil {
		return nil, diag.FromErr(err)
	}
	if len(bids) == 0 {
		return nil, diag.FromErr(errors.New("no bids on deployment"))
	}
	tflog.Info(ctx, fmt.Sprintf("Received %d bids", len(bids)))
	return bids, nil
}

func sendManifest(ctx context.Context, akash *client.AkashClient, seqs client.Seqs, provider string, manifestLocation string) diag.Diagnostics {
	tflog.Info(ctx, fmt.Sprintf("Sending manifest %s to %s", manifestLocation, provider))
	res, err := akash.SendManifest(seqs.Dseq, provider, manifestLocation)
	if err != nil {
		tflog.Error(ctx, "Error sending manifest")
		return diag.FromErr(err)
	}
	tflog.Debug(ctx, fmt.Sprintf("Result: %s", res))
	return nil
}

func setCreatedState(d *schema.ResourceData, address string, seqs client.Seqs, provider string) diag.Diagnostics {
	if err := d.Set("deployment_dseq", seqs.Dseq); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("deployment_gseq", seqs.Gseq); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("deployment_oseq", seqs.Oseq); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("deployment_owner", address); err != nil {
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
	provider, err := akash.FindCheapest(bids)
	if err != nil {
		diag.FromErr(err)
		return ""
	}

	tflog.Debug(ctx, fmt.Sprintf("Selected provider %s", provider))
	return provider
}

func createLease(ctx context.Context, akash *client.AkashClient, seqs client.Seqs, provider string) diag.Diagnostics {
	tflog.Info(ctx, "Creating lease")
	lease, err := akash.CreateLease(seqs, provider)
	if err != nil {
		tflog.Error(ctx, "Failed creating the lease")
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

	leaseStatus, err := akash.GetLeaseStatus(client.Seqs{
		Dseq: deploymentId[DeploymentIdDseq],
		Gseq: deploymentId[DeploymentIdGseq],
		Oseq: deploymentId[DeploymentIdOseq],
	}, deploymentId[DeploymentIdProvider])
	if err != nil {
		return diag.FromErr(err)
	}

	services := extractServicesFromLeaseStatus(*leaseStatus)

	if err := d.Set("services", services); err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func extractServicesFromLeaseStatus(leaseStatus types.LeaseStatus) []interface{} {
	var services []interface{}

	for key, value := range leaseStatus.Services {
		service := make(map[string]interface{})
		service["service_name"] = key
		service["service_uri"] = strings.Join(value.URIs, " | ")

		services = append(services, service)
	}
	return services
}

func resourceDeploymentUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	akash := m.(*client.AkashClient)

	deploymentId := strings.Split(d.Id(), IdSeparator)

	seqs := client.Seqs{
		Dseq: deploymentId[DeploymentIdDseq],
		Gseq: deploymentId[DeploymentIdGseq],
		Oseq: deploymentId[DeploymentIdOseq],
	}
	provider := deploymentId[DeploymentIdProvider]

	if d.HasChange("sdl") {
		manifestLocation, err := CreateTemporaryDeploymentFile(ctx, d.Get("sdl").(string))

		// Update the deployment
		if err := akash.UpdateDeployment(seqs.Dseq, manifestLocation); err != nil {
			return diag.FromErr(err)
		}

		if diagnostics := sendManifest(ctx, akash, seqs, provider, manifestLocation); diagnostics != nil {
			return diagnostics
		}

		err = d.Set("last_updated", time.Now().Format(time.RFC850))
		if err != nil {
			return diag.FromErr(err)
		}
	}

	if d.HasChange("provider_filters") {
		tflog.Warn(ctx, "Ignoring filters on resource update")
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
