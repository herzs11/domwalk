package types

import (
	"fmt"
	"log"
	"testing"
	"time"

	"domwalk/db"
	"gorm.io/gorm"
)

func TestMain(m *testing.M) {
	ClearTables()
	CreateTables()
	fmt.Println("HERE")
	m.Run()
	gdb, err := db.GormDB.DB()
	if err != nil {
		log.Fatal(err)
	}
	gdb.Close()
}

func TestDomains(t *testing.T) {
	doms := []string{
		"hanselhonda.com",
		"levi.com",
		"cetac.com",
	}

	for _, dom := range doms {
		d, err := NewDomain(dom)
		// skip domains that are already in the table
		if err != nil {
			t.Errorf("Error parsing domain %s: %s\n", dom, err)
		}
		d.GetDNSRecords()
		d.GetRedirectDomains()
		d.GetCertSANs()
		d.GetDomainsFromSitemap()
		db.Mut.Lock()

		db.GormDB.Session(&gorm.Session{FullSaveAssociations: true}).Save(d)
		time.Sleep(1 * time.Second)
		// d2, err := NewDomain("levistrauss.com")
		// if err != nil {
		// 	t.Errorf("Error parsing domain %s: %s\n", d2.DomainName, err)
		// }
		// d2.GetCertSANs()
		// db.GormDB.Session(&gorm.Session{FullSaveAssociations: true}).Save(d2)
		db.Mut.Unlock()
	}

}
