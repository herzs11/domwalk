package domains

import (
	"time"
)

type MatchedDomain struct {
	CreatedAt  time.Time `json:"createdAt,omitempty"`
	UpdatedAt  time.Time `json:"updatedAt,omitempty"`
	DomainName string    `json:"matchedDomain,omitempty"`
}
