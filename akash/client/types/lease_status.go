package types

import "net"

type LeaseStatus struct {
	Services       map[string]Service `json:"services"`
	ForwardedPorts []ForwardedPort    `json:"forwarded_ports"`
	IPs            map[string]IP      `json:"ips"`
}

type IP struct {
	Address net.IP `json:"address"`
}

type ForwardedPort struct {
}

type Service struct {
	URIs              []string `json:"uris"`
	Name              string   `json:"name"`
	Available         int32    `json:"available"`
	Total             int32    `json:"total"`
	Replicas          int32    `json:"replicas"`
	UpdatedReplicas   int32    `json:"updated_replicas"`
	ReadyReplicas     int32    `json:"ready_replicas"`
	AvailableReplicas int32    `json:"available_replicas"`
}
