package providers_api

import (
	"encoding/json"
	"log"
	"net/http"
	"terraform-provider-akash/akash/client/types"
)

type Provider struct {
	Address    string            `json:"address"`
	Active     bool              `json:"active"`
	Uptime     uptime            `json:"uptime"`
	Attributes map[string]string `json:"extraAttributes"`
}

type uptime struct {
	Percentage float32 `json:"percentage"`
	Since      string  `json:"since"`
}

type ProvidersClient struct {
	host string
}

func New(host string) *ProvidersClient {
	return &ProvidersClient{
		host: host,
	}
}

func (c *ProvidersClient) GetAllProviders() []types.Provider {
	req, err := http.NewRequest(http.MethodGet, c.host, nil)
	if err != nil {
		log.Fatalf("error: %s", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatalf("error: %s", err)
	}

	if resp.StatusCode != http.StatusOK {
		log.Fatalf("error: %s", resp.Status)
	}

	var result []Provider
	dec := json.NewDecoder(resp.Body)
	if err := dec.Decode(&result); err != nil {
		log.Fatalf("error: can't decode - %s", err)
	}

	providers := make([]types.Provider, 0, len(result))

	for _, provider := range result {

		// TODO: Fix bad design
		providers = append(providers, types.Provider{
			Address:    provider.Address,
			Active:     provider.Active,
			Uptime:     provider.Uptime.Percentage,
			Attributes: provider.Attributes,
		})
	}

	return providers
}

func (c *ProvidersClient) GetActiveProviders() []types.Provider {
	providers := c.GetAllProviders()
	activeProviders := make([]types.Provider, 0, len(providers))

	for _, p := range providers {
		if p.Active {
			activeProviders = append(activeProviders, p)
		}
	}

	return activeProviders
}
