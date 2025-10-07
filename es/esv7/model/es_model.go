package model

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"github.com/doubunv/common-pkg/commonTool"
	"github.com/doubunv/common-pkg/es/esv7"
	"github.com/doubunv/common-pkg/es/esv7/core"
	"github.com/elastic/go-elasticsearch/v7/esapi"
	"strings"
)

type EsModel struct {
	Ctx context.Context
	Db  *core.Es
}

func NewEsModel(ctx context.Context, Db *core.Es) *EsModel {
	return &EsModel{
		Ctx: ctx,
		Db:  Db,
	}
}

func (model *EsModel) GetDb() *core.Es {
	if model.Db != nil {
		return model.Db
	}

	return nil
}

func (model *EsModel) IndexAuto(index string) error {
	create, err := model.GetDb().Indices.Create(index)
	if err != nil {
		return err
	}
	if create.IsError() {
		return errors.New(create.String())
	}

	return nil
}

func (model *EsModel) InsertSchema(data interface{}) error {
	dataJson, err := json.Marshal(data)
	if err != nil {
		return err
	}

	var (
		indexName string
		idKey     string
	)

	if str, ok := data.(esv7.IndexTable); ok {
		indexName = str.IndexName()
	}
	if str, ok := data.(esv7.Schema); ok {
		idKey = str.GetId()
		if idKey == "" {
			idKey = commonTool.GenUUID().String()
		}
	}

	if indexName == "" || idKey == "" {
		return errors.New("Not IndexTable. ")
	}

	response, err := model.GetDb().Create(indexName, idKey, bytes.NewBuffer(dataJson))
	if err != nil {
		return err
	}

	if response.IsError() {
		return errors.New(response.String())
	}

	return nil
}

func (model *EsModel) FindOne(id string, res interface{}) error {
	var (
		indexName string
	)

	if str, ok := res.(esv7.IndexTable); ok {
		indexName = str.IndexName()
	}

	resEs, err := model.GetDb().Get(indexName, id)
	if err != nil {
		return err
	}

	if err := core.ParseGetResponse(model.Ctx, &res, resEs); err != nil {
		return err
	}

	return nil
}
func (model *EsModel) Delete(id string, res interface{}) error {
	var (
		indexName string
	)
	if str, ok := res.(esv7.IndexTable); ok {
		indexName = str.IndexName()
	}

	resEs, err := model.GetDb().Delete(indexName, id)
	if err != nil {
		return err
	}
	if resEs.IsError() {
		return errors.New(resEs.String())
	}

	return nil
}

func (model *EsModel) UpdateSchema(data interface{}) error {
	dataJson, err := json.Marshal(data)
	if err != nil {
		return err
	}

	var (
		indexName string
		idKey     string
	)

	if str, ok := data.(esv7.IndexTable); ok {
		indexName = str.IndexName()
	}
	if str, ok := data.(esv7.Schema); ok {
		idKey = str.GetId()
	}

	if indexName == "" || idKey == "" {
		return errors.New("Not IndexTable or Schema. ")
	}

	response, err := model.GetDb().Update(indexName, idKey, bytes.NewBuffer(dataJson))
	if err != nil {
		return err
	}

	if response.IsError() {
		return errors.New(response.String())
	}

	return nil
}
func (model *EsModel) Search(res interface{}, res2 interface{}, query map[string]interface{}) (es *esapi.Response, num int64, err error) {
	var (
		indexName string
	)

	if str, ok := res.(esv7.IndexTable); ok {
		indexName = str.IndexName()
	}

	querydata, err := json.Marshal(query)
	if err != nil {
		return nil, 0, err
	}
	resEs, err := model.GetDb().Search(
		model.GetDb().Search.WithIndex(indexName),
		model.GetDb().Search.WithBody(strings.NewReader(string(querydata))),
	)
	if err != nil {
		return nil, 0, err
	}
	total, err := core.ParseSearchResponse(context.Background(), res2, resEs)
	if err != nil {
		return nil, 0, err
	}

	return resEs, total, err
}
