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
)

var UPDATE_MAP = map[string]interface{}{
	"updated_at":             gorm.Expr("MAX(updated_at, excluded.updated_at)"),
	"last_ran_web_redirect":  gorm.Expr("MAX(last_ran_web_redirect, excluded.last_ran_web_redirect)"),
	"last_ran_dns":           gorm.Expr("MAX(last_ran_dns, excluded.last_ran_dns)"),
	"last_ran_cert_sans":     gorm.Expr("MAX(last_ran_cert_sans, excluded.last_ran_cert_sans)"),
	"last_ran_sitemap_parse": gorm.Expr("MAX(last_ran_sitemap_parse, excluded.last_ran_sitemap_parse)"),
}

func main() {
	// Add flags for number of workers, whether or not to clear local db
	var input_file string
	clr := flag.Bool("clear", false, "Clear local database")
	workers := flag.Int("workers", 15, "Number of workers to use")
	enrich := flag.Bool("enrich", false, "Enrich data")
	flag.StringVar(&input_file, "file", "", "File to read domains from")
	flag.Parse()
	if *clr {
		types.ClearTables()
	}
	if input_file == "" && !*enrich {
		log.Fatal("File flag of enrich flag is required")
	}
	if input_file != "" {
		executeDomainsFromFile(input_file, *workers)
	}
	if *enrich {
		enrichDomainsFromDB(*workers)
	}

}

func enrichDomainsFromDB(workers int) {
	types.CreateTables()
	// Use worker pool with 15 workers to execute all domains
	doms := []types.Domain{}
	db.GormDB.Where("last_ran_cert_sans = '0001-01-01 00:00:00+00:00'").Limit(1000).Find(&doms)
	jobs := make(chan *types.Domain, len(doms))
	var wg sync.WaitGroup
	wg.Add(workers)
	for w := 1; w <= workers; w++ {
		go enrichDomainWorker(w, jobs, &wg)
	}
	for _, dom := range doms {
		jobs <- &dom
	}
	close(jobs)
	wg.Wait()
}

func executeDomainsFromFile(filename string, workers int) {
	f, err := os.Open(filename)
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

	types.CreateTables()
	// Use worker pool with 15 workers to execute all domains
	doms := data[1:]
	jobs := make(chan string, len(doms))
	var wg sync.WaitGroup
	wg.Add(workers)
	for w := 1; w <= workers; w++ {
		go createNewDomainWorker(w, jobs, &wg)
	}
	for _, dom := range doms {
		jobs <- dom[0]
	}
	close(jobs)
	wg.Wait()
}

func createNewDomainWorker(id int, jobs <-chan string, wg *sync.WaitGroup) {
	defer wg.Done()
	for domain := range jobs {
		d, err := types.NewDomain(domain)
		// skip domains that are already in the table
		if db.GormDB.First(&types.Domain{}, "domain_name = ?", d.DomainName).Error == nil {
			log.Printf("Domain %s already in database\n", domain)
			continue
		}
		if err != nil {
			log.Printf("Error parsing domain %s: %s\n", domain, err)
		}
		d.GetDNSRecords()
		d.GetRedirectDomains()
		d.GetCertSANs()
		d.GetDomainsFromSitemap()
		db.Mut.Lock()
		db.GormDB.Session(&gorm.Session{FullSaveAssociations: true}).Save(&d)
		db.Mut.Unlock()
		log.Printf("Worker %d Finished processing domain: %s\n", id, domain)
	}
}

func enrichDomainWorker(id int, jobs <-chan *types.Domain, wg *sync.WaitGroup) {
	defer wg.Done()
	for d := range jobs {
		log.Printf("Worker %d started processing domain: %s\n", id, d.DomainName)
		d.GetDNSRecords()
		d.GetRedirectDomains()
		d.GetCertSANs()
		d.GetDomainsFromSitemap()
		db.Mut.Lock()
		db.GormDB.Save(d)
		db.Mut.Unlock()
	}
}
