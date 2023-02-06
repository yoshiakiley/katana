package core

import (
	"testing"
)

type Action struct {
	Uid     string `json:"uid"`
	Version string `json:"version"`
	Id      string `json:"id"`
}

func Test_Object_Clone(t *testing.T) {
	object := NewObject(Action{}, &JSONMarshaler{})

	object.Set(Action{
		Uid:     "123",
		Version: "123",
		Id:      "123",
	})

	newObj, err := object.Clone()
	if err != nil {
		t.Fatalf("%s", err)
	}
	old := object.Get()
	new := newObj.Get()

	if old.Uid != new.Uid || old.Version != new.Version {
		t.Failed()
	}

	t.Logf("ok")
}
