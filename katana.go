package katana

import (
	"context"
	"github.com/yoshiakiley/katana/common"
	"github.com/yoshiakiley/katana/store/mongo"
)

func InitStore(ctx context.Context, initStore common.StoreType, addr string) error {
	switch initStore {
	case common.StoreTypeMongo:
		return mongo.InitMongoCli(ctx, addr)
	}
	return nil
}
