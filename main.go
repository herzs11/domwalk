package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"os"
	"sync"

	"domwalk/db"
	"domwalk/types"
	"gorm.io/gorm"
)

var (
	N_PROCESSED  = 0
	N_TO_PROCESS = 0
	UPDATE_MAP   = map[string]interface{}{
		"updated_at":             gorm.Expr("MAX(updated_at, excluded.updated_at)"),
		"last_ran_web_redirect":  gorm.Expr("MAX(last_ran_web_redirect, excluded.last_ran_web_redirect)"),
		"last_ran_dns":           gorm.Expr("MAX(last_ran_dns, excluded.last_ran_dns)"),
		"last_ran_cert_sans":     gorm.Expr("MAX(last_ran_cert_sans, excluded.last_ran_cert_sans)"),
		"last_ran_sitemap_parse": gorm.Expr("MAX(last_ran_sitemap_parse, excluded.last_ran_sitemap_parse)"),
	}
	Reset = "\033[0m"
	Green = "\033[32m"
)

func main() {
	// Add flags for number of workers, whether or not to clear local db
	var input_file string
	clr := flag.Bool("clear", false, "Clear local database")
	workers := flag.Int("workers", 15, "Number of workers to use")
	enrich := flag.Bool("enrich", false, "Enrich data")
	dns := flag.Bool("dns", false, "Enrich DNS data")
	cert_sans := flag.Bool("certSans", false, "Enrich Cert SANs data")
	web_redirect := flag.Bool("redirect", false, "Enrich Web redirect data")
	domain := flag.String("domain", "", "Domain to process")
	flag.StringVar(&input_file, "file", "", "File to read domains from")
	flag.Parse()
	if *clr {
		types.ClearTables()
	}
	if *domain != "" {
		d, err := types.NewDomain(*domain)
		if err != nil {
			log.Fatalf("Error parsing domain %s: %s\n", *domain, err)
		}
		d.GetDNSRecords()
		d.GetRedirectDomains()
		d.GetCertSANs()
		d.GetDomainsFromSitemap()
		db.Mut.Lock()
		db.GormDB.Session(&gorm.Session{FullSaveAssociations: true}).Save(&d)
		db.Mut.Unlock()
		return
	}
	if input_file != "" {
		executeDomainsFromFile(input_file, *workers)
		return
	}
	if *enrich {
		enrichDomainsFromDB(*workers)
		return
	}
	if *web_redirect {
		N_PROCESSED = 0
		enrichWebRedirectDomainsFromDB(*workers)
	}
	if *dns {
		N_PROCESSED = 0
		enrichDNSDomainsFromDB(*workers)
	}
	if *cert_sans {
		N_PROCESSED = 0
		enrichCertSansDomainsFromDB(*workers)
	}
}

func enrichCertSansDomainsFromDB(workers int) {
	types.CreateTables()
	// Use worker pool with 15 workers to execute all domains
	doms := []types.Domain{}
	db.GormDB.Where("last_ran_cert_sans = '0001-01-01 00:00:00+00:00'").Limit(3000).Find(&doms)
	N_TO_PROCESS = len(doms)
	jobs := make(chan *types.Domain, len(doms))
	var wg sync.WaitGroup
	wg.Add(workers)
	for w := 1; w <= workers; w++ {
		go enrichCertSansDomainWorker(w, jobs, &wg)
	}
	for _, dom := range doms {
		jobs <- &dom
	}
	close(jobs)
	wg.Wait()
}

func enrichDNSDomainsFromDB(workers int) {
	types.CreateTables()
	// Use worker pool with 15 workers to execute all domains
	doms := []types.Domain{}
	db.GormDB.Where("last_ran_dns = '0001-01-01 00:00:00+00:00'").Limit(3000).Find(&doms)
	N_TO_PROCESS = len(doms)
	jobs := make(chan *types.Domain, len(doms))
	var wg sync.WaitGroup
	wg.Add(workers)
	for w := 1; w <= workers; w++ {
		go enrichDNSDomainWorker(w, jobs, &wg)
	}
	for _, dom := range doms {
		jobs <- &dom
	}
	close(jobs)
	wg.Wait()
}

