package dict

import (
	"fmt"
	"reflect"
	"testing"

	"gopkg.in/yaml.v2"
)

func Test_shift(t *testing.T) {
	path := "a.b.c"
	prefix, remain := shift(path)
	if prefix != "a" || remain != "b.c" {
		t.Fatal("expected not equal")
	}
}

func Test_Delete(t *testing.T) {
	data := map[string]interface{}{
		"a": map[string]interface{}{
			"b": map[string]interface{}{
				"c": 123,
			},
		},
	}
	expected := map[string]interface{}{
		"a": map[string]interface{}{
			"b": map[string]interface{}{},
		},
	}
	Delete(data, "a.b.c")

	if !reflect.DeepEqual(data, expected) {
		t.Fatal("expected not equal")
	}
}

func Test_Set(t *testing.T) {
	data := map[string]interface{}{
		"a": map[string]interface{}{
			"b": map[string]interface{}{
				"c": 123,
			},
		},
	}
	expected := map[string]interface{}{
		"a": map[string]interface{}{
			"b": map[string]interface{}{
				"c": 456,
			},
		},
	}

	Set(data, "a.b.c", 456)

	if !reflect.DeepEqual(data, expected) {
		t.Fatal("expected not equal")
	}
}

func Test_Get(t *testing.T) {
	data := map[string]interface{}{
		"a": map[string]interface{}{
			"b": map[string]interface{}{
				"c": 123,
			},
		},
	}
	value := Get(data, "a.b.c")
	if value.(int) != 123 {
		t.Fatal("expected not equal")
	}
}

const deployment = `
`

func Test_Get_Deployment(t *testing.T) {
	object := make(map[string]interface{})
	err := yaml.Unmarshal([]byte(deployment), &object)
	if err != nil {
		t.Fatal(err)
	}
	v := Get(object, "spec.selector.matchLabels")
	if v == nil {
		t.Fatal("get expected data is nil")
	}
	fmt.Printf("%s\n", v)

}

var src = map[string]interface{}{
	"metadata": map[string]string{"A": "B"},
	"spec": map[string]interface{}{
		"ops": []string{"c", "r", "u", "d"},
		"a":   1,
	},
}

var dest = map[string]interface{}{
	"metadata": map[string]string{"A": "B"},
	"spec": map[string]interface{}{
		"ops": []string{"c", "r", "u", "d"},
		"a":   2,
		"b":   2,
		"c":   3,
	},
}

func Test_Compare(t *testing.T) {
	_, ok := CompareMergeObject(src, dest, "spec.a")
	if !ok {
		t.Fatal("expected not equal")
	}

	expected := map[string]interface{}{
		"metadata": map[string]string{"A": "B"},
		"spec": map[string]interface{}{
			"ops": []string{"c", "r", "u", "d"},
			"a":   2,
		},
	}

	if !reflect.DeepEqual(src, expected) {
		t.Fatalf("expected equal,%v", src)
	}
	_, ok = CompareMergeObject(src, dest, "spec")

	if !ok {
		t.Fatal("expected not equal")
	}

	if !reflect.DeepEqual(src, dest) {
		t.Fatalf("expected equal,%v", src)
	}
}
