package cmd

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"cloud.google.com/go/bigquery"
	"domwalk/db"
	"domwalk/types"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/idtoken"
	"google.golang.org/api/option"
)

const SNAPSHOT_JOB_URL = "https://snapshot-domwalk-593149879404.us-central1.run.app"

var dataset *bigquery.Dataset

type syncConfig struct {
	Dataset            string
	Domains            bool
	CertSansDomains    bool
	DNS                bool
	WebRedirectDomains bool
	Sitemaps           bool
}

func pullFromBQ(cfg syncConfig) {
	fmt.Println("Pulling data from BigQuery")
	// TODO: Implement pulling data from BigQuery
}

func pushToBQ(cfg syncConfig) {
	dataset = db.BQConn.Dataset(cfg.Dataset)
	defer db.BQConn.Close()
	var (
		offset    int
		limit     int
		tableName string
	)

	if cfg.Domains {
		tableName = "domains"
		if err := recreateTable(types.DomainBQ{}, tableName); err != nil {
			log.Printf("Failed to truncate table: %v", err)
		}
		offset = 0
		limit = 1000
		for {
			chunk := []types.Domain{}
			err := db.GormDB.Limit(limit).Offset(offset).Find(&chunk).Error
			if err != nil {
				log.Fatalf("Failed to get chunk: %v", err)
			}
			if len(chunk) == 0 {
				break
			}
			bqds := make([]types.DomainBQ, len(chunk))
			for i, m := range chunk {
				bqds[i] = m.ToBQ()
			}
			loadToBigQuery(bqds, tableName)
			offset += limit
		}
	}
	if cfg.CertSansDomains {
		offset = 0
		limit = 5000
		tableName = "cert_sans_domains"
		if err := recreateTable(types.MatchedDomainBQ{}, tableName); err != nil {
			log.Printf("Failed to truncate table: %v", err)
		}
		for {
			chunk := []types.CertSansDomain{}
			bqr := []types.MatchedDomainBQ{}
			err := db.GormDB.Limit(limit).Offset(offset).Find(&chunk).Error
			if err != nil {
				log.Fatalf("Failed to get chunk: %v", err)
			}
			if len(chunk) == 0 {
				break
			}
			for _, m := range chunk {
				bqr = append(bqr, m.ToBQ())
			}
			loadToBigQuery(bqr, tableName)
			offset += limit
		}
	}
	if cfg.DNS {
		tableName = "mx_records"
		offset = 0
		limit = 5000
		if err := recreateTable(types.MXRecordBQ{}, tableName); err != nil {
			log.Printf("Failed to truncate table: %v", err)
		}
		for {
			bqr := []types.MXRecordBQ{}
			chunk := []types.MXRecord{}
			err := db.GormDB.Limit(limit).Offset(offset).Find(&chunk).Error
			if err != nil {
				log.Fatalf("Failed to get chunk: %v", err)
			}
			if len(chunk) == 0 {
				break
			}
			for _, m := range chunk {
				bqr = append(bqr, m.ToBQ())
			}
			loadToBigQuery(bqr, tableName)
			offset += limit
		}

		tableName = "a_records"
		offset = 0
		limit = 5000
		if err := recreateTable(types.ARecordBQ{}, tableName); err != nil {
			log.Printf("Failed to truncate table: %v", err)
		}
		for {
			bqr := []types.ARecordBQ{}
			chunk := []types.ARecord{}
			err := db.GormDB.Limit(limit).Offset(offset).Find(&chunk).Error
			if err != nil {
				log.Fatalf("Failed to get chunk: %v", err)
			}
			if len(chunk) == 0 {
				break
			}
			for _, m := range chunk {
				bqr = append(bqr, m.ToBQ())
			}
			loadToBigQuery(bqr, tableName)
			offset += limit

		}

		tableName = "aaaa_records"
		offset = 0
		limit = 1000
		if err := recreateTable(types.AAAARecordBQ{}, tableName); err != nil {
			log.Fatalf("Failed to truncate table: %v", err)
		}
		for {
			bqr := []types.AAAARecordBQ{}
			chunk := []types.AAAARecord{}
			err := db.GormDB.Limit(limit).Offset(offset).Find(&chunk).Error
			if err != nil {
				log.Fatalf("Failed to get chunk: %v", err)
			}
			if len(chunk) == 0 {
				break
			}
			for _, m := range chunk {
				bqr = append(bqr, m.ToBQ())
			}
			loadToBigQuery(bqr, tableName)
			offset += limit
		}

		tableName = "soa_records"
		offset = 0
		limit = 1000
		if err := recreateTable(types.SOARecordBQ{}, tableName); err != nil {
			log.Fatalf("Failed to truncate table: %v", err)
		}
		for {
			bqr := []types.SOARecordBQ{}
			chunk := []types.SOARecord{}
			err := db.GormDB.Limit(limit).Offset(offset).Find(&chunk).Error
			if err != nil {
				log.Fatalf("Failed to get chunk: %v", err)
			}
			if len(chunk) == 0 {
				break
			}
			for _, m := range chunk {
				bqr = append(bqr, m.ToBQ())
			}
			loadToBigQuery(bqr, tableName)
			offset += limit

		}
	}
	if cfg.WebRedirectDomains {
		tableName = "web_redirect_domains"
		offset = 0
		limit = 1000
		if err := recreateTable(types.MatchedDomainBQ{}, tableName); err != nil {
			log.Printf("Failed to truncate table: %v", err)
		}
		for {
			chunk := []types.WebRedirectDomain{}
			bqr := []types.MatchedDomainBQ{}
			err := db.GormDB.Limit(limit).Offset(offset).Find(&chunk).Error
			if err != nil {
				log.Fatalf("Failed to get chunk: %v", err)
			}
			if len(chunk) == 0 {
				break
			}
			for _, m := range chunk {
				bqr = append(bqr, m.ToBQ())
			}
			loadToBigQuery(bqr, tableName)
			offset += limit

		}
	}
	if cfg.Sitemaps {
		tableName = "sitemaps"
		offset = 0
		limit = 1000
		if err := recreateTable(types.SitemapBQ{}, tableName); err != nil {
			log.Fatalf("Failed to truncate table: %v", err)
		}
		for {
			chunk := []types.Sitemap{}
			bqr := []types.SitemapBQ{}
			err := db.GormDB.Limit(limit).Offset(offset).Find(&chunk).Error
			if err != nil {
				log.Fatalf("Failed to get chunk: %v", err)
			}
			if len(chunk) == 0 {
				break
			}
			for _, m := range chunk {
				bqr = append(bqr, m.ToBQ())
			}
			loadToBigQuery(bqr, tableName)
			offset += limit
		}

		tableName = "sitemap_web_domains"
		offset = 0
		limit = 5000
		if err := recreateTable(types.MatchedDomainBQ{}, tableName); err != nil {
			log.Fatalf("Failed to truncate table: %v", err)
		}
		for {
			chunk := []types.SitemapWebDomain{}
			bqr := []types.MatchedDomainBQ{}
			err := db.GormDB.Limit(limit).Offset(offset).Find(&chunk).Error
			if err != nil {
				log.Fatalf("Failed to get chunk: %v", err)
			}
			if len(chunk) == 0 {
				break
			}
			for _, m := range chunk {
				bqr = append(bqr, m.ToBQ())
			}
			loadToBigQuery(bqr, tableName)
			offset += limit

		}

		tableName = "sitemap_contact_domains"
		offset = 0
		limit = 5000
		if err := recreateTable(types.MatchedDomainBQ{}, tableName); err != nil {
			log.Fatalf("Failed to truncate table: %v", err)
		}
		for {
			chunk := []types.SitemapContactDomain{}
			bqr := []types.MatchedDomainBQ{}
			err := db.GormDB.Limit(limit).Offset(offset).Find(&chunk).Error
			if err != nil {
				log.Fatalf("Failed to get chunk: %v", err)
			}
			if len(chunk) == 0 {
				break
			}
			for _, m := range chunk {
				bqr = append(bqr, m.ToBQ())
			}
			loadToBigQuery(bqr, tableName)
			offset += limit
		}
	}
}

