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

	dom := []types.Domain{}
	if err := truncateTable("domains"); err != nil {
		log.Fatalf("Failed to truncate table: %v", err)
	}
	db.GormDB.Find(&dom)
	for _, d := range dom {
		loadToBigQuery(d.ToBQ(), "domains")
	}

	csan := []types.CertSansDomain{}
	if err := truncateTable("cert_sans_domains"); err != nil {
		log.Fatalf("Failed to truncate table: %v", err)
	}
	db.GormDB.Find(&csan)
	for _, c := range csan {
		loadToBigQuery(c, "cert_sans_domains")
	}

	mx := []types.MXRecord{}
	if err := truncateTable("mx_records"); err != nil {
		log.Fatalf("Failed to truncate table: %v", err)
	}
	db.GormDB.Find(&mx)
	for _, m := range mx {
		loadToBigQuery(m, "mx_records")
	}

	a := []types.ARecord{}
	if err := truncateTable("a_records"); err != nil {
		log.Fatalf("Failed to truncate table: %v", err)
	}
	db.GormDB.Find(&a)
	for _, ip := range a {
		loadToBigQuery(ip, "a_records")
	}

	aaaa := []types.AAAARecord{}
	if err := truncateTable("aaaa_records"); err != nil {
		log.Fatalf("Failed to truncate table: %v", err)
	}
	db.GormDB.Find(&aaaa)
	loadToBigQuery(aaaa, "aaaa_records")

	soa := []types.SOARecord{}
	if err := truncateTable("soa_records"); err != nil {
		log.Fatalf("Failed to truncate table: %v", err)
	}
	db.GormDB.Find(&soa)
	loadToBigQuery(soa, "soa_records")

	web := []types.WebRedirectDomain{}
	if err := truncateTable("web_redirects"); err != nil {
		log.Fatalf("Failed to truncate table: %v", err)
	}
	db.GormDB.Find(&web)
	loadToBigQuery(web, "web_redirects")

	sitemap := []types.Sitemap{}
	if err := truncateTable("sitemaps"); err != nil {
		log.Fatalf("Failed to truncate table: %v", err)
	}
	db.GormDB.Find(&sitemap)
	loadToBigQuery(sitemap, "sitemaps")

	sitemapWeb := []types.SitemapWebDomain{}
	if err := truncateTable("sitemap_web_domains"); err != nil {
		log.Fatalf("Failed to truncate table: %v", err)
	}
	db.GormDB.Find(&sitemapWeb)
	loadToBigQuery(sitemapWeb, "sitemap_web_domains")

	sitemapContact := []types.SitemapContactDomain{}
	if err := truncateTable("sitemap_contact_domains"); err != nil {
		log.Fatalf("Failed to truncate table: %v", err)
	}
	db.GormDB.Find(&sitemapContact)
	loadToBigQuery(sitemapContact, "sitemap_contact_domains")

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
		log.Fatalf("Failed to wait for delete job: %v", err)
	}
	if err := status.Err(); err != nil {
		log.Fatalf("Delete job failed: %v", err)
	}
	return nil
}
