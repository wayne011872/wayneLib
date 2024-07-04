package search

import (
	"bytes"
	"context"
	"errors"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/esapi"
)

func Index(ctx context.Context, clt *elasticsearch.Client, dao SearchDao) error {
	body, err := dao.Body()
	if err != nil {
		return err
	}
	indexReqA := esapi.IndexRequest{
		Index:      dao.Index(),
		DocumentID: dao.Id(),
		Body:       body,
		Refresh:    "true",
	}
	resp, err := indexReqA.Do(ctx, clt)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return errors.New(getResponseBody(resp.Body))
	}
	return nil
}

func DeleteDocument(ctx context.Context, clt *elasticsearch.Client, dao SearchDao) error {
	deleteReq := esapi.DeleteRequest{
		Index:      dao.Index(),
		DocumentID: dao.Id(),
	}
	resp, err := deleteReq.Do(ctx, clt)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return errors.New(getResponseBody(resp.Body))
	}
	return nil
}

func UpdateDocument(ctx context.Context, clt *elasticsearch.Client, dao SearchDao) error {
	body, err := dao.Body()
	if err != nil {
		return err
	}
	buf := new(bytes.Buffer)
	buf.ReadFrom(body)
	bodyTpl := []byte(`{"doc": %s,"doc_as_upsert": true}`)
	newBody := bytes.Replace(
		bodyTpl, []byte("%s"), buf.Bytes(), 1)
	updateReq := esapi.UpdateRequest{
		Index:      dao.Index(),
		DocumentID: dao.Id(),
		Body:       bytes.NewReader(newBody),
	}
	resp, err := updateReq.Do(ctx, clt)
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return errors.New(getResponseBody(resp.Body))
	}
	return err
}

func IsDucumentExist(ctx context.Context, clt *elasticsearch.Client, dao SearchDao) (bool, error) {
	existsReq := esapi.ExistsRequest{
		Index:      dao.Index(),
		DocumentID: dao.Id(),
	}
	resp, err := existsReq.Do(ctx, clt)
	if err != nil {
		return false, err
	}
	return resp.StatusCode != http.StatusNotFound, nil
}

func getResponseBody(rc io.ReadCloser) string {
	body, _ := ioutil.ReadAll(rc)
	rc.Close()
	return string(body)
}
