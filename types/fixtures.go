package types

import (
	"domwalk/db"
)

func ClearTables() {
	db.GormDB.Migrator().DropTable(
		&Domain{}, &MXRecord{}, &ARecord{}, &AAAARecord{}, &SOARecord{},
		&WebRedirect{}, &CertSAN{}, &Sitemap{}, &SitemapWebDomain{},
		&SitemapContactDomain{},
	)
}

func CreateTables() {
	db.GormDB.Migrator().AutoMigrate(
		&Domain{}, &MXRecord{}, &ARecord{}, &AAAARecord{}, &SOARecord{},
		&WebRedirect{}, &CertSAN{}, &Sitemap{}, &SitemapWebDomain{},
		&SitemapContactDomain{},
	)
}
