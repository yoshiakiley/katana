package mongo

import (
	"context"
	"fmt"
	"time"

	"github.com/yoshiakiley/katana/core"
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

type MongoCli[R core.IObject] struct {
	cli *mongo.Client
}

func (m *MongoCli[R]) Close() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return m.cli.Disconnect(ctx)
}

func NewMongoCli[R core.IObject]() *MongoCli[R] {
	return &MongoCli[R]{
		cli: cli,
	}
}

func InitMongoCli(ctx context.Context, uri string) error {
	var err error
	if cli == nil {
		cli, err = connect(ctx, uri)
		if err != nil {
			return fmt.Errorf("connect to mongo: %w", err)
		}
	}
	return nil
}
