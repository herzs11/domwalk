package types

import (
	"github.com/herzs11/domwalk/domains"
)

type RequestParams struct {
	DomainNames []string `json:"domain_names"`
	ProcessConfig
	NoResponse         bool `json:"no_response,omitempty"`
	OnlyMatchedDomains bool `json:"only_matched_domains,omitempty"`
}

type ProcessConfig struct {
	Workers int `json:"workers,omitempty"`
	domains.EnrichmentConfig
}
