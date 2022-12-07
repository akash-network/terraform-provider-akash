package akash

import (
	"context"
	"errors"
	"fmt"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"strings"
	"terraform-provider-akash/akash/client"
	"terraform-provider-akash/akash/client/types"
	"terraform-provider-akash/akash/extensions"
	"terraform-provider-akash/akash/filtering"
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
						"providers": {
							Type:     schema.TypeList,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
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
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"uris": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"available":          {Type: schema.TypeInt, Computed: true},
						"total":              {Type: schema.TypeInt, Computed: true},
						"replicas":           {Type: schema.TypeInt, Computed: true},
						"updated_replicas":   {Type: schema.TypeInt, Computed: true},
						"available_replicas": {Type: schema.TypeInt, Computed: true},
						"ready_replicas":     {Type: schema.TypeInt, Computed: true},
						"ips": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"port":          {Type: schema.TypeInt, Computed: true},
									"ip":            {Type: schema.TypeString, Computed: true},
									"external_port": {Type: schema.TypeInt, Computed: true},
									"protocol":      {Type: schema.TypeString, Computed: true},
								},
							},
						},
					},
				},
			},
			"forwarded_ports": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{},
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

	var diags diag.Diagnostics

	manifestLocation, err := CreateTemporaryDeploymentFile(ctx, d.Get("sdl").(string))

	if err != nil {
		diags = append(diags)
		return diag.FromErr(err)
	}

	seqs, err := akash.CreateDeployment(manifestLocation)
	if err != nil {
		return diag.FromErr(err)
	}

	bids, diagnostics := queryBids(ctx, akash, seqs)
	if diagnostics != nil {
		tflog.Warn(ctx, "No bids on deployment")
		if err := akash.DeleteDeployment(seqs.Dseq, akash.Config.AccountAddress); err != nil {
			return append(diagnostics, diag.FromErr(err)...)
		}
		return diagnostics
	}

	provider, err := selectProvider(ctx, d, bids)
	if err != nil {
		if err := akash.DeleteDeployment(seqs.Dseq, akash.Config.AccountAddress); err != nil {
			return diag.FromErr(err)
		}

		return diag.FromErr(err)
	}

	if diagnostics := createLease(ctx, akash, seqs, provider); diagnostics != nil {
		tflog.Warn(ctx, "Could not create lease, deleting deployment")
		err := akash.DeleteDeployment(seqs.Dseq, akash.Config.AccountAddress)
		if err != nil {
			// TODO: Add diagnostic warning saying deployment was not deleted
			return append(diagnostics, diag.FromErr(err)...)
		}
		return diagnostics
	}

	if diagnostics := sendManifest(ctx, akash, seqs, provider, manifestLocation); diagnostics != nil {
		tflog.Warn(ctx, "Could not send manifest, deleting deployment")
		err := akash.DeleteDeployment(seqs.Dseq, akash.Config.AccountAddress)
		if err != nil {
			return append(diagnostics, diag.FromErr(err)...)
		}
		return diagnostics
	}
	tflog.Info(ctx, "Setting created state")
	if diagnostics := setCreatedState(d, akash.Config.AccountAddress, seqs, provider); diagnostics != nil {
		tflog.Warn(ctx, "Could not set state to created, deleting deployment")
		err := akash.DeleteDeployment(seqs.Dseq, akash.Config.AccountAddress)
		if err != nil {
			return append(diagnostics, diag.FromErr(err)...)
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
	_, err := akash.SendManifest(seqs.Dseq, provider, manifestLocation)
	if err != nil {
		tflog.Error(ctx, "Error sending manifest")
		return diag.FromErr(err)
	}
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

// This function gets all the configured filters and applies them. In the end it selects the cheapest provider.
func selectProvider(ctx context.Context, d *schema.ResourceData, bids types.Bids) (string, error) {
	filterPipeline := filtering.NewFilterPipeline(bids)

	if f, ok := d.GetOk("provider_filters"); ok {
		tflog.Info(ctx, "Filters provided")

		filters := f.([]interface{})
		filter := filters[0].(map[string]interface{})
		if !ok {
			return "", errors.New("at least one field is expected inside filters")
		}

		bidsProviders := bids.GetProviderAddresses()

		uncastProviders, ok := filter["providers"].([]interface{})
		if !ok {
			tflog.Debug(ctx, fmt.Sprintf("Could not convert: %+v\n", filter["providers"]))
			return "", errors.New("could not get 'providers' filter")
		}

		preferredProviders := make([]string, len(uncastProviders))
		for _, uncastProvider := range uncastProviders {
			preferredProviders = append(preferredProviders, uncastProvider.(string))
		}

		if extensions.ContainsAny(bidsProviders, preferredProviders) {
			tflog.Info(ctx, "Accepting preferred provider's bid")
			// Add pipe to get bids of preferred providers
			filterPipeline.Pipe(func(bids types.Bids) (types.Bids, error) {
				return bids.FindAllByProviders(preferredProviders), nil
			})
		} else {
			tflog.Warn(ctx, "Preferred provider did not bid")
			if enforced, ok := filter["enforce"]; ok && enforced.(bool) {
				tflog.Warn(ctx, "Could not find the preferred provider, deleting deployment")
				return "", errors.New("preferred providers did not bid")
			} else {
				tflog.Warn(ctx, "Not enforcing filters, selecting another provider")
			}
		}
	} else {
		tflog.Info(ctx, "Filters were not provided")
	}

	bid, err := filterPipeline.Reduce(filtering.Cheapest)
	if err != nil {
		return "", err
	}

	tflog.Info(ctx, fmt.Sprintf("Selected %s for %fuakt", bid.Id.Provider, bid.Price.Amount))
	return bid.Id.Provider, nil
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

	if err := d.Set("deployment_dseq", deployment.DeploymentInfo.DeploymentId.Dseq); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("deployment_owner", deployment.DeploymentInfo.DeploymentId.Owner); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("deployment_state", deployment.DeploymentInfo.State); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("escrow_account_owner", deployment.EscrowAccount.Owner); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("escrow_account_balance_denom", deployment.EscrowAccount.Balance.Denom); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("escrow_account_balance_amount", deployment.EscrowAccount.Balance.Amount); err != nil {
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

	tflog.Info(ctx, fmt.Sprintf("Extracted %d services from lease-status", len(services)))
	tflog.Debug(ctx, fmt.Sprintf("Services: %+v", services))

	if err := d.Set("services", services); err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func extractServicesFromLeaseStatus(leaseStatus types.LeaseStatus) []map[string]interface{} {
	// TODO: Force ordering of services to provide some predictability to the position of the services
	services := make([]map[string]interface{}, 0)

	for key, value := range leaseStatus.Services {
		service := make(map[string]interface{})
		service["name"] = key
		service["uris"] = value.URIs
		service["replicas"] = value.Replicas
		service["updated_replicas"] = value.UpdatedReplicas
		service["available_replicas"] = value.AvailableReplicas
		service["ready_replicas"] = value.ReadyReplicas
		service["available"] = value.Available
		service["total"] = value.Total

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
