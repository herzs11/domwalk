package main

import (
	"encoding/csv"
	"flag"
	"log"
	"os"
	"sync"

	"domwalk/db"
	"domwalk/types"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func clearTables() {
	db.GormDB.Migrator().DropTable(
		&types.Domain{}, &types.MXRecord{}, &types.ARecord{}, &types.AAAARecord{}, &types.SOARecord{},
		&types.WebRedirect{}, &types.CertSAN{}, &types.Sitemap{}, &types.SitemapWebDomain{},
		&types.SitemapContactDomain{},
	)
}

func main() {
	// Add flags for number of workers, whether or not to clear local db
	var input_file string
	clr := flag.Bool("clear", false, "Clear local database")
	workers := flag.Int("workers", 15, "Number of workers to use")
	flag.StringVar(&input_file, "file", "", "File to read domains from")
	flag.Parse()
	if *clr {
		clearTables()
	}
	if input_file == "" {
		log.Fatal("File flag is required")
	}
	f, err := os.Open(input_file)
	if err != nil {
		log.Fatal(err)
	}

	// remember to close the file at the end of the program
	defer f.Close()

	// read csv values using csv.Reader
	csvReader := csv.NewReader(f)
	data, err := csvReader.ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	err = db.GormDB.AutoMigrate(
		&types.Domain{}, &types.MXRecord{}, &types.ARecord{}, &types.AAAARecord{}, &types.SOARecord{},
		&types.WebRedirect{}, &types.CertSAN{}, &types.Sitemap{}, &types.SitemapWebDomain{},
		&types.SitemapContactDomain{},
	)
	if err != nil {
		panic(err)
	}
	// Use worker pool with 15 workers to execute all domains
	doms := data[200:300]
	jobs := make(chan string, len(doms))
	var wg sync.WaitGroup
	wg.Add(*workers)
	for w := 1; w <= *workers; w++ {
		go executeDomainWorker(w, jobs, &wg)
	}
	for _, dom := range doms {
		jobs <- dom[0]
	}
	close(jobs)
	wg.Wait()
}

func executeDomainWorker(id int, jobs <-chan string, wg *sync.WaitGroup) {
	defer wg.Done()
	for domain := range jobs {
		d, err := types.NewDomain(domain)
		// skip domains that are already in the table
		// if db.GormDB.First(&types.Domain{}, "domain_name = ?", d.DomainName).Error == nil {
		// 	log.Printf("Domain %s already in database\n", domain)
		// 	continue
		// }
		if err != nil {
			log.Printf("Error parsing domain %s: %s\n", domain, err)
		}
		d.GetDNSRecords()
		d.GetRedirectDomains()
		d.GetCertSANs()
		d.GetDomainsFromSitemap()
		db.Mut.Lock()
		db.GormDB.Clauses(
			clause.OnConflict{
				UpdateAll: true,
			},
		).Session(&gorm.Session{FullSaveAssociations: true}).Create(&d)
		db.Mut.Unlock()
		log.Printf("Worker %d Finished processing domain: %s\n", id, domain)
	}
}
