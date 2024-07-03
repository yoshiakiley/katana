package katana

import (
	"context"
	"github.com/yoshiakiley/katana/store/mongo"
)

const StoreTypeMongo = "mongo"

func InitStore(ctx context.Context, initStore string, addr string) error {
	switch initStore {
	case StoreTypeMongo:
		return mongo.InitMongoCli(ctx, addr)
	}
	return nil
}
