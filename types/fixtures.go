package types

import (
	"time"

	"domwalk/db"
	"gorm.io/gorm"
)

type MatchedDomain struct {
	gorm.Model
	DomainID        uint   `json:"reqDomainName,omitempty" bigquery:"domain_name" gorm:"index:,unique,composite:dom_match"`
	MatchedDomainID uint   `json:"certSanDomainName,omitempty" bigquery:"matched_domain_name" gorm:"index:,unique,composite:dom_match"`
	Domain          Domain `json:"certSanDomain,omitempty" gorm:"foreignKey:MatchedDomainID" bigquery:"-"`
}

func (m *MatchedDomain) ToBQ() MatchedDomainBQ {
	return MatchedDomainBQ{
		ID:              int(m.ID),
		CreatedAt:       m.CreatedAt,
		UpdatedAt:       m.UpdatedAt,
		DomainID:        int(m.DomainID),
		MatchedDomainID: int(m.MatchedDomainID),
	}
}

type MatchedDomainBQ struct {
	ID              int       `bigquery:"id"`
	CreatedAt       time.Time `bigquery:"created_at"`
	UpdatedAt       time.Time `bigquery:"updated_at"`
	DomainID        int       `bigquery:"domain_id"`
	MatchedDomainID int       `bigquery:"matched_domain_id"`
}

func ClearTables() {
	err := db.GormDB.Migrator().DropTable(
		&Domain{}, &MXRecord{}, &ARecord{}, &AAAARecord{}, &SOARecord{},
		&WebRedirectDomain{}, &Sitemap{}, &CertSansDomain{}, &SitemapWebDomain{}, &SitemapContactDomain{},
	)
	if err != nil {
		panic(err)
	}
}

func CreateTables() {
	err := db.GormDB.Migrator().AutoMigrate(
		&Domain{}, &MXRecord{}, &ARecord{}, &AAAARecord{}, &SOARecord{},
		&WebRedirectDomain{}, &Sitemap{}, &CertSansDomain{}, &SitemapWebDomain{}, &SitemapContactDomain{},
	)
	if err != nil {
		panic(err)
	}
}
