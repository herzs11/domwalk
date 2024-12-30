package bq

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"
	
	"cloud.google.com/go/bigquery"
	"github.com/herzs11/domwalk/domains"
	"google.golang.org/api/googleapi"
	"google.golang.org/api/iterator"
)

type BQStore struct {
	Mut *sync.RWMutex
	*bigquery.Client
	Dataset *bigquery.Dataset
	Table   *bigquery.Table
}

func NewBQStore(projectID, datasetID, tableName string) (*BQStore, error) {
	ctx := context.Background()
	client, err := bigquery.NewClient(ctx, projectID)
	if err != nil {
		return nil, err
	}
	dataset := client.Dataset(datasetID)
	if _, err := dataset.Metadata(ctx); err != nil {
		return nil, err
	}
	table := dataset.Table(tableName)
	if _, err := table.Metadata(ctx); err != nil {
		if e, ok := err.(*googleapi.Error); ok {
			if e.Code == http.StatusNotFound {
				log.Printf("Table %s not found, creating it now", tableName)
				schema, err := bigquery.InferSchema(DomainBQ{})
				if err != nil {
					return nil, err
				}
				tableMetadata := &bigquery.TableMetadata{
					Schema: schema,
				}
				if err := table.Create(ctx, tableMetadata); err != nil {
					return nil, err
				}
				time.Sleep(2 * time.Second)
			} else {
				return nil, err
			}
		} else {
			return nil, err
		}
	}
	return &BQStore{
		Mut:     &sync.RWMutex{},
		Client:  client,
		Dataset: dataset,
		Table:   table,
	}, nil
}

func (bq *BQStore) recreateMergeTable(ctx context.Context) error {
	qry := `create table domwalk.domain_mrg
				(
					created_at              TIMESTAMP,
					updated_at              TIMESTAMP,
					domain_name             STRING,
					non_public_domain       BOOL,
					hostname                STRING,
					subdomain               STRING,
					suffix                  STRING,
					successful_web_landing  BOOL,
					web_redirect_url_final  STRING,
					last_ran_web_redirect   TIMESTAMP,
					last_ran_dns            TIMESTAMP,
					last_ran_cert_sans      TIMESTAMP,
					last_ran_sitemap_parse  TIMESTAMP,
					a_records               ARRAY <STRUCT < created_at TIMESTAMP, updated_at TIMESTAMP, ip STRING>>,
					aaaa_records            ARRAY <STRUCT < created_at TIMESTAMP, updated_at TIMESTAMP, ip_v6 STRING>>,
					mx_records              ARRAY <STRUCT < created_at TIMESTAMP, updated_at TIMESTAMP, mx STRING>>,
					soa_records             ARRAY <STRUCT < created_at TIMESTAMP, updated_at TIMESTAMP, ns STRING, mbox STRING,
															serial INT64>>,
					cert_org_names ARRAY<STRING>,
					sitemaps                ARRAY <STRUCT < created_at TIMESTAMP, updated_at TIMESTAMP, sitemap_loc STRING>>,
					web_redirect_domains    ARRAY <STRUCT < created_at TIMESTAMP, updated_at TIMESTAMP, domain_name STRING>>,
					cert_sans               ARRAY <STRUCT < created_at TIMESTAMP, updated_at TIMESTAMP, domain_name STRING>>,
					sitemap_web_domains     ARRAY <STRUCT < created_at TIMESTAMP, updated_at TIMESTAMP, domain_name STRING>>,
					sitemap_contact_domains ARRAY <STRUCT < created_at TIMESTAMP, updated_at TIMESTAMP, domain_name STRING>>
				);`
	_, err := bq.Client.Query(qry).Read(ctx)
	return err
}

func (bq *BQStore) PutDomains(ctx context.Context, doms []*domains.Domain) error {
	var dbq []DomainBQ
	now := time.Now()
	for _, d := range doms {
		d.UpdatedAt = now
		dbq = append(dbq, newDomainBQ(d))
	}
	qry := bq.Client.Query(
		fmt.Sprintf(
			`MERGE INTO %s.%s t
					USING (
						SELECT * FROM UNNEST(@d)
						) s
						ON t.domain_name = s.domain_name
					WHEN MATCHED THEN
						UPDATE SET t.updated_at = GREATEST(t.updated_at, s.updated_at),
									t.last_ran_web_redirect = GREATEST(t.last_ran_web_redirect, s.last_ran_web_redirect),
									t.last_ran_dns = GREATEST(t.last_ran_dns, s.last_ran_dns),
									t.last_ran_cert_sans = GREATEST(t.last_ran_cert_sans, s.last_ran_cert_sans),
									t.last_ran_sitemap_parse = GREATEST(t.last_ran_sitemap_parse, s.last_ran_sitemap_parse),
									t.a_records = s.a_records,
									t.aaaa_records = s.aaaa_records,
									t.mx_records = s.mx_records,
									t.soa_records = s.soa_records,
									t.sitemaps = s.sitemaps,
									t.web_redirect_domains = s.web_redirect_domains,
									t.cert_sans = s.cert_sans,
									t.cert_org_names = s.cert_org_names,
									t.sitemap_web_domains = s.sitemap_web_domains,
									t.sitemap_contact_domains = s.sitemap_contact_domains
					WHEN NOT MATCHED THEN INSERT ROW;`,
			bq.Dataset.DatasetID, bq.Table.TableID,
		),
	)
	qry.Parameters = []bigquery.QueryParameter{
		{Name: "d", Value: dbq},
	}
	bq.Mut.Lock()
	defer bq.Mut.Unlock()
	_, err := qry.Read(ctx)
	return err
}

func (bq *BQStore) GetDomains(ctx context.Context, query string) ([]*domains.Domain, error) {
	var doms []*domains.Domain
	bq.Mut.RLock()
	q := bq.Client.Query(query)
	it, err := q.Read(ctx)
	if err != nil {
		return nil, err
	}
	for {
		var d DomainBQ
		err := it.Next(&d)
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}
		doms = append(doms, d.parse())
	}
	bq.Mut.RUnlock()
	return doms, nil
}

func (bq *BQStore) GetDomainsByNames(ctx context.Context, doms []string) ([]*domains.Domain, error) {
	var domObjs []*domains.Domain
	var domsFound = make(map[string]bool)
	bq.Mut.RLock()
	qry := bq.Client.Query(
		`SELECT * FROM ` + fmt.Sprintf(
			"%s.%s", bq.Dataset.DatasetID, bq.Table.TableID,
		) + ` WHERE domain_name IN UNNEST(@dns)`,
	)
	qry.Parameters = []bigquery.QueryParameter{
		{Name: "dns", Value: doms},
	}
	it, err := qry.Read(ctx)
	if err != nil {
		return nil, err
	}
	for {
		var d DomainBQ
		err := it.Next(&d)
		if errors.Is(err, iterator.Done) {
			break
		}
		if err != nil {
			return nil, err
		}
		domObjs = append(domObjs, d.parse())
		domsFound[d.DomainName] = true
	}
	bq.Mut.RUnlock()
	for _, dom := range doms {
		if _, exists := domsFound[dom]; !exists {
			domsFound[dom] = true
			d, err := domains.NewDomain(dom)
			if err != nil {
				log.Printf("Error parsing domain %s: %s\n", dom, err)
				continue
			}
			domObjs = append(domObjs, d)
		}
	}
	return domObjs, nil
}
