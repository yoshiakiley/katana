package service

import "github.com/yoshiakiley/katana/core"

type manager map[string]any

var m = &manager{}

func Set[R core.IObject](key string, svc *Service[R]) {
	m.set(key, svc)
}

func Get[R core.IObject](key string) (*Service[R], bool) {
	svc, exists := m.get(key)
	if !exists {
		return nil, false
	}
	return svc.(*Service[R]), true
}

func (m manager) set(k string, v any) {
	m[k] = v
}

func (m manager) get(k string) (any, bool) {
	svc, exists := m[k]
	return svc, exists
}