func enrichDomainsFromDB(workers int) {
	types.CreateTables()
	// Use worker pool with 15 workers to execute all domains
	doms := []types.Domain{}
	db.GormDB.Where("last_ran_cert_sans = '0001-01-01 00:00:00+00:00'").Limit(1000).Find(&doms)
	N_TO_PROCESS = len(doms)
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
	doms := data[6000:15000]
	N_TO_PROCESS = len(doms)
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
		if db.GormDB.Where("domain_name = ?", domain).First(&types.Domain{}).Error == nil {
			db.Mut.Lock()
			N_PROCESSED++
			db.Mut.Unlock()
			fmt.Printf("%sDomain Already Processed: %d out of %d %s\n", Green, N_PROCESSED, N_TO_PROCESS, Reset)
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
		N_PROCESSED++
		fmt.Printf("%sProcessed: %d out of %d %s\n", Green, N_PROCESSED, N_TO_PROCESS, Reset)
		db.Mut.Unlock()
	}
}

func enrichDomainWorker(id int, jobs <-chan *types.Domain, wg *sync.WaitGroup) {
	defer wg.Done()
	for d := range jobs {
		d.GetDNSRecords()
		d.GetRedirectDomains()
		d.GetCertSANs()
		d.GetDomainsFromSitemap()
		db.Mut.Lock()
		db.GormDB.Save(d)
		N_PROCESSED++
		fmt.Printf("%sProcessed: %d out of %d %s\n", Green, N_PROCESSED, N_TO_PROCESS, Reset)
		db.Mut.Unlock()
	}
}

func enrichDNSDomainWorker(id int, jobs <-chan *types.Domain, wg *sync.WaitGroup) {
	defer wg.Done()
	for d := range jobs {
		d.GetDNSRecords()
		db.Mut.Lock()
		db.GormDB.Save(d)
		N_PROCESSED++
		fmt.Printf("%sProcessed: %d out of %d %s\n", Green, N_PROCESSED, N_TO_PROCESS, Reset)
		db.Mut.Unlock()
	}
}

func enrichCertSansDomainWorker(id int, jobs <-chan *types.Domain, wg *sync.WaitGroup) {
	defer wg.Done()
	for d := range jobs {
		d.GetCertSANs()
		db.Mut.Lock()
		db.GormDB.Save(d)
		N_PROCESSED++
		fmt.Printf("%sProcessed: %d out of %d %s\n", Green, N_PROCESSED, N_TO_PROCESS, Reset)
		db.Mut.Unlock()
	}
}

func enrichWebRedirectDomainsFromDB(workers int) {
	types.CreateTables()
	// Use worker pool with 15 workers to execute all domains
	doms := []types.Domain{}
	db.GormDB.Where("last_ran_web_redirect = '0001-01-01 00:00:00+00:00'").Limit(1000).Find(&doms)
	N_TO_PROCESS = len(doms)
	jobs := make(chan *types.Domain, len(doms))
	var wg sync.WaitGroup
	wg.Add(workers)
	for w := 1; w <= workers; w++ {
		go enrichWebRedirectDomainWorker(w, jobs, &wg)
	}
	for _, dom := range doms {
		jobs <- &dom
	}
	close(jobs)
	wg.Wait()
}

func enrichWebRedirectDomainWorker(id int, jobs <-chan *types.Domain, wg *sync.WaitGroup) {
	defer wg.Done()
	for d := range jobs {
		d.GetRedirectDomains()
		db.Mut.Lock()
		db.GormDB.Save(d)
		N_PROCESSED++
		fmt.Printf("%sProcessed: %d out of %d %s\n", Green, N_PROCESSED, N_TO_PROCESS, Reset)
		db.Mut.Unlock()
	}
}
