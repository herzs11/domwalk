package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"cloud.google.com/go/bigquery"
	"domwalk/db"
	"domwalk/types"
)

var dataset *bigquery.Dataset

func main() {
	err := db.CreateBigQueryConn()
	if err != nil {
		log.Fatal(err)
	}
	dataset = db.BQConn.Dataset("domwalk")
	defer db.BQConn.Close()

	tableName := "domains"
	if err := recreateTable(types.DomainBQ{}, tableName); err != nil {
		log.Printf("Failed to truncate table: %v", err)
	}
	offset := 0
	limit := 1000
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

	offset = 0
	limit = 5000
	tableName = "cert_sans_domains"
	if err := recreateTable(types.CertSansDomain{}, tableName); err != nil {
		log.Printf("Failed to truncate table: %v", err)
	}
	for {
		chunk := []types.CertSansDomain{}
		err := db.GormDB.Limit(limit).Offset(offset).Find(&chunk).Error
		if err != nil {
			log.Fatalf("Failed to get chunk: %v", err)
		}
		if len(chunk) == 0 {
			break
		}
		loadToBigQuery(chunk, tableName)
		offset += limit
	}

	tableName = "mx_records"
	offset = 0
	limit = 5000
	if err := recreateTable(types.MXRecord{}, tableName); err != nil {
		log.Printf("Failed to truncate table: %v", err)
	}
	for {
		chunk := []types.MXRecord{}
		err := db.GormDB.Limit(limit).Offset(offset).Find(&chunk).Error
		if err != nil {
			log.Fatalf("Failed to get chunk: %v", err)
		}
		if len(chunk) == 0 {
			break
		}
		loadToBigQuery(chunk, tableName)
		offset += limit
	}

	tableName = "a_records"
	offset = 0
	limit = 5000
	if err := recreateTable(types.ARecord{}, tableName); err != nil {
		log.Printf("Failed to truncate table: %v", err)
	}
	for {
		chunk := []types.ARecord{}
		err := db.GormDB.Limit(limit).Offset(offset).Find(&chunk).Error
		if err != nil {
			log.Fatalf("Failed to get chunk: %v", err)
		}
		if len(chunk) == 0 {
			break
		}
		loadToBigQuery(chunk, tableName)
		offset += limit

	}

	tableName = "aaaa_records"
	offset = 0
	limit = 1000
	if err := recreateTable(types.AAAARecord{}, tableName); err != nil {
		log.Fatalf("Failed to truncate table: %v", err)
	}
	for {
		chunk := []types.AAAARecord{}
		err := db.GormDB.Limit(limit).Offset(offset).Find(&chunk).Error
		if err != nil {
			log.Fatalf("Failed to get chunk: %v", err)
		}
		if len(chunk) == 0 {
			break
		}
		loadToBigQuery(chunk, tableName)
		offset += limit
	}

	tableName = "soa_records"
	offset = 0
	limit = 1000
	if err := recreateTable(types.SOARecord{}, tableName); err != nil {
		log.Fatalf("Failed to truncate table: %v", err)
	}
	for {
		chunk := []types.SOARecord{}
		err := db.GormDB.Limit(limit).Offset(offset).Find(&chunk).Error
		if err != nil {
			log.Fatalf("Failed to get chunk: %v", err)
		}
		if len(chunk) == 0 {
			break
		}
		loadToBigQuery(chunk, tableName)
		offset += limit

	}

	tableName = "web_redirect_domains"
	offset = 0
	limit = 1000
	if err := recreateTable(types.WebRedirectDomain{}, tableName); err != nil {
		log.Printf("Failed to truncate table: %v", err)
	}
	for {
		chunk := []types.WebRedirectDomain{}
		err := db.GormDB.Limit(limit).Offset(offset).Find(&chunk).Error
		if err != nil {
			log.Fatalf("Failed to get chunk: %v", err)
		}
		if len(chunk) == 0 {
			break
		}
		loadToBigQuery(chunk, tableName)
		offset += limit

	}

	tableName = "sitemaps"
	offset = 0
	limit = 1000
	if err := recreateTable(types.Sitemap{}, tableName); err != nil {
		log.Fatalf("Failed to truncate table: %v", err)
	}
	for {
		chunk := []types.Sitemap{}
		err := db.GormDB.Limit(limit).Offset(offset).Find(&chunk).Error
		if err != nil {
			log.Fatalf("Failed to get chunk: %v", err)
		}
		if len(chunk) == 0 {
			break
		}
		loadToBigQuery(chunk, tableName)
		offset += limit
	}

	tableName = "sitemap_web_domains"
	offset = 0
	limit = 5000
	if err := recreateTable(types.SitemapWebDomain{}, tableName); err != nil {
		log.Fatalf("Failed to truncate table: %v", err)
	}
	for {
		chunk := []types.SitemapWebDomain{}
		err := db.GormDB.Limit(limit).Offset(offset).Find(&chunk).Error
		if err != nil {
			log.Fatalf("Failed to get chunk: %v", err)
		}
		if len(chunk) == 0 {
			break
		}
		loadToBigQuery(chunk, tableName)
		offset += limit

	}

	tableName = "sitemap_contact_domains"
	offset = 0
	limit = 5000
	if err := recreateTable(types.SitemapContactDomain{}, tableName); err != nil {
		log.Fatalf("Failed to truncate table: %v", err)
	}
	for {
		chunk := []types.SitemapContactDomain{}
		err := db.GormDB.Limit(limit).Offset(offset).Find(&chunk).Error
		if err != nil {
			log.Fatalf("Failed to get chunk: %v", err)
		}
		if len(chunk) == 0 {
			break
		}
		loadToBigQuery(chunk, tableName)
		offset += limit
	}
}
func loadToBigQuery(model interface{}, tableName string) {

	ctx := context.Background()
	table := dataset.Table(tableName)
	// Ensure table exists
	for {
		if _, err := table.Metadata(ctx); err != nil {
			log.Printf("Table %s not found, sleeping for 3 seconds\n", tableName)
			time.Sleep(3 * time.Second)
		} else {
			break
		}
	}
	inserter := table.Inserter()
	if err := inserter.Put(ctx, model); err != nil {
		log.Printf("Failed to insert data into BigQuery table: %v\n", err)
		return
	}

	fmt.Printf("Data loaded into BigQuery table: %s.%s.%s\n", "unum-marketing-data-assets", "domwalk", tableName)
	time.Sleep(1 * time.Second)
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
	time.Sleep(10 * time.Second)
	return nil
}
