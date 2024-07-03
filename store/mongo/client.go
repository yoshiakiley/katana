package mongo

import (
	"context"
	"fmt"
	"github.com/yoshiakiley/katana/store"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var cli *mongo.Client

var marshaler = &BSONMarshaler{}

type BSONMarshaler struct{}

func (b BSONMarshaler) Marshal(v any) ([]byte, error) {
	return bson.Marshal(v)
}

func (b BSONMarshaler) Unmarshal(data []byte, v any) error {
	return bson.Unmarshal(data, v)
}

func InitMongoCli(ctx context.Context, uri string) error {
	var err error
	if cli == nil {
		cli, err = connect(ctx, uri)
		if err != nil {
			return fmt.Errorf("connect to mongo: %w", err)
		}
	}
	store.Store = "mongo"
	return nil
}
