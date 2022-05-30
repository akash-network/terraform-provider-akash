package types

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
