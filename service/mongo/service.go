package mongo

import (
	"github.com/yoshiakiley/katana/core"
	"github.com/yoshiakiley/katana/service"
	"github.com/yoshiakiley/katana/store"
	"github.com/yoshiakiley/katana/store/mongo"
)

func NewService[R core.IObject](schema, collection string) *service.Service[R] {
	svc := &service.Service[R]{
		Schema:     schema,
		Collection: collection,
	}
	svc.Set(store.NewStore[R](mongo.NewMongoCli[R](), nil))

	return svc
}
