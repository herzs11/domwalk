package cloud_functions

import (
	"sync"

	"domwalk/cloud_functions/types"
	"domwalk/domains"
)

func enrichDomains(doms []*domains.Domain, cfg types.ProcessConfig) {
	jobs := make(chan *domains.Domain, len(doms))
	var wg sync.WaitGroup
	wg.Add(cfg.Workers)
	for w := 1; w <= cfg.Workers; w++ {
		go enrichDomainWorker(w, jobs, &wg, cfg.EnrichmentConfig)
	}
	for _, dom := range doms {
		jobs <- dom
	}
	close(jobs)
	wg.Wait()
}

func enrichDomainWorker(id int, jobs <-chan *domains.Domain, wg *sync.WaitGroup, cfg domains.EnrichmentConfig) {
	defer wg.Done()
	for domain := range jobs {
		domain.Enrich(cfg)
	}
}
