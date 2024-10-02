package types

import (
	"testing"

	"domwalk/db"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func TestMain(m *testing.M) {
	ClearTables()
	CreateTables()
	m.Run()
}

func TestDomains(t *testing.T) {
	doms := []string{
		"google.com",
		"levi.com",
		"cetac.com",
	}

	for _, dom := range doms {
		d, err := NewDomain(dom)
		// skip domains that are already in the table
		// if db.GormDB.First(&types.Domain{}, "domain_name = ?", d.DomainName).Error == nil {
		// 	log.Printf("Domain %s already in database\n", domain)
		// 	continue
		// }
		if err != nil {
			t.Errorf("Error parsing domain %s: %s\n", dom, err)
		}
		d.GetDNSRecords()
		d.GetRedirectDomains()
		d.GetCertSANs()
		d.GetDomainsFromSitemap()
		db.GormDB.Clauses(
			clause.OnConflict{
				UpdateAll: true,
			},
		).Session(&gorm.Session{FullSaveAssociations: true}).Create(&d)
	}
}
