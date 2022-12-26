package types

type Provider struct {
	Address    string            `json:"address"`
	Active     bool              `json:"active"`
	Uptime     float32           `json:"uptime"`
	Attributes map[string]string `json:"attributes"`
}
