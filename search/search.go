package search

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"

	"github.com/elastic/go-elasticsearch/v7"
)

type Q map[string]any

type SearchDao interface {
	Id() string
	Index() string
	Body() (*bytes.Reader, error)
	GetMapping() string
}

type SearchResult struct {
	Status string
	Total  int
	Took   int
	Hits   []*hit
}

func (sr *SearchResult) AddHit(id string, s map[string]any, score float64) {
	sr.Hits = append(sr.Hits, &hit{ID: id, Source: s, Score: score})
}

type hit struct {
	ID     string
	Source map[string]any
	Score  float64
}

func Search(ctx context.Context, clt *elasticsearch.Client, dao SearchDao, q Q) (*SearchResult, error) {
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(q); err != nil {
		return nil, err
	}
	res, err := clt.Search(
		clt.Search.WithContext(ctx),
		clt.Search.WithIndex(dao.Index()),
		clt.Search.WithBody(&buf),
		clt.Search.WithTrackTotalHits(true),
	)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.IsError() {
		var e map[string]interface{}
		if err := json.NewDecoder(res.Body).Decode(&e); err != nil {
			return nil, err
		} else {
			// Print the response status and error information.
			return nil, fmt.Errorf(
				"[%s] %s: %s",
				res.Status(),
				e["error"].(map[string]interface{})["type"],
				e["error"].(map[string]interface{})["reason"],
			)
		}
	}
	var r map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		return nil, fmt.Errorf("error parsing the response body: %w", err)
	}

	result := SearchResult{
		Status: res.Status(),
		Total:  int(r["hits"].(map[string]interface{})["total"].(map[string]interface{})["value"].(float64)),
		Took:   int(r["took"].(float64)),
	}

	// Print the ID and document source for each hit.
	for _, hit := range r["hits"].(map[string]interface{})["hits"].([]interface{}) {
		hitMap := hit.(map[string]interface{})
		if hitMap["_score"] != nil {
			result.AddHit(hitMap["_id"].(string), hitMap["_source"].(map[string]interface{}), hitMap["_score"].(float64))
		} else {
			result.AddHit(hitMap["_id"].(string), hitMap["_source"].(map[string]interface{}), 0)
		}
	}
	return &result, nil
}
