package search

import (
	"context"
	"strings"

	"github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/esapi"
)

func CreateIndex(ctx context.Context, clt *elasticsearch.Client, dao SearchDao) error {
	indexReq := esapi.IndicesCreateRequest{
		Index: dao.Index(),
		Body:  strings.NewReader(dao.GetMapping()),
	}
	_, err := indexReq.Do(ctx, clt)
	if err != nil {
		return err
	}
	return nil
}
