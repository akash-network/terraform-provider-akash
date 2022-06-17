package types

type LeaseStatus struct {
	Services map[string]Service
}

type Service struct {
	Name      string   `json:"name"`
	Available int32    `json:"available"`
	URIs      []string `json:"uris"`
}
