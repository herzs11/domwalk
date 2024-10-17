package types

import (
	"domwalk/domains"
)

type RequestParams struct {
	DomainNames []string `json:"domain_names"`
	ProcessConfig
	OnlyMatchedDomains bool `json:"only_matched_domains,omitempty"`
}

type ProcessConfig struct {
	Workers int `json:"workers,omitempty"`
	domains.EnrichmentConfig
}
