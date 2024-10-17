package domains

import (
	"time"

	"domwalk/db"
)

type MatchedDomain struct {
	CreatedAt  time.Time `json:"createdAt,omitempty"`
	UpdatedAt  time.Time `json:"updatedAt,omitempty"`
	DomainName string    `json:"matchedDomain,omitempty"`
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
