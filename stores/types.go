package stores

import (
	"context"
	"fmt"
)

type Storer interface {
	Put(context.Context, string, any) error
}

type BQObject interface {
	ToBQ() any
}

type BQStore struct {
}

func (bq *BQStore) Put(ctx context.Context, tableName string, obj any) error {
	bqObj, ok := obj.(BQObject)
	if !ok {
		return fmt.Errorf("object does not implement BQObject")
	}

	return nil
}
