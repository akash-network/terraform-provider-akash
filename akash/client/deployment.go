package client

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"net/http"
	"terraform-provider-akash/akash/client/cli"
	"terraform-provider-akash/akash/client/types"
	"time"
)

type Seqs struct {
	Dseq string
	Gseq string
	Oseq string
}

func (ak *AkashClient) GetDeployments() ([]map[string]interface{}, error) {
	client := &http.Client{Timeout: 10 * time.Second}
	address := ak.Config.AccountAddress

	req, err := http.NewRequest("GET", fmt.Sprintf("%s/akash/deployment/v1beta2/deployments/list?filters.owner=%s", "https://akash.c29r3.xyz/api", address), nil)
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

	req, err := http.NewRequest("GET", fmt.Sprintf("%s/akash/deployment/v1beta2/deployments/info?id.owner=%s&id.dseq=%s", "https://akash.c29r3.xyz/api", owner, dseq), nil)
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

func (ak *AkashClient) CreateDeployment(manifestLocation string) (Seqs, error) {

	tflog.Info(ak.ctx, "Creating deployment")
	// Create deployment using the file created with the SDL
	attributes, err := transactionCreateDeployment(ak, manifestLocation)
	if err != nil {
		tflog.Error(ak.ctx, "Failed creating deployment")
		tflog.Debug(ak.ctx, fmt.Sprintf("%s", err))
		return Seqs{}, err
	}

	dseq, _ := attributes.Get("dseq")
	gseq, _ := attributes.Get("gseq")
	oseq, _ := attributes.Get("oseq")

	tflog.Info(ak.ctx, fmt.Sprintf("Deployment created with DSEQ=%s GSEQ=%s OSEQ=%s", dseq, gseq, oseq))

	return Seqs{dseq, gseq, oseq}, nil
}

// Perform the transaction to create the deployment and return either the DSEQ or an error.
func transactionCreateDeployment(ak *AkashClient, manifestLocation string) (types.TransactionEventAttributes, error) {
	cmd := cli.AkashCli(ak).Tx().Deployment().Create().Manifest(manifestLocation).
		DefaultGas().AutoAccept().SetFrom(ak.Config.KeyName).SetKeyringBackend(ak.Config.KeyringBackend).
		SetChainId(ak.Config.ChainId).SetNode(ak.Config.Node).OutputJson()

	transaction := types.Transaction{}
	if err := cmd.DecodeJson(&transaction); err != nil {
		return nil, err
	}

	if len(transaction.Logs) == 0 {
		return nil, errors.New(fmt.Sprintf("something went wrong: %s", transaction.RawLog))
	}

	return transaction.Logs[0].Events[0].Attributes, nil
}

func (ak *AkashClient) DeleteDeployment(dseq string, owner string) error {
	cmd := cli.AkashCli(ak).Tx().Deployment().Close().
		SetDseq(dseq).SetOwner(owner).SetFrom(ak.Config.KeyName).
		DefaultGas().SetChainId(ak.Config.ChainId).SetKeyringBackend(ak.Config.KeyringBackend).
		SetNode(ak.Config.Node).AutoAccept().OutputJson()

	out, err := cmd.Raw()
	if err != nil {
		return err
	}

	tflog.Debug(ak.ctx, fmt.Sprintf("Response: %s", out))

	return nil
}

func (ak *AkashClient) UpdateDeployment(dseq string, manifestLocation string) error {
	cmd := cli.AkashCli(ak).Tx().Deployment().Update().Manifest(manifestLocation).
		SetDseq(dseq).SetFrom(ak.Config.KeyName).SetNode(ak.Config.Node).
		SetKeyringBackend(ak.Config.KeyringBackend).SetChainId(ak.Config.ChainId).
		GasAuto().SetGasAdjustment(1.5).SetGasPrices().SetSignMode("amino-json").AutoAccept().OutputJson()

	out, err := cmd.Raw()
	if err != nil {
		return err
	}

	tflog.Debug(ak.ctx, fmt.Sprintf("Response: %s", out))

	return nil
}
