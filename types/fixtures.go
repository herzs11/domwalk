package types

import (
	"time"

	"domwalk/db"
)

type MatchedDomain struct {
	CreatedAt         time.Time `bigquery:"created_at"`
	UpdatedAt         time.Time `bigquery:"updated_at"`
	DomainName        string    `json:"reqDomainName,omitempty" bigquery:"domain_name" gorm:"index:dom_match,unique"`
	MatchedDomainName string    `json:"certSanDomainName,omitempty" bigquery:"matched_domain_name" gorm:"index:dom_match,unique"`
	Domain            Domain    `json:"certSanDomain,omitempty" gorm:"foreignKey:MatchedDomainName" bigquery:"-"`
}

func ClearTables() {
	db.GormDB.Migrator().DropTable(
		&Domain{}, &MXRecord{}, &ARecord{}, &AAAARecord{}, &SOARecord{},
		&WebRedirect{}, &MatchedDomain{}, &Sitemap{}, &SitemapWebDomain{},
		&SitemapContactDomain{},
	)
}

func CreateTables() {
	db.GormDB.Migrator().AutoMigrate(
		&Domain{}, &MXRecord{}, &ARecord{}, &AAAARecord{}, &SOARecord{},
		&WebRedirect{}, &MatchedDomain{}, &Sitemap{}, &SitemapWebDomain{},
		&SitemapContactDomain{},
	)
}
