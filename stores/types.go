package stores

import (
	"context"

	"github.com/herzs11/domwalk/domains"
)

type DomainStorer interface {
	PutDomains(context.Context, []domains.Domain) error
	GetDomains(context.Context, string) ([]domains.Domain, error)
}
