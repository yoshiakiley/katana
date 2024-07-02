package store

import (
	"context"
	"fmt"
	"github.com/yoshiakiley/katana/core"
	"github.com/yoshiakiley/katana/store/mongo"
)

type QueryOpType uint

type CompareType struct {
	Op    QueryOpType
	Value any
}

type ValueType struct {
	Value string
	Regex bool
}

type PaginationType struct {
	Limit int
	Skip  int
}

const (
	LT QueryOpType = iota
	LTE
	GT
	GTE
	ALL
	OR
	Ne
)

var DataNotFound = fmt.Errorf("dataNotFound")

// IKVStore IStore interface method for internal implementation, only KV storage or rdbms storage is implemented here
type IKVStore[R core.IObject] interface {
	Create(context.Context, R, map[string]any) error
	//
	//Update(ctx context.Context, new R, q map[string]any) (R, error)
	//
	//Delete(context.Context, map[string]any) error
	//
	//List(context.Context, map[string]any) ([]R, error)
	//
	//Get(context.Context, map[string]any) (R, error)
	//
	//Count(context.Context, map[string]any) (int64, error)
	//
	//Watch(context.Context, string, map[string]any) (<-chan core.Event, <-chan error)
}

type Store struct {
	*mongo.MongoCli
}

func NewStore(*mongo.MongoCli) *Store {
	return &Store{}
}
