package db

import (
	"context"
	"fmt"
	"log"

	"cloud.google.com/go/bigquery"
)

var BQConn *bigquery.Client

func CreateBigQueryConn() error {
	ctx := context.Background()
	// Set your Google Cloud Project ID

	// Initialize the BigQuery client
	client, err := bigquery.NewClient(ctx, bigquery.DetectProjectID)
	fmt.Println(client.Project())
	if err != nil {
		log.Fatalf("bigquery.NewClient: %v", err)
		return err
	}

	BQConn = client

	return nil
}
