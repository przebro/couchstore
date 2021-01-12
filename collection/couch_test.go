package collection

import (
	"context"
	"testing"

	"github.com/przebro/couchdb/connection"
	"github.com/przebro/couchdb/database"
	"github.com/przebro/databazaar/collection"
	tst "github.com/przebro/databazaar/collection/testing"
	"github.com/przebro/databazaar/selector"
)

type InvalidDocument struct {
	ID    int    `json:"_id,omitempty"`
	REV   string `json:"_rev,omitempty"`
	Title string `json:"title"`
}

var col *couchDBCollection = nil

var (
	singleRecord        = tst.TestDocument{Title: "Blade Runner", Score: 8.1, Year: 1982, Oscars: false}
	singleRecordInvalid = InvalidDocument{ID: 12345, Title: "Invalid Title"}
	singleRecordWithID  tst.TestDocument
	testCollection      []tst.TestDocument
	conn                *connection.Connection
)

const host = "127.0.0.1"
const port = 5300

func init() {

	builder := connection.NewBuilder()
	builder = builder.WithAddress(host, port).
		WithAuthentication("COOKIE", "admin", "notsecure")

	conn, _ = builder.Build(true)
	_, db, _ := database.CreateDatabase(context.Background(), "databazaar", conn.GetClient())

	col = &couchDBCollection{database: db}

	singleRecordWithID, testCollection = tst.GetSingleRecord("../data/testdata.json")
}

func TestMain(m *testing.M) {

	m.Run()
	database.DropDatabase(context.Background(), "databazaar", conn.GetClient())
}

func TestInsertSingle(t *testing.T) {

	r, err := col.Create(context.Background(), &singleRecord)
	if err != collection.ErrEmptyOrInvalidID {
		t.Error(err)
	}

	r, err = col.Create(context.Background(), &singleRecordWithID)
	if err != nil {
		t.Error(err)
	}

	if r.ID != singleRecordWithID.ID {
		t.Error("unexpected result")
	}

	r, err = col.Create(context.Background(), &singleRecordInvalid)
	if err != collection.ErrEmptyOrInvalidID {
		t.Error("Unexpected result")
	}
}

func TestInsertMany(t *testing.T) {

	tc := []interface{}{}

	for x := range testCollection {
		tc = append(tc, testCollection[x])
	}
	results, err := col.CreateMany(context.Background(), tc)
	if err != nil {
		t.Error("unexpected result:", err)
	}

	if len(results) != len(testCollection) {
		t.Error("Unexpected result:", len(results), "expected:", len(testCollection))
	}

}

func TestGet(t *testing.T) {

	rs := tst.TestDocument{}

	err := col.Get(context.Background(), singleRecordWithID.ID, &rs)

	if err != nil {
		t.Error("unexpected result")
	}

	if rs.ID != singleRecordWithID.ID {
		t.Error("unexpected result, expected:", singleRecordWithID.ID, "actual:", rs.ID)
	}

}

func TestSelect(t *testing.T) {

	result := []tst.TestDocument{}

	sel := selector.Gte("year", selector.Int(1986))

	r, err := col.Select(context.Background(), sel, nil)
	if err != nil {
		t.Error("unexpected result:", err)
	}

	r.All(context.Background(), &result)

	if len(result) != 5 {
		t.Error("unexpected result: collection len expected:", 5, "actual:", len(result))
	}

	r, err = col.Select(context.Background(), sel, nil)
	cnt := 0
	for r.Next(context.Background()) {
		doc := tst.TestDocument{}
		r.Decode(&doc)
		cnt++
	}
	if cnt != 5 {
		t.Error("unexpected result:", cnt)
	}

	err = r.Close()
	if err != nil {
		t.Error("unexpected result")
	}

}

func TestUpdate(t *testing.T) {

	singleRecord.Score = 7.3

	err := col.Update(context.Background(), &singleRecord)
	if err == nil {
		t.Error("unexpected result")
	}

	doc := tst.TestDocument{
		ID:     "movie_13",
		Oscars: true,
		Score:  7.9,
		Year:   1999,
		Title:  "The Matrix",
	}
	result, err := col.Create(context.Background(), &doc)
	if err != nil {
		t.Error("unexpected result:", err)
	}

	doc.ID = result.ID
	doc.REV = result.Revision
	doc.Score = 2.3
	err = col.Update(context.Background(), &doc)
	if err != nil {
		t.Error("unexpected result:", err)
	}
}

func TestGetQuerableCollection(t *testing.T) {

	_, err := col.AsQuerable()
	if err != nil {
		t.Error("Unexpected result:", err)
	}

}

func TestAll(t *testing.T) {

	r, err := col.All(context.Background())
	if err != nil {
		t.Error("unexpected result:", err)
	}

	doc := tst.TestDocument{}
	numdocs := 0

	for r.Next(context.Background()) {
		numdocs++
		r.Decode(&doc)
	}

	if numdocs == 0 {
		t.Error("unexpected result")
	}

}

func TestDelete(t *testing.T) {
	err := col.Delete(context.Background(), "single_record")
	if err != nil {
		t.Error("unexpected result:", err)
	}
}
