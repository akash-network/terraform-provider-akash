package client

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"terraform-provider-hashicups/akash/client/types"
	"time"
)

const AKASH_BINARY = "../bin/akash"

func GetDeployments() ([]map[string]interface{}, error) {
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

func GetDeployment(dseq string, owner string) (map[string]interface{}, error) {
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

// CreateDeployment TODO: Extract the different operations inside here to different methods on different clients.
// Transactions, Deployments...
func CreateDeployment(ctx context.Context, sdl string) (string, error) {
	tflog.Debug(ctx, "Creating temporary deployment file")

	err := ioutil.WriteFile("deployment.yaml", []byte(sdl), 0666)
	if err != nil {
		return "", err
	}

	tflog.Debug(ctx, "Creating deployment")
	// Create deployment using the file created with the SDL
	dseq, err := transactionCreateDeployment()
	if err != nil {
		tflog.Error(ctx, "Failed creating deployment")
		tflog.Debug(ctx, fmt.Sprintf("%s", err))
		return "", err
	}
	tflog.Info(ctx, "Deployment created with DSEQ "+dseq)

	return dseq, nil
}

// Perform the transaction to create the deployment and return either the DSEQ or an error.
func transactionCreateDeployment() (string, error) {
	cmd := exec.Command(
		AKASH_BINARY,
		"tx",
		"deployment",
		"create",
		"deployment.yaml",
		"--fees",
		"5000uakt",
		"-y",
		"--from",
		os.Getenv("AKASH_KEY_NAME"),
		"-o json",
	)

	var errb bytes.Buffer
	cmd.Stderr = &errb
	out, err := cmd.Output()
	if err != nil {
		return "", errors.New(errb.String())
	}

	transaction := types.Transaction{}
	err = json.NewDecoder(strings.NewReader(string(out))).Decode(&transaction)
	if err != nil {
		return "", err
	}

	return transaction.Logs[0].Events[0].Attributes.Get("dseq")
}

func DeleteDeployment(dseq string, owner string) error {
	/*	out, err := exec.Command("akash tx deployment close --dseq " + dseq + " --owner " + owner + " --from pktminerwallet -y --fees 5000uakt").Output()
		if err != nil {
			return err
		}

		err = json.NewDecoder(strings.NewReader(string(out))).Decode(&out)
		if err != nil {
			return err
		}*/

	return nil
}
