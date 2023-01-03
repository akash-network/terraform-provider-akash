package providers_api

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"path"
	"terraform-provider-akash/akash/client/types"
)

type provider struct {
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

// New creates a new ProviderClient based on the given host.
func New(host string) *ProvidersClient {
	return &ProvidersClient{
		host: host,
	}
}

// GetAllProviders gets all the providers from the providers' API. Returns error in case something goes wrong.
func (c *ProvidersClient) GetAllProviders() ([]types.Provider, error) {
	addr := path.Join(c.host, "provider/")
	req, err := http.NewRequest(http.MethodGet, addr, nil)
	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(fmt.Sprintf("Response status code %d", resp.StatusCode))
	}

	var result []provider
	dec := json.NewDecoder(resp.Body)
	if err := dec.Decode(&result); err != nil {
		return nil, err
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

	return providers, err
}

// GetActiveProviders gets the active providers from the providers' API.
func (c *ProvidersClient) GetActiveProviders() ([]types.Provider, error) {
	providers, err := c.GetAllProviders()
	if err != nil {
		return nil, err
	}
	activeProviders := make([]types.Provider, 0, len(providers))

	for _, p := range providers {
		if p.Active {
			activeProviders = append(activeProviders, p)
		}
	}

	return activeProviders, nil
}
