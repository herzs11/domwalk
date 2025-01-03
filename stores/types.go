package stores

import (
	"context"

	"github.com/herzs11/go-doms/domain"
)

type DomainStorer interface {
	PutDomains(context.Context, []domain.Domain) error
	GetDomains(context.Context, string) ([]domain.Domain, error)
}
