package types

type LeaseStatus struct {
	Services map[string]Service
}

type Service struct {
	URIs      []string `json:"uris"`
	Name      string   `json:"name"`
	Available int32    `json:"available"`
}
