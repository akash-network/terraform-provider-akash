package client

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os/exec"
	"strings"
	"terraform-provider-hashicups/akash/client/types"
	"time"
)

const AKASH_BINARY = "./bin/akash"

type DeploymentId struct {
	Dseq  string `json:"dseq"`
	Owner string `json:"owner"`
}

type DeploymentInfo struct {
	State        string       `json:"state"`
	DeploymentId DeploymentId `json:"deployment_id"`
}

type EscrowAccountBalance struct {
	Denom  string `json:"denom"`
	Amount string `json:"amount"`
}

type EscrowAccount struct {
	Owner   string               `json:"owner"`
	State   string               `json:"state"`
	Balance EscrowAccountBalance `json:"balance"`
}

type Deployment struct {
	DeploymentInfo DeploymentInfo `json:"deployment"`
	EscrowAccount  EscrowAccount  `json:"escrow_account"`
}

type DeploymentResponse struct {
	Deployments []Deployment `json:"deployments"`
}

func GetDeployments() ([]map[string]interface{}, error) {
	client := &http.Client{Timeout: 10 * time.Second}
	address := "akash1qyfg4zl2dku8ry7gjkhf88vnc3zrn6vmnzlvr9"

	req, err := http.NewRequest("GET", fmt.Sprintf("%s/akash/deployment/v1beta2/deployments/list?filters.owner=%s", "http://135.181.181.122:1518", address), nil)
	if err != nil {
		return nil, err
	}

	r, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()

	parsed := DeploymentResponse{}

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

	deployment := Deployment{}

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

func CreateDeployment(sdl string) (map[string]interface{}, error) {
	err := ioutil.WriteFile("deployment.yaml", []byte(sdl), 0666)
	if err != nil {
		return nil, err
	}

	// Create deployment using the file created with the SDL
	dseq, err := transactionCreateDeployment(err)
	if err != nil {
		return nil, err
	}

	// Check bids on deployments and choose one
	bids, err := queryBidList(dseq)
	if err != nil {
		return nil, err
	}
	provider := bids[0].Id.Provider

	// Create lease and send manifest

	d := make(map[string]interface{})
	d["deployment_state"] = "active"
	d["deployment_dseq"] = dseq
	d["deployment_owner"] = "akashdokfmdjmf023n32423"

	return d, nil
}

// Perform the transaction to create the deployment and return either the DSEQ or an error.
func transactionCreateDeployment(err error) (string, error) {
	out, err := exec.Command(AKASH_BINARY + " tx deployment create deployment.yaml --fees 5000uakt -y --from $AKASH_KEY_NAME -o json").Output()
	if err != nil {
		return "", err
	}

	transaction := types.Transaction{}
	err = json.NewDecoder(strings.NewReader(string(out))).Decode(&transaction)
	if err != nil {
		return "", err
	}

	return transaction.Logs[0].Events[0].Attributes.Get("dseq")
}

func queryBidList(dseq string) (types.Bids, error) {
	out, err := exec.Command(fmt.Sprintf("%s query bid list --dseq %s -o json", AKASH_BINARY, dseq)).Output()
	if err != nil {
		return nil, err
	}

	bids := types.Bids{}
	err = json.NewDecoder(strings.NewReader(string(out))).Decode(&bids)
	if err != nil {
		return nil, err
	}

	return bids, nil
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
