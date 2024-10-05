package main

import (
	"context"
	"fmt"
	"log"

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
	if err := truncateTable(tableName); err != nil {
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
	limit = 1000
	tableName = "cert_sans_domains"
	if err := truncateTable(tableName); err != nil {
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
	limit = 1000
	if err := truncateTable(tableName); err != nil {
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
	limit = 1000
	if err := truncateTable(tableName); err != nil {
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
	if err := truncateTable(tableName); err != nil {
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
	if err := truncateTable(tableName); err != nil {
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
	if err := truncateTable(tableName); err != nil {
		log.Fatalf("Failed to truncate table: %v", err)
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
	if err := truncateTable(tableName); err != nil {
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
	limit = 1000
	if err := truncateTable(tableName); err != nil {
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
	limit = 1000
	if err := truncateTable(tableName); err != nil {
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
	inserter := table.Inserter()
	if err := inserter.Put(ctx, model); err != nil {
		log.Printf("Failed to insert data into BigQuery table: %v\n", err)
		return
	}

	fmt.Printf("Data loaded into BigQuery table: %s.%s.%s\n", "unum-marketing-data-assets", "domwalk", tableName)

}

func truncateTable(tableName string) error {
	ctx := context.Background()
	q := db.BQConn.Query(
		fmt.Sprintf(
			"DELETE FROM `%s.%s.%s` WHERE true", "unum-marketing-data-assets", "domwalk", tableName,
		),
	)
	job, err := q.Run(ctx)
	if err != nil {
		log.Fatalf("Failed to run delete query: %v", err)
	}
	status, err := job.Wait(ctx)
	if err != nil {
		log.Printf("Failed to wait for delete job: %v", err)
		return err
	}
	if err := status.Err(); err != nil {
		log.Printf("Delete job failed: %v", err)
	}
	return nil
}
