package search

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strconv"
	"time"

	"github.com/elastic/go-elasticsearch/v7"
	"github.com/tidwall/gjson"
)

type PageSearch interface {
	Search(ctx context.Context, page, size int) (*SearchResult, error)
}

func NewPageSearch(clt *elasticsearch.Client, dao SearchDao, q Q) PageSearch {
	return &pageSearchImpl{clt: clt, dao: dao, q: q}
}

type pageSearchImpl struct {
	clt *elasticsearch.Client
	dao SearchDao
	q   Q
}

const _MAX_SIZE = 1000
const _ES_SEARCH_LIMIT = 10000

func (ps *pageSearchImpl) Search(ctx context.Context, page, size int) (*SearchResult, error) {
	if size > _MAX_SIZE {
		return nil, errors.New("size must be less than 1000")
	}
	from := (page - 1) * size
	if from+size <= _ES_SEARCH_LIMIT {
		ps.q["from"] = from
		ps.q["size"] = size
		return Search(ctx, ps.clt, ps.dao, ps.q)
	}
	after, err := ps.findSkipId(ctx, ps.q, from)
	if err != nil {
		return nil, err
	}
	ps.q["size"] = size
	ps.q["search_after"] = after
	return Search(ctx, ps.clt, ps.dao, ps.q)

}

func (ps *pageSearchImpl) findSkipId(ctx context.Context, q Q, skip int) ([]any, error) {
	const scrollTime = time.Second * 5
	q["size"] = _ES_SEARCH_LIMIT
	if _, ok := q["sort"]; !ok {
		q["sort"] = []Q{
			{"_score": "desc"},
			{"_id": "asc"},
		}
	}
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(q); err != nil {
		return nil, err
	}
	res, err := ps.clt.Search(
		ps.clt.Search.WithContext(ctx),
		ps.clt.Search.WithIndex(ps.dao.Index()),
		ps.clt.Search.WithBody(&buf),
		ps.clt.Search.WithScroll(scrollTime),
	)
	if err != nil {
		return nil, err
	}
	json := read(res.Body)

	res.Body.Close()
	scrollID := gjson.Get(json, "_scroll_id").String()
	count := (skip/_ES_SEARCH_LIMIT - 1)
	skip = skip % _ES_SEARCH_LIMIT
	if skip > 0 {
		count++
	}

	for i := 0; i < count; i++ {
		if scrollID == "" {
			return nil, nil
		}
		res, err = ps.clt.Scroll(
			ps.clt.Scroll.WithScrollID(scrollID),
			ps.clt.Scroll.WithScroll(scrollTime),
		)
		if err != nil {
			return nil, err
		}
		json = read(res.Body)
		res.Body.Close()
		hits := gjson.Get(json, "hits.hits")
		if len(hits.Array()) < 1 {
			return nil, errors.New("no more data")
		}
	}

	var index string
	if skip == 0 {
		index = "9999"
	} else {
		index = strconv.Itoa(skip - 1)
	}
	arr := gjson.Get(json, fmt.Sprintf("hits.hits.%s.sort", index)).Array()
	result := make([]any, len(arr))
	for i, v := range arr {
		result[i] = v.Value()
	}
	return result, nil
}

func read(r io.Reader) string {
	var b bytes.Buffer
	b.ReadFrom(r)
	return b.String()
}
