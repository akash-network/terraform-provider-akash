package client

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"net/http"
	"os"
	"terraform-provider-akash/akash/client/cli"
	"terraform-provider-akash/akash/client/types"
	"time"
)

func (ak *AkashClient) GetDeployments() ([]map[string]interface{}, error) {
	client := &http.Client{Timeout: 10 * time.Second}
	address := os.Getenv("AKASH_ACCOUNT_ADDRESS")

	req, err := http.NewRequest("GET", fmt.Sprintf("%s/akash/deployment/v1beta2/deployments/list?filters.owner=%s", "http://135.181.181.122:1518", address), nil)
	if err != nil {
		return nil, err
	}

	r, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()

	parsed := types.DeploymentResponse{}

	err = json.NewDecoder(r.Body).Decode(&parsed)
	if err != nil {
		return nil, err
	}

	parsedDeployments := parsed.Deployments
	deployments := make([]map[string]interface{}, 0)

	for _, deployment := range parsedDeployments {
		d := make(map[string]interface{})
		d["deployment_state"] = deployment.DeploymentInfo.State
		d["deployment_dseq"] = deployment.DeploymentInfo.DeploymentId.Dseq
		d["deployment_owner"] = deployment.DeploymentInfo.DeploymentId.Owner

		deployments = append(deployments, d)
	}

	return deployments, nil
}

func (ak *AkashClient) GetDeployment(dseq string, owner string) (map[string]interface{}, error) {
	client := &http.Client{Timeout: 10 * time.Second}

	req, err := http.NewRequest("GET", fmt.Sprintf("%s/akash/deployment/v1beta2/deployments/info?id.owner=%s&id.dseq=%s", "http://135.181.181.122:1518", owner, dseq), nil)
	if err != nil {
		return nil, err
	}

	r, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()

	deployment := types.Deployment{}

	err = json.NewDecoder(r.Body).Decode(&deployment)
	if err != nil {
		return nil, err
	}

	d := make(map[string]interface{})
	d["deployment_state"] = deployment.DeploymentInfo.State
	d["deployment_dseq"] = deployment.DeploymentInfo.DeploymentId.Dseq
	d["deployment_owner"] = deployment.DeploymentInfo.DeploymentId.Owner
	d["escrow_account_owner"] = deployment.EscrowAccount.Owner
	d["escrow_account_state"] = deployment.EscrowAccount.State
	d["escrow_account_balance_amount"] = deployment.EscrowAccount.Balance.Amount
	d["escrow_account_balance_denom"] = deployment.EscrowAccount.Balance.Denom

	return d, nil
}

func (ak *AkashClient) CreateDeployment(manifestLocation string) (string, error) {

	tflog.Debug(ak.ctx, "Creating deployment")
	// Create deployment using the file created with the SDL
	dseq, err := transactionCreateDeployment(ak.ctx, manifestLocation)
	if err != nil {
		tflog.Error(ak.ctx, "Failed creating deployment")
		tflog.Debug(ak.ctx, fmt.Sprintf("%s", err))
		return "", err
	}
	tflog.Info(ak.ctx, "Deployment created with DSEQ "+dseq)

	return dseq, nil
}

// Perform the transaction to create the deployment and return either the DSEQ or an error.
func transactionCreateDeployment(ctx context.Context, manifestLocation string) (string, error) {
	cmd := cli.AkashCli(ctx).Tx().Deployment().Create().Manifest(manifestLocation).
		SetFees(5000).AutoAccept().SetFrom(os.Getenv("AKASH_KEY_NAME")).OutputJson()

	transaction := types.Transaction{}
	if err := cmd.DecodeJson(&transaction); err != nil {
		return "", err
	}

	if len(transaction.Logs) == 0 {
		return "", errors.New(fmt.Sprintf("something went wrong: %s", transaction.RawLog))
	}

	return transaction.Logs[0].Events[0].Attributes.Get("dseq")
}

func (ak *AkashClient) DeleteDeployment(ctx context.Context, dseq string, owner string) error {
	cmd := cli.AkashCli(ctx).Tx().Deployment().Close().
		SetDseq(dseq).SetOwner(owner).SetFrom(os.Getenv("AKASH_KEY_NAME")).
		DefaultGas().AutoAccept().OutputJson()

	out, err := cmd.Raw()
	if err != nil {
		return err
	}

	tflog.Debug(ctx, fmt.Sprintf("Response: %s", out))

	return nil
}

func (ak *AkashClient) UpdateDeployment(dseq string, manifestLocation string) error {
	cmd := cli.AkashCli(ak.ctx).Tx().Deployment().Update().Manifest(manifestLocation).
		SetDseq(dseq).SetFrom(os.Getenv("AKASH_KEY_NAME")).DefaultGas().AutoAccept().OutputJson()

	out, err := cmd.Raw()
	if err != nil {
		return err
	}

	tflog.Debug(ak.ctx, fmt.Sprintf("Response: %s", out))

	return nil
}
