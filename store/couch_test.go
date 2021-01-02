package store

import (
	"context"
	"testing"

	"github.com/przebro/databazaar/store"
)

var (
	couchSrv = "couchdb;127.0.0.1:5300"
)

func TestConnection(t *testing.T) {

	_, err := store.NewStore(couchSrv + "/?username=abc&password=nopassword&auth=COOKIE")

	if err == nil {
		t.Error("unxcpected result")
	}

	_, err = store.NewStore(couchSrv + "/?username=abc&password=nopassword&auth=BASIC")
	if err == nil {
		t.Error("unxcpected result")
	}

	_, err = store.NewStore(couchSrv + "/?username=admin&password=notsecure&auth=COOKIE")

	if err != nil {
		t.Error("unxcpected result:", err)
	}

	_, err = store.NewStore(couchSrv + "/?username=admin&password=notsecure")

	if err != nil {
		t.Error("unxcpected result:", err)
	}

	_, err = store.NewStore(couchSrv + "/?username=admin&password=notsecure&auth=INVALID")

	if err == nil {
		t.Error("unxcpected result:", err)
	}
}

func TestStatus(t *testing.T) {

	store, err := store.NewStore(couchSrv + "/?username=admin&password=notsecure&auth=COOKIE")
	if err != nil {
		t.Error("unxcpected result:", err)
	}

	status, err := store.Status(context.Background())
	if err != nil {
		t.Error("unexpected result:", err)
	}

	if status != `200 OK` {
		t.Error("unexpected result:", status)
	}

}

func TestCreateCollection(t *testing.T) {

	store, err := store.NewStore(couchSrv + "/?username=admin&password=notsecure&auth=COOKIE")
	if err != nil {
		t.Error("unxcpected result:", err)
	}

	_, err = store.CreateCollection(context.Background(), "databazaar")
	if err != nil {
		t.Log("unexpected result:", err)
	}
}

func TestGetCollection(t *testing.T) {

	store, err := store.NewStore(couchSrv + "/?username=admin&password=notsecure&auth=COOKIE")
	if err != nil {
		t.Error("unxcpected result:", err)
	}

	_, err = store.Collection(context.Background(), "databazaar")
	if err != nil {
		t.Error("unexpected result:", err)
	}
}
