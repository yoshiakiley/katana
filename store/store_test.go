package store

import (
	"context"
	"testing"
)

type R struct{}

func (R) Kind(string) {}
func (R) UUID(string) {}

func Test_Store_NewStore(t *testing.T) {
	type K string
	type Q map[string]any

	store := NewStore[R](&MockStore[R]{}, nil)
	if &store == nil {
		t.Failed()
	}

	ctx := context.Background()
	_, err := store.List(ctx, Q{"a": "b"})
	if err != MockExpectListError {
		t.Failed()
	}

	err = store.Delete(ctx, map[string]any{"a": []CompareType{{LTE, 100}, {GTE, 1}}})
	if err != MockExpectListError {
		t.Failed()
	}

	t.Logf("ok")
}
