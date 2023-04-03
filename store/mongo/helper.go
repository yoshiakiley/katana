package mongo

import (
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"reflect"
	"time"

	"github.com/yoshiakiley/katana/common"
	"github.com/yoshiakiley/katana/store"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsontype"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.opentelemetry.io/contrib/instrumentation/go.mongodb.org/mongo-driver/mongo/otelmongo"
)

func getCtx(client *mongo.Client) (context.Context, context.CancelFunc, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	if err := client.Connect(ctx); err != nil {
		return nil, cancel, err
	}
	return ctx, cancel, nil
}

func connect(ctx context.Context, uri string) (*mongo.Client, error) {
	cliOpt := options.Client()
	cliOpt.SetMaxPoolSize(2000)
	cliOpt.SetMinPoolSize(1)
	cliOpt.SetMaxConnIdleTime(time.Second)

	cliOpt.SetRegistry(
		bson.NewRegistryBuilder().
			RegisterTypeMapEntry(
				bsontype.DateTime,
				reflect.TypeOf(time.Time{})).
			Build(),
	)
	cliOpt.ApplyURI(uri).SetMonitor(otelmongo.NewMonitor())
	mcli, err := mongo.NewClient(cliOpt)
	if err != nil {
		return nil, err
	}
	ctx, cancel, err := getCtx(mcli)
	defer cancel()
	if err != nil {
		return nil, err
	}
	if err := mcli.Ping(ctx, nil); err != nil {
		return nil, err
	}
	return mcli, nil
}

type query struct {
	DB          string `json:"db"`
	Coll        string `json:"coll"`
	Version     uint64 `json:"version"`
	MergeFields []string
	Q           bson.M         `json:"q"`
	Limit       int            `json:"limit"`
	Skip        int            `json:"skip"`
	Sort        map[string]any `json:"sort"`
	Fields      []string       `json:"fields"`
}

const (
	LT  = "$lt"
	LTE = "$lte"
	GT  = "$gt"
	GTE = "$gte"
	ALL = "$all"
	OR  = "$or"
	Ne  = "$ne"
	IN  = "$in"
)

// TODO: 增加字段检查
func parseQ[Q map[string]any](q Q) *query {
	var db string
	if _, exist := q[common.Schema]; exist {
		db = q[common.Schema].(string)
	}

	var coll string
	if _, exists := q[common.Schema]; exists {
		coll = q[common.Collection].(string)
	}

	var version uint64
	if ver, exist := q[common.Version]; exist {
		_version, ok := ver.(uint64)
		if !ok {
			goto NEXT
		}
		delete(q, common.Version)
		version = _version
	}
NEXT:
	var fields []string
	var returnFields []string
	var limit int
	var skip int
	var sort map[string]any

	if p, exist := q[common.ReturnFields]; exist {
		returnFields = p.([]string)
		delete(q, common.ReturnFields)
	}

	if p, exist := q[common.MergeFields]; exist {
		fields = p.([]string)
		delete(q, common.MergeFields)
	}

	delete(q, common.Schema)
	delete(q, common.Collection)

	d := make(bson.M)
	if uid, exist := q[common.Uid]; exist {
		delete(q, common.Uid)
		switch uid.(type) {
		case []string:
			d[common.Uid] = bson.M{IN: uid}
		default:
			d[common.Uid] = uid
		}

	}

	if queryOr, exists := q[common.QueryOR]; exists {
		delete(q, common.QueryOR)
		or := make([]map[string]any, 0)
		for k, v := range queryOr.(map[string]store.ValueType) {
			if v.Value == "" {
				continue
			}
			if v.Regex {
				or = append(or, map[string]any{
					k: bson.M{"$regex": primitive.Regex{Pattern: ".*" + v.Value + ".*", Options: "i"}}})
				continue
			}
			or = append(or, map[string]any{
				k: v.Value,
			})
		}
		if len(or) != 0 {
			d[OR] = or
		}
	}

	if querySort, exists := q[common.Sort]; exists {
		delete(q, common.Sort)
		sort = querySort.(map[string]any)
	}

	for k, v := range q {
		switch v.(type) {
		case store.CompareType:
			compare := v.(store.CompareType)
			switch compare.Op {
			case store.LT:
				d[k] = bson.M{LT: compare.Value}
			case store.LTE:
				d[k] = bson.M{LTE: compare.Value}
			case store.GT:
				d[k] = bson.M{GT: compare.Value}
			case store.GTE:
				d[k] = bson.M{GTE: compare.Value}
			case store.ALL:
				d[k] = bson.M{ALL: compare.Value}
			case store.OR:
				d[k] = bson.M{OR: compare.Value}
			case store.Ne:
				d[k] = bson.M{Ne: compare.Value}

			}
		case []store.CompareType:
			condictions := make(bson.M)
			for _, compare := range v.([]store.CompareType) {
				switch compare.Op {
				case store.LT:
					condictions[LT] = compare.Value
				case store.LTE:
					condictions[LTE] = compare.Value
				case store.GT:
					condictions[GT] = compare.Value
				case store.GTE:
					condictions[GTE] = compare.Value
				case store.ALL:
					condictions[ALL] = compare.Value
				case store.OR:
					d[k] = bson.M{OR: compare.Value}
				}
			}
			d[k] = condictions
		case store.PaginationType:
			pagination := v.(store.PaginationType)
			limit = pagination.Limit
			skip = pagination.Skip

		case store.ValueType:
			vValue := v.(store.ValueType).Value
			d[k] = bson.M{"$regex": primitive.Regex{Pattern: ".*" + vValue + ".*", Options: "i"}}

		case string:
			if v != "" {
				d[k] = v
			}
		case []string, []int, []float64, []float32, []any:
			d[k] = bson.M{IN: v}
		default:
			d[k] = v
		}
	}
	return &query{
		DB:          db,
		Coll:        coll,
		MergeFields: fields,
		Version:     version,
		Q:           d,
		Limit:       limit,
		Skip:        skip,
		Sort:        sort,
		Fields:      returnFields,
	}
}

func example() {
	parseQ(map[string]any{"DB": "base"})
}
