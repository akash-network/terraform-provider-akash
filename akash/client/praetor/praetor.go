package praetor

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"terraform-provider-akash/akash/client/types"
)

type Response[T any] struct {
	Status string `json:"status"`
	Data   T      `json:"data"`
}

type ProvidersData struct {
	Count struct {
		Active   uint32 `json:"active"`
		Inactive uint32 `json:"inactive"`
		Total    uint32 `json:"total"`
	} `json:"count"`
	Providers Providers `json:"providers"`
}

type Providers struct {
	ActiveProviders []Provider `json:"active_providers"`
}

type Provider struct {
	Address    string      `json:"address"`
	Uptime     float32     `json:"uptime"`
	Attributes []Attribute `json:"attributes"`
}

type Attribute struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

func GetAllProviders() []types.Provider {
	req, err := http.NewRequest(http.MethodGet, "https://api.praetorapp.com/providers", nil)
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

	var praetorResp Response[ProvidersData]
	dec := json.NewDecoder(resp.Body)
	if err := dec.Decode(&praetorResp); err != nil {
		log.Fatalf("error: can't decode - %s", err)
	}

	providers := make([]types.Provider, 0, len(praetorResp.Data.Providers.ActiveProviders))

	for _, provider := range praetorResp.Data.Providers.ActiveProviders {
		unwrapped := unwrapAttributes(provider.Attributes)

		providers = append(providers, types.Provider{
			Address:    provider.Address,
			Active:     true,
			Uptime:     provider.Uptime,
			Attributes: unwrapped,
		})
	}

	return providers
}

func GetActiveProviders() []types.Provider {
	providers := GetAllProviders()
	activeProviders := make([]types.Provider, 0, len(providers))

	for _, p := range providers {
		if p.Active {
			activeProviders = append(activeProviders, p)
		}
	}

	return activeProviders
}

func unwrapAttributes(attributes []Attribute) map[string]string {
	out := make(map[string]string, len(attributes))
	for _, attr := range attributes {
		out[attr.Key] = attr.Value
	}

	return out
}

func main() {
	providers := GetAllProviders()

	for _, p := range providers {
		fmt.Println(p.Address)
	}
}
