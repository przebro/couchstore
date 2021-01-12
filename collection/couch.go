package collection

import (
	"context"
	"fmt"

	"github.com/przebro/couchdb/database"
	query "github.com/przebro/couchstore/query"

	"github.com/przebro/databazaar/collection"

	"github.com/przebro/databazaar/result"
	"github.com/przebro/databazaar/selector"
)

type couchDBCollection struct {
	database database.CouchDatabase
}

func Collection(database database.CouchDatabase) collection.DataCollection {
	return &couchDBCollection{database: database}
}

func (col *couchDBCollection) Create(ctx context.Context, document interface{}) (*result.BazaarResult, error) {

	var res *result.BazaarResult = &result.BazaarResult{}

	id, _, err := collection.RequiredFields(document)
	if err != nil {
		return nil, err
	}
	if id == "" {
		return nil, collection.ErrEmptyOrInvalidID
	}

	ir, err := col.database.Insert(ctx, document)

	if err == nil {
		err = ir.Decode(&res)
	}
	return res, err
}
func (col *couchDBCollection) CreateMany(ctx context.Context, docs []interface{}) ([]result.BazaarResult, error) {

	var res []result.BazaarResult = []result.BazaarResult{}
	ir, err := col.database.InsertMany(ctx, docs)

	if err == nil {
		err = ir.Decode(&res)
	}

	return res, err
}
func (col *couchDBCollection) Get(ctx context.Context, id string, result interface{}) error {

	ir, err := col.database.Get(ctx, id)

	if err == nil {
		err = ir.Decode(result)
	}

	return err
}
func (col *couchDBCollection) Select(ctx context.Context, s selector.Expr, fld selector.Fields) (collection.BazaarCursor, error) {

	qbuilder := query.NewBuilder()
	dataSelector := qbuilder.Build(s)
	crsr, err := col.database.Select(ctx, dataSelector, fld, map[database.FindOption]interface{}{database.OptionStat: true})
	if err != nil {
		return nil, err
	}

	return &couchCursor{crsr}, nil
}

func (col *couchDBCollection) All(ctx context.Context) (collection.BazaarCursor, error) {

	qbuilder := query.NewBuilder()
	dataSelector := qbuilder.Build(selector.Ne("_id", selector.String("")))
	fmt.Println(dataSelector)

	crsr, err := col.database.Select(ctx, dataSelector, nil, map[database.FindOption]interface{}{database.OptionStat: true})
	if err != nil {
		return nil, err
	}

	return &couchCursor{crsr}, nil
}

func (col *couchDBCollection) Update(ctx context.Context, doc interface{}) error {
	_, err := col.database.Update(ctx, doc)
	return err
}
func (col *couchDBCollection) Delete(ctx context.Context, id string) error {

	r, err := col.database.Revision(ctx, id)

	if err != nil {
		return err
	}

	revs := struct {
		Data []struct {
			Rev    string `json:"rev"`
			Status string `json:"status"`
		} `json:"_revs_info"`
	}{}

	err = r.Decode(&revs)

	revisions := []string{}
	for x := 0; x < len(revs.Data); x++ {
		revisions = append(revisions, revs.Data[x].Rev)
	}

	_, err = col.database.Purge(ctx, id, revisions)

	return err
}

func (col *couchDBCollection) AsQuerable() (collection.QuerableCollection, error) {
	return col, nil
}
