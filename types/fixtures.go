package types

import (
	"time"

	"domwalk/db"
)

type MatchedDomain struct {
	CreatedAt         time.Time `bigquery:"created_at"`
	UpdatedAt         time.Time `bigquery:"updated_at"`
	DomainName        string    `json:"reqDomainName,omitempty" bigquery:"domain_name" gorm:"index:,unique,composite:dom_match"`
	MatchedDomainName string    `json:"certSanDomainName,omitempty" bigquery:"matched_domain_name" gorm:"index:,unique,composite:dom_match"`
	Domain            Domain    `json:"certSanDomain,omitempty" gorm:"foreignKey:MatchedDomainName" bigquery:"-"`
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
	// db.GormDB.Find(&WebRedirectDomain{}).
}
