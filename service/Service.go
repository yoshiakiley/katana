package service

import (
	"context"
	"github.com/yoshiakiley/katana/common"
	"github.com/yoshiakiley/katana/core"
	"github.com/yoshiakiley/katana/store"
)

type IService interface {
	Watch(ctx context.Context, kind string, q map[string]any) (<-chan core.Event, <-chan error)
}

type Service[R core.IObject] struct {
	Schema     string
	Collection string
	*store.Store[R]
}

func (s *Service[R]) query(qs ...map[string]any) map[string]any {
	p := make(map[string]any)
	for _, q := range qs {
		for k, v := range q {
			p[k] = v
		}
	}
	p[common.Schema] = s.Schema
	p[common.Collection] = s.Collection
	return p
}

func (s *Service[R]) Create(ctx context.Context, r R) error {
	return s.Store.Create(ctx, r, s.query())
}

func (s *Service[R]) Update(ctx context.Context, new R, q map[string]any) (R, error) {
	return s.Store.Update(ctx, new, s.query(q))
}

func (s *Service[R]) Delete(ctx context.Context, q map[string]any) error {
	return s.Store.Delete(ctx, s.query(q))
}

func (s *Service[R]) List(ctx context.Context, q map[string]any) ([]R, error) {
	return s.Store.List(ctx, s.query(q))
}

func (s *Service[R]) Get(ctx context.Context, q map[string]any) (R, error) {
	return s.Store.Get(ctx, s.query(q))
}

func (s *Service[R]) Count(ctx context.Context, q map[string]any) (int64, error) {
	return s.Store.Count(ctx, s.query(q))
}

func (s *Service[R]) Watch(ctx context.Context, kind string, q map[string]any) (<-chan core.Event, <-chan error) {
	return s.Store.Watch(ctx, kind, s.query(q))
}

func (s *Service[R]) Set(store *store.Store[R]) {
	s.Store = store
}
