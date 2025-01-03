package bq

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/herzs11/go-doms/domain"
)

func TestDomains(t *testing.T) {
	bqs, err := NewBQStore("unum-marketing-data-assets", "domwalk", "domains")
	if err != nil {
		t.Fatal(err)
	}
	d := []string{
		"piibr.com",
	}
	doms, err := bqs.GetDomainsByNames(context.Background(), d)
	if err != nil {
		t.Fatal(err)
	}
	if len(doms) == 0 {
		t.Fatal("No domains found")
	}
	cfg := domain.EnrichmentConfig{
		CertSans:         true,
		DNS:              true,
		Sitemap:          true,
		WebRedirect:      true,
		MinFreshnessDate: time.Now(),
	}
	for _, dom := range doms {
		dom.Enrich(cfg)
	}
	fmt.Println(doms[0].GetAllMatchedDomains())
}

func TestMergeTable(t *testing.T) {
	bqs, err := NewBQStore("unum-marketing-data-assets", "domwalk", "domains")
	if err != nil {
		t.Fatal(err)
	}
	err = bqs.createDomainTable(context.Background(), "domains")
	if err != nil {
		t.Fatal(err)
	}
}
