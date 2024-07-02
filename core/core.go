package core

import (
	"encoding/json"
	"github.com/yoshiakiley/katana/utils"

	"github.com/yoshiakiley/katana/common"
	"github.com/yoshiakiley/katana/dict"
)

var defaultDiscardedKeys = []string{"kind"}

type Marshaler interface {
	Unmarshal([]byte, any) error
	Marshal(any) ([]byte, error)
}

type JSONMarshaler struct{}

func (j JSONMarshaler) Marshal(v any) ([]byte, error) {
	return json.Marshal(v)
}

func (j JSONMarshaler) Unmarshal(data []byte, v any) error {
	return json.Unmarshal(data, v)
}

type IObject interface {
}

type Object[T IObject] struct {
	Item T
	Marshaler
}

func NewObject[T IObject](item T, m Marshaler) Object[T] {
	if m == nil {
		m = &JSONMarshaler{}
	}
	return Object[T]{Item: item, Marshaler: m}
}

func (o *Object[T]) Set(item T) { o.Item = item }

func (o *Object[T]) Get() T { return o.Item }

func (o *Object[T]) CompareMergeObject(other Object[T], paths ...string) (map[string]any, bool, error) {
	self, err := o.ToMap()
	if err != nil {
		return nil, false, err
	}
	_other, err := other.ToMap()
	if err != nil {
		return nil, false, err
	}

	for _, k := range defaultDiscardedKeys {
		dict.Delete(self, k)
	}

	updateMap, isUpdate := dict.CompareMergeObject(self, _other, paths...)
	return updateMap, isUpdate, nil
}

func (o *Object[T]) Unmarshal(a []byte) error {
	var item = new(T)
	if err := o.Marshaler.Unmarshal(a, item); err != nil {
		return err
	}
	o.Item = *item
	return nil
}

func (o *Object[T]) From(i any) error {
	b, err := o.Marshaler.Marshal(i)
	if err != nil {
		return err
	}
	return o.Unmarshal(b)
}

func (o *Object[T]) ToMap2(generate bool, discardedKeys ...string) (map[string]any, error) {
	r, err := o.ToMap(discardedKeys...)
	if err != nil {
		return nil, err
	}

	if generate {
		r[common.Version] = utils.GetVersion()
	}
	return r, nil
}

func (o *Object[T]) ToMap(discardedKeys ...string) (map[string]any, error) {
	bs, err := o.Marshaler.Marshal(&o.Item)
	if err != nil {
		return nil, err
	}
	var r map[string]any
	if err = o.Marshaler.Unmarshal(bs, &r); err != nil {
		return nil, err
	}
	for _, k := range discardedKeys {
		dict.Delete(r, k)
	}
	return r, nil
}

func (o *Object[T]) Clone() (*Object[T], error) {
	var obj Object[T]
	var target T
	src, err := o.Marshaler.Marshal(&o.Item)
	if err != nil {
		return nil, err
	}
	if err := o.Marshaler.Unmarshal(src, &target); err != nil {
		return nil, err
	}
	obj.Item = target
	return &obj, nil
}

func Decode(src, target any) error {
	switch target.(type) {
	case []byte:
		return json.Unmarshal(target.([]byte), &src)
	case string:
		return json.Unmarshal([]byte(target.(string)), &src)
	default:
		bs, err := json.Marshal(target)
		if err != nil {
			return err
		}
		return json.Unmarshal(bs, &src)
	}
}

type EventType = string

const (
	ADDED    EventType = "ADDED"
	MODIFIED EventType = "MODIFIED"
	DELETED  EventType = "DELETED"
	REMOVED  EventType = "REMOVED"
)

type Event struct {
	Type   EventType `json:"type"`
	Object any       `json:"object"`
}
