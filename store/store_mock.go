package store

import (
	"context"
	"errors"

	"github.com/yoshiakiley/katana/core"
)

var MockExpectListError = errors.New("MockExpectListError")

type MockStore[R any] struct{}

// Create Object
func (m MockStore[R]) Create(context.Context, R, map[string]any) error { return nil }

func (m MockStore[R]) Update(ctx context.Context, new R, q map[string]any) (R, error) {
	return new, nil
}

func (m MockStore[R]) Count(ctx context.Context, q map[string]any) (int64, error) {
	return 0, nil
}

func (m MockStore[R]) Delete(context.Context, map[string]any) error { return MockExpectListError }

func (m MockStore[R]) List(context.Context, map[string]any) ([]R, error) {
	return nil, MockExpectListError
}

func (m MockStore[R]) Get(context.Context, map[string]any) (R, error) {
	var r R
	return r, nil
}

func (m MockStore[R]) Watch(context.Context, string, map[string]any) (<-chan core.Event, <-chan error) {
	return nil, nil
}
