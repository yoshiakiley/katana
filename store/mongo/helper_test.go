package mongo

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/yoshiakiley/katana/store"
	"go.mongodb.org/mongo-driver/bson"
)

func Test_parseQ(t *testing.T) {
	query := map[string]any{
		"version": store.CompareType{Op: store.GT, Value: 1},
	}
	rs := parseQ(query)
	if reflect.DeepEqual(rs.Q, bson.M{GT: 1}) {
		t.Errorf("parseQ() error: %v", rs.Q)
	}

	query2 := map[string]any{
		"version": []store.CompareType{{Op: store.GT, Value: 1}, {Op: store.LTE, Value: 9999}},
	}

	rs2 := parseQ(query2)
	if reflect.DeepEqual(rs2.Q, bson.M{GT: 1, LTE: 9999}) {
		t.Errorf("parseQ() error: %v", rs2.Q)
	}
}

func Test_parse(t *testing.T) {
	query := map[string]any{
		"version": store.CompareType{Op: store.OR, Value: 1},
	}
	rs := parseQ(query)
	fmt.Println(rs.Q)
}
