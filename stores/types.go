package stores

import (
	"context"

	"dev.azure.com/Unum/Mkt_Analytics/_git/domwalk/domains"
)

type DomainStorer interface {
	PutDomains(context.Context, []domains.Domain) error
	GetDomains(context.Context, string) ([]domains.Domain, error)
}
