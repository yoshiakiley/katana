package store

import (
	"context"
	"fmt"

	"github.com/yoshiakiley/katana/core"
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
	// Create Query map[string]any{"version":[]CompareType{{GTE:1},{LTE:100}} } == version >= 1 and version <=100
	// Create Object
	Create(context.Context, R, map[string]any) error

	Update(ctx context.Context, new R, q map[string]any) (R, error)

	Delete(context.Context, map[string]any) error

	List(context.Context, map[string]any) ([]R, error)

	Get(context.Context, map[string]any) (R, error)

	Count(context.Context, map[string]any) (int64, error)

	Watch(context.Context, string, map[string]any) (<-chan core.Event, <-chan error)
}

type IRDBMSStore[R core.IObject] interface {
}

type Store[R core.IObject] struct {
	IKVStore[R]
	IRDBMSStore[R]
}

func NewStore[R core.IObject](kv IKVStore[R], rdbms IRDBMSStore[R]) *Store[R] {
	return &Store[R]{kv, rdbms}
}

func Get[R core.IObject](store any) *Store[R] {
	return store.(*Store[R])
}
