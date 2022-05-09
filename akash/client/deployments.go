package client

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type DeploymentId struct {
	Dseq  string `json:"dseq"`
	Owner string `json:"owner"`
}

type DeploymentInfo struct {
	State        string       `json:"state"`
	DeploymentId DeploymentId `json:"deployment_id"`
}

type Deployment struct {
	DeploymentInfo DeploymentInfo `json:"deployment"`
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
	d := make(map[string]interface{})
	d["deployment_state"] = "active"
	d["deployment_dseq"] = "12345"
	d["deployment_owner"] = "akashdokfmdjmf023n32423"

	return d, nil
}

func CreateDeployment(sdl string) (map[string]interface{}, error) {
	err := ioutil.WriteFile("deployment.yaml", []byte(sdl), 0666)
	if err != nil {
		return nil, err
	}

	// Create deployment using the file created with the SDL
	/*	out, err := exec.Command("akash tx deployment create deployment.yaml -o json").Output(); if err != nil {
			return nil, err
		}

		err = json.NewDecoder(strings.NewReader(string(out))).Decode(&out); if err != nil {
			return nil, err
		}*/

	d := make(map[string]interface{})
	d["deployment_state"] = "active"
	d["deployment_dseq"] = "12345"
	d["deployment_owner"] = "akashdokfmdjmf023n32423"

	return d, nil
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
