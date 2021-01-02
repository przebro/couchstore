package collection

import (
	"context"

	"github.com/przebro/couchdb/response"
)

type couchCursor struct {
	crsr *response.CouchMultiResult
}

func (c *couchCursor) All(ctx context.Context, v interface{}) error {
	return c.crsr.All(ctx, v)
}
func (c *couchCursor) Next(ctx context.Context) bool {
	return c.crsr.Next(ctx)
}
func (c *couchCursor) Decode(v interface{}) error {
	return c.crsr.Decode(v)
}
func (c *couchCursor) Close() error {
	return c.crsr.Close(context.Background())
}
