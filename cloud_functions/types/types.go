package types

import (
	"dev.azure.com/Unum/Mkt_Analytics/_git/domwalk/domains"
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
