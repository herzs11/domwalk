package cmd

import (
	"encoding/csv"
	"errors"
	"fmt"
	"os"
	"sync"
	"time"

	"domwalk/db"
	"domwalk/types"
	"github.com/fatih/color"
	"gorm.io/gorm"
)

var (
	N_TO_PROCESS  int
	NUM_PROCESSED int
)

type enrichmentConfig struct {
	CertSans         bool
	DNS              bool
	SitemapWeb       bool
	SitemapContact   bool
	WebRedirect      bool
	Limit            int
	Offset           int
	NWorkers         int
	MinFreshnessDate time.Time
}

func readDomainsFromFile(filename string) ([]string, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("error opening file: %w", err)
	}

	// remember to close the file at the end of the program
	defer f.Close()

	// read csv values using csv.Reader
	csvReader := csv.NewReader(f)
	data, err := csvReader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("error reading csv data: %w", err)
	}
	doms := make([]string, 0, len(data))
	for _, row := range data {
		doms = append(doms, row[0])
	}
	return doms, nil
}

func enrichDomainNames(domains []string, cfg enrichmentConfig) {
	doms := []*types.Domain{}
	for _, dom := range domains {
		d, err := types.NewDomain(dom)
		if err != nil {
			color.Yellow("Error parsing domain %s: %s\n", dom, err)
		}
		doms = append(doms, d)
	}
	enrichDomains(doms, cfg)
}

func enrichDBDomains(cfg enrichmentConfig) {
	if cfg.DNS {
		color.Yellow("Enriching domains with DNS records\n")
		var domains []*types.Domain
		db.GormDB.Where("last_ran_dns <= ?", cfg.MinFreshnessDate).Find(&domains)
		enrichDomains(
			domains, enrichmentConfig{
				DNS:              true,
				Limit:            cfg.Limit,
				Offset:           cfg.Offset,
				NWorkers:         cfg.NWorkers,
				MinFreshnessDate: cfg.MinFreshnessDate,
			},
		)
	}
	if cfg.WebRedirect {
		color.Yellow("Enriching domains with web redirects\n")
		var domains []*types.Domain
		db.GormDB.Where("last_ran_web_redirect <= ?", cfg.MinFreshnessDate).Find(&domains)
		enrichDomains(
			domains, enrichmentConfig{
				WebRedirect:      true,
				Limit:            cfg.Limit,
				Offset:           cfg.Offset,
				NWorkers:         cfg.NWorkers,
				MinFreshnessDate: cfg.MinFreshnessDate,
			},
		)
	}
	if cfg.CertSans {
		color.Yellow("Enriching domains with certificate SANs\n")
		var domains []*types.Domain
		db.GormDB.Where("last_ran_cert_sans <= ?", cfg.MinFreshnessDate).Find(&domains)
		enrichDomains(
			domains, enrichmentConfig{
				CertSans:         true,
				Limit:            cfg.Limit,
				Offset:           cfg.Offset,
				NWorkers:         cfg.NWorkers,
				MinFreshnessDate: cfg.MinFreshnessDate,
			},
		)
	}
	if cfg.SitemapWeb {
		color.Yellow("Enriching domains with web domains from sitemaps\n")
		var domains []*types.Domain
		db.GormDB.Where("last_ran_sitemap_parse <= ?", cfg.MinFreshnessDate).Find(&domains)
		enrichDomains(
			domains, enrichmentConfig{
				SitemapWeb:       true,
				Limit:            cfg.Limit,
				Offset:           cfg.Offset,
				NWorkers:         cfg.NWorkers,
				MinFreshnessDate: cfg.MinFreshnessDate,
			},
		)
	}
}

func enrichDomains(domains []*types.Domain, cfg enrichmentConfig) {
	if cfg.Offset+cfg.Limit > len(domains) {
		cfg.Limit = len(domains) - cfg.Offset
	}
	domains = domains[cfg.Offset : cfg.Offset+cfg.Limit]
	N_TO_PROCESS = len(domains)
	NUM_PROCESSED = 0
	jobs := make(chan *types.Domain, len(domains))
	var wg sync.WaitGroup
	wg.Add(cfg.NWorkers)
	for w := 1; w <= cfg.NWorkers; w++ {
		go enrichDomainWorker(w, jobs, &wg, cfg)
	}
	for _, dom := range domains {
		jobs <- dom
	}
	close(jobs)
	wg.Wait()
}

func enrichDomainWorker(id int, jobs <-chan *types.Domain, wg *sync.WaitGroup, cfg enrichmentConfig) {
	defer wg.Done()
	for domain := range jobs {
		var d = types.Domain{}
		if err := db.GormDB.Where("domain_name = ?", domain.DomainName).First(&d).Error; errors.Is(
			err, gorm.ErrRecordNotFound,
		) {
			d = *domain
		}
		if d.LastRanDns.Unix() <= cfg.MinFreshnessDate.Unix() && cfg.DNS {
			d.GetDNSRecords()
		}
		if d.LastRanWebRedirect.Unix() <= cfg.MinFreshnessDate.Unix() && cfg.WebRedirect {
			d.GetRedirectDomains()
		}
		if d.LastRanCertSans.Unix() <= cfg.MinFreshnessDate.Unix() && cfg.CertSans {
			d.GetCertSANs()
		}
		if d.LastRanSitemapParse.Unix() <= cfg.MinFreshnessDate.Unix() && cfg.SitemapWeb {
			d.GetDomainsFromSitemap()
		}
		db.Mut.Lock()
		db.GormDB.Session(&gorm.Session{FullSaveAssociations: true}).Save(&d)
		NUM_PROCESSED++
		color.Green("Worker %d: Processed domain %s, %d out of %d\n", id, domain, NUM_PROCESSED, N_TO_PROCESS)
		db.Mut.Unlock()
	}
}
