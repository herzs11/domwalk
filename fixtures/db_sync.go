package main

import (
	"context"
	"fmt"
	"log"

	"domwalk/db"
	"domwalk/types"
	"gorm.io/gorm"
)

// TODO: Sync sqlite with bigquery

func getAllRecords() {
	// Get all records from sqlite
	doms := []types.Domain{}
	err := db.GormDB.Find(&doms).Error
	if err != nil {
		log.Fatalf("Failed to get rows: %v", err)
	}
	stmt := &gorm.Statement{DB: db.GormDB}
	err = stmt.Parse(&doms)
	if err != nil {
		log.Fatalf("Failed to parse model: %v", err)
	}
	fmt.Println(stmt.Schema.Table)

}

func main() {
	err := db.CreateBigQueryConn()
	if err != nil {
		log.Fatal(err)
	}
	getAllRecords()
	dom := []types.Domain{}
	db.GormDB.Find(&dom)
	var domsToUpload []types.DomainBQ
	for _, d := range dom {
		domsToUpload = append(domsToUpload, d.ToBQ())
	}
	loadToBigQuery(domsToUpload, "domains")

	csan := []types.CertSAN{}
	db.GormDB.Find(&csan)
	loadToBigQuery(csan, "cert_sans")

	mx := []types.MXRecord{}
	db.GormDB.Find(&mx)
	loadToBigQuery(mx, "mx_records")

	a := []types.ARecord{}
	db.GormDB.Find(&a)
	loadToBigQuery(a, "a_records")

	aaaa := []types.AAAARecord{}
	db.GormDB.Find(&aaaa)
	loadToBigQuery(aaaa, "aaaa_records")

	soa := []types.SOARecord{}
	db.GormDB.Find(&soa)
	loadToBigQuery(soa, "soa_records")

	web := []types.WebRedirect{}
	db.GormDB.Find(&web)
	loadToBigQuery(web, "web_redirects")

	sitemap := []types.Sitemap{}
	db.GormDB.Find(&sitemap)
	loadToBigQuery(sitemap, "sitemaps")

	sitemapWeb := []types.SitemapWebDomain{}
	db.GormDB.Find(&sitemapWeb)
	loadToBigQuery(sitemapWeb, "sitemap_web_domains")

	sitemapContact := []types.SitemapContactDomain{}
	db.GormDB.Find(&sitemapContact)
	loadToBigQuery(sitemapContact, "sitemap_contact_domains")

	defer db.BQConn.Close()
}

func loadToBigQuery(model interface{}, tableName string) {
	dataset := db.BQConn.Dataset("domwalk")
	if err := truncateTable(tableName); err != nil {
		log.Fatalf("Failed to truncate table: %v", err)
	}
	ctx := context.Background()
	table := dataset.Table(tableName)
	inserter := table.Inserter()
	if err := inserter.Put(ctx, model); err != nil {
		log.Printf("Failed to insert data into BigQuery table: %v\n", err)
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
		log.Fatalf("Failed to wait for delete job: %v", err)
	}
	if err := status.Err(); err != nil {
		log.Fatalf("Delete job failed: %v", err)
	}
	return nil
}
