package core

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/doubunv/common-pkg/es/esaws"
	"github.com/elastic/go-elasticsearch/v7/esapi"
	"github.com/elastic/go-elasticsearch/v8/typedapi/core/get"
	"github.com/elastic/go-elasticsearch/v8/typedapi/core/search"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
	"io"
	"net/http"
)

var EsNotFound = errors.New("Es not found. ")

func ParseGetResponse(ctx context.Context, v interface{}, resEs *esapi.Response) error {
	if resEs.StatusCode == http.StatusNotFound {
		return EsNotFound
	}

	if resEs.IsError() {
		return errors.New(resEs.String())
	}
	defer resEs.Body.Close()

	body, err := io.ReadAll(resEs.Body)
	if err != nil {
		return err
	}

	var esBody get.Response
	if err := json.Unmarshal(body, &esBody); err != nil {
		return err
	}

	if err := json.Unmarshal(esBody.Source_, v); err != nil {
		return err
	}

	if item, ok := v.(esv7.Schema); ok {
		item.SetId(esBody.Id_)
	}

	return nil
}

func ParseSearchResponse(ctx context.Context, v interface{}, resEs *esapi.Response) (total int64, aggregate map[string]types.Aggregate, err error) {
	if resEs.StatusCode == http.StatusNotFound {
		return 0, nil, EsNotFound
	}

	if resEs.IsError() {
		return 0, nil, errors.New(resEs.String())
	}
	defer resEs.Body.Close()

	body, err := io.ReadAll(resEs.Body)
	if err != nil {
		return 0, nil, err
	}

	var esBody search.Response
	if err := json.Unmarshal(body, &esBody); err != nil {
		return 0, nil, err
	}

	var hits []map[string]interface{}
	for _, hit := range esBody.Hits.Hits {
		var temp map[string]interface{}
		err := json.Unmarshal(hit.Source_, &temp)
		if err != nil {
			return 0, nil, err
		}
		hits = append(hits, temp)

	}
	jsonHit, err := json.Marshal(hits)
	if err != nil {
		return 0, nil, err
	}
	err = json.Unmarshal(jsonHit, v)
	if err != nil {
		return 0, nil, err
	}

	return esBody.Hits.Total.Value, esBody.Aggregations, nil
}
