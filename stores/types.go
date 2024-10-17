package stores

import (
	"context"

	"domwalk/domains"
)

type DomainStorer interface {
	PutDomains(context.Context, []domains.Domain) error
	GetDomains(context.Context, string) ([]domains.Domain, error)
}
