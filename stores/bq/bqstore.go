package bq

import (
	"context"
	"log"
	"net/http"
	"sync"
	"time"

	"cloud.google.com/go/bigquery"
	"domwalk/domains"
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

func (bq *BQStore) PutDomains(ctx context.Context, doms []domains.Domain) error {
	var dbq []DomainBQ
	for _, d := range doms {
		dbq = append(dbq, newDomainBQ(d))
	}
	return bq.Table.Inserter().Put(ctx, dbq)
}

func (bq *BQStore) GetDomains(ctx context.Context, query string) ([]domains.Domain, error) {
	var doms []domains.Domain
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
	return doms, nil
}
