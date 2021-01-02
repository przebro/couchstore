package store

import (
	"context"
	"errors"

	couchdb "github.com/przebro/couchbazaar/collection"
	"github.com/przebro/couchdb/client"
	"github.com/przebro/couchdb/connection"
	"github.com/przebro/couchdb/database"

	"github.com/przebro/databazaar/store"

	o "github.com/przebro/databazaar/store"

	"github.com/przebro/databazaar/collection"
)

const (
	couchstore     = "couchdb"
	authType       = "auth"
	authTypeBasic  = "basic"
	authTypeCookie = "cookie"
)

func init() {
	store.RegisterStoreFactory(couchstore, initCouchDB)
}

type CouchStore struct {
	connection *connection.Connection
}

func initCouchDB(opt o.ConnectionOptions) (store.DataStore, error) {

	var atype client.AuthType
	if opt.Options[authType] == "" {
		atype = client.None
	} else {
		atype = client.AuthType(opt.Options[authType])

	}

	switch atype {
	case
		client.Basic,
		client.Cookie,
		client.None:
	default:
		return nil, errors.New("invalid authentication type")

	}

	builder := connection.NewBuilder()
	builder = builder.WithAddress(opt.Host, opt.Port).
		WithAuthentication(atype, opt.Options[o.UsernameOption], opt.Options[o.PasswordOption])

	connection, err := builder.Build(true)

	if err != nil {
		return nil, err
	}

	return &CouchStore{connection: connection}, nil
}

func (s *CouchStore) Collection(ctx context.Context, name string) (collection.DataCollection, error) {

	_, data, err := database.GetDatabsase(ctx, name, s.connection.GetClient())
	if err != nil {
		return nil, err
	}

	return couchdb.Collection(data), nil

}

func (s *CouchStore) CreateCollection(ctx context.Context, name string) (collection.DataCollection, error) {

	_, data, err := database.CreateDatabase(ctx, name, s.connection.GetClient())

	if err != nil {
		return nil, err
	}

	return couchdb.Collection(data), nil
}

//Status - Gets status of a connection
func (s *CouchStore) Status(ctx context.Context) (string, error) {

	r, err := s.connection.Up(ctx)

	if err != nil {
		return "", err
	}

	return r.Status, nil
}
func (s *CouchStore) Close(ctx context.Context) {
	s.Close(ctx)

}