func loadToBigQuery(model interface{}, tableName string) {

	ctx := context.Background()
	table := dataset.Table(tableName)

	inserter := table.Inserter()
	for {
		if err := inserter.Put(ctx, model); err != nil {
			log.Printf("Failed to insert data into BigQuery table: %v, sleeping for 3 seconds\n", err)
			time.Sleep(3 * time.Second)

		} else {
			break
		}
	}

	fmt.Printf("Data loaded into BigQuery table: %s\n", table.FullyQualifiedName())
}

func recreateTable(model interface{}, tableName string) error {
	ctx := context.Background()
	table := dataset.Table(tableName)
	if err := table.Delete(ctx); err != nil {
		log.Printf("Failed to delete table: %v, trying to create\n", err)
	}
	time.Sleep(5 * time.Second)
	sch, err := bigquery.InferSchema(model)
	if err != nil {
		log.Printf("Failed to infer schema: %v", err)
		return err
	}
	metaData := &bigquery.TableMetadata{
		Schema: sch,
	}
	if err := table.Create(ctx, metaData); err != nil {
		log.Printf("Failed to create table: %v", err)
		return err
	}
	fmt.Printf("Table %s created\n", tableName)
	return nil
}

func snapshotDomains() error {

	// client, err := google.DefaultClient(context.Background(), "https://snapshot-domwalk-593149879404.us-central1.run.app")
	ctx := context.Background()

	// Construct the GoogleCredentials object which obtains the default configuration from your
	// working environment.
	credentials, err := google.FindDefaultCredentials(ctx)
	if err != nil {
		return fmt.Errorf("failed to generate default credentials: %w", err)
	}

	ts, err := idtoken.NewTokenSource(ctx, SNAPSHOT_JOB_URL, option.WithCredentials(credentials))
	if err != nil {
		return fmt.Errorf("failed to create NewTokenSource: %w", err)
	}

	// Get the ID token.
	// Once you've obtained the ID token, you can use it to make an authenticated call
	// to the target audience.
	tok, err := ts.Token()
	if err != nil {
		return fmt.Errorf("failed to receive token: %w", err)
	}
	req, err := http.NewRequest("GET", SNAPSHOT_JOB_URL, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", tok.AccessToken))
	resp, err := http.DefaultClient.Do(req)
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	fmt.Println(string(body))
	return nil
}
