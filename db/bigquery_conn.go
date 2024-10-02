package db

import (
	"context"
	"log"

	"cloud.google.com/go/bigquery"
)

var BQConn *bigquery.Client

func CreateBigQueryConn() error {
	ctx := context.Background()

	// Set your Google Cloud Project ID
	projectID := "unum-marketing-data-assets"

	// Initialize the BigQuery client
	client, err := bigquery.NewClient(ctx, projectID)
	if err != nil {
		log.Fatalf("bigquery.NewClient: %v", err)
		return err
	}

	BQConn = client

	return nil
}
