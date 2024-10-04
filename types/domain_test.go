package types

import (
	"log"
	"testing"
	"time"

	"domwalk/db"
	"gorm.io/gorm"
)

func TestMain(m *testing.M) {
	ClearTables()
	CreateTables()
	m.Run()
	gdb, err := db.GormDB.DB()
	if err != nil {
		log.Fatal(err)
	}
	gdb.Close()
}

func TestDomains(t *testing.T) {
	doms := []string{
		// "hanselhonda.com",
		"levi.com",
		// "cetac.com",
	}

	for _, dom := range doms {
		d, err := NewDomain("levistrauss.com")
		// skip domains that are already in the table
		if err != nil {
			t.Errorf("Error parsing domain %s: %s\n", dom, err)
		}
		d.GetDNSRecords()
		d.GetRedirectDomains()
		d.GetCertSANs()
		d.GetDomainsFromSitemap()
		d.CertSANs = append(
			d.CertSANs, CertSansDomain{
				MatchedDomain{
					DomainName:        d.DomainName,
					MatchedDomainName: "teledynecetac.com",
				},
			},
		)
		db.Mut.Lock()

		db.GormDB.Session(&gorm.Session{FullSaveAssociations: true}).Save(d)
		time.Sleep(1 * time.Second)
		d2, err := NewDomain(dom)
		if err != nil {
			t.Errorf("Error parsing domain %s: %s\n", d2.DomainName, err)
		}
		d2.GetCertSANs()
		db.GormDB.Session(&gorm.Session{FullSaveAssociations: true}).Save(d2)
		db.Mut.Unlock()
	}

}
