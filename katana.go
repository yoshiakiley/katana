package katana

import (
	"context"
	"github.com/yoshiakiley/katana/core"
	"github.com/yoshiakiley/katana/service"
	"github.com/yoshiakiley/katana/store"
	"github.com/yoshiakiley/katana/store/mongo"
)

type StoreType string

const StoreTypeMongo StoreType = "mongo"

func NewMongoService[R core.IObject](schema, collection string) *service.Service[R] {
	svc := &service.Service[R]{
		Schema:     schema,
		Collection: collection,
	}
	svc.Set(store.NewStore[R](mongo.NewMongoCli[R](), nil))

	return svc
}

func InitStore(ctx context.Context, initStore StoreType, addr string) error {
	switch initStore {
	case StoreTypeMongo:
		return mongo.InitMongoCli(ctx, addr)
	}
	return nil
}
