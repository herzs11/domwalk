package main

import (
	"log"
	"sync"

	"domwalk/db"
	"domwalk/types"
	"gorm.io/gorm"
)

func clearTables() {
	db.GormDB.Migrator().DropTable(
		&types.Domain{}, &types.MXRecord{}, &types.ARecord{}, &types.AAAARecord{}, &types.SOARecord{},
		&types.WebRedirect{}, &types.CertSAN{}, &types.Sitemap{}, &types.SitemapWebDomain{},
		&types.SitemapContactDomain{},
	)
}

func main() {

	err := db.GormDB.AutoMigrate(
		&types.Domain{}, &types.MXRecord{}, &types.ARecord{}, &types.AAAARecord{}, &types.SOARecord{},
		&types.WebRedirect{}, &types.CertSAN{}, &types.Sitemap{}, &types.SitemapWebDomain{},
		&types.SitemapContactDomain{},
	)
	if err != nil {
		panic(err)
	}
	doms := []string{
		"fleetup.com",
		"google.com",
		"cetac.com",
		"levi.com",
	}
	wg := sync.WaitGroup{}
	for _, dom := range doms {
		wg.Add(1)
		go executeDomain(dom, &wg)
	}
	wg.Wait()
}

func executeDomain(domain string, wg *sync.WaitGroup) {
	log.Printf("Processing domain: %s\n", domain)
	d, err := types.NewDomain(domain)
	if err != nil {
		panic(err)
	}
	d.GetDNSRecords()
	d.GetRedirectDomains()
	d.GetCertSANs()
	d.GetDomainsFromSitemap()
	db.GormDB.Session(&gorm.Session{FullSaveAssociations: true}).Create(d)
	log.Printf("Finished processing domain: %s\n", domain)
	wg.Done()
}
