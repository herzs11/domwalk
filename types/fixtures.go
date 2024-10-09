package types

import (
	"time"

	"domwalk/db"
	"gorm.io/gorm"
)

type MatchedDomain struct {
	gorm.Model
	DomainName        string `json:"reqDomainName,omitempty" bigquery:"domain_name" gorm:"index:,unique,composite:dom_match"`
	MatchedDomainName string `json:"certSanDomainName,omitempty" bigquery:"matched_domain_name" gorm:"index:,unique,composite:dom_match"`
	Domain            Domain `json:"certSanDomain,omitempty" gorm:"foreignKey:MatchedDomainName" bigquery:"-"`
}

func (m *MatchedDomain) ToBQ() MatchedDomainBQ {
	return MatchedDomainBQ{
		ID:                int(m.ID),
		CreatedAt:         m.CreatedAt,
		UpdatedAt:         m.UpdatedAt,
		DomainName:        m.DomainName,
		MatchedDomainName: m.MatchedDomainName,
	}
}

type MatchedDomainBQ struct {
	ID                int       `bigquery:"id"`
	CreatedAt         time.Time `bigquery:"created_at"`
	UpdatedAt         time.Time `bigquery:"updated_at"`
	DomainName        string    `bigquery:"domain_name"`
	MatchedDomainName string    `bigquery:"matched_domain_name"`
}

func (m *MatchedDomainBQ) ToGorm() MatchedDomain {
	return MatchedDomain{
		Model: gorm.Model{
			CreatedAt: m.CreatedAt,
			UpdatedAt: m.UpdatedAt,
		},
		DomainName:        m.DomainName,
		MatchedDomainName: m.MatchedDomainName,
	}
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
