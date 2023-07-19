package mongo

import (
	"context"
	"fmt"
	"github.com/yoshiakiley/katana/common"
	"github.com/yoshiakiley/katana/core"
	"github.com/yoshiakiley/katana/dict"
	"github.com/yoshiakiley/katana/store"
	"github.com/yoshiakiley/katana/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"reflect"
	"strconv"
	"time"
)

var nilfilter = map[string]any{}

func (m *MongoCli[R]) Create(ctx context.Context, r R, q map[string]any) error {
	query := parseQ(q)
	nObj := core.NewObject(r, marshaler)
	nMap, err := nObj.ToMap2(true)

	if nMap[common.Uid].(string) == "" {
		nMap[common.Uid] = utils.GetUUID()
	}
	if err != nil {
		return err
	}
	_, err = m.cli.
		Database(query.DB).
		Collection(query.Coll).
		InsertOne(ctx, nMap)
	if err != nil {
		return err
	}
	return nil
}

func (m *MongoCli[R]) Update(ctx context.Context, new R, q map[string]any) (R, error) {
	query := parseQ(q)
	singleResult := m.cli.Database(query.DB).
		Collection(query.Coll).
		FindOne(ctx, query.Q)

	if singleResult.Err() == mongo.ErrNoDocuments {
		nObj := core.NewObject(new, marshaler)
		nMap, err := nObj.ToMap2(true)
		if err != nil {
			return new, err
		}
		if nMap[common.Uid] == "" {
			nMap[common.Uid] = utils.GetUUID()
		}
		_, err = m.cli.Database(query.DB).
			Collection(query.Coll).
			InsertOne(ctx, nMap)
		if err != nil {
			return new, err
		}
		return m.GetByQuery(ctx, query)
	}

	nObj := core.NewObject(new, marshaler)
	nMap, err := nObj.ToMap2(true)
	if err != nil {
		return new, err
	}

	updateMap := map[string]any{}
	for _, key := range query.MergeFields {
		updateMap[key] = nMap[key]
	}
	updateMap[common.Version] = utils.GetVersion()

	update := bson.M{"$set": updateMap}

	_, err = m.cli.Database(query.DB).
		Collection(query.Coll).
		UpdateOne(ctx,
			query.Q,
			update,
		)
	if err != nil {
		return new, err
	}
	return m.GetByQuery(ctx, query)
}

func (m *MongoCli[R]) List(ctx context.Context, q map[string]any) ([]R, error) {
	var targets []R
	fOpts := options.Find()
	query := parseQ(q)

	if query.Skip != 0 {
		fOpts.SetSkip(int64(query.Skip))
	}
	if query.Limit != 0 {
		fOpts.SetLimit(int64(query.Limit))
	}

	project := bson.M{}
	for _, field := range query.Fields {
		project[field] = 1
	}
	fOpts.SetProjection(project)

	fOpts.SetSort(query.Sort)
	cursor, err := m.cli.Database(query.DB).
		Collection(query.Coll).
		Find(ctx, query.Q, fOpts)
	if err != nil {
		return nil, err
	}
	if err := cursor.All(ctx, &targets); err != nil {
		return nil, err
	}

	return targets, nil
}

func (m *MongoCli[R]) Get(ctx context.Context, q map[string]any) (R, error) {
	var t R
	query := parseQ(q)
	singleResult := m.cli.Database(query.DB).
		Collection(query.Coll).
		FindOne(ctx, query.Q)

	if err := singleResult.Decode(&t); err != nil {
		if err == mongo.ErrNoDocuments {
			return t, store.DataNotFound
		}
		return t, err
	}
	return t, nil
}

func (m *MongoCli[R]) GetByQuery(ctx context.Context, query *query) (R, error) {
	var t R
	singleResult := m.cli.Database(query.DB).
		Collection(query.Coll).
		FindOne(ctx, query.Q)

	if err := singleResult.Decode(&t); err != nil {
		if err == mongo.ErrNoDocuments {
			return t, store.DataNotFound
		}
		return t, err
	}
	return t, nil
}

func (m *MongoCli[R]) Count(ctx context.Context, q map[string]any) (int64, error) {
	query := parseQ(q)

	count, err := m.cli.Database(query.DB).Collection(query.Coll).CountDocuments(ctx, query.Q, options.Count())
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (m *MongoCli[R]) Delete(ctx context.Context, q map[string]any) error {
	query := parseQ(q)
	if _, err := m.cli.Database(query.DB).
		Collection(query.Coll).
		DeleteMany(ctx, query.Q); err != nil {
		return err
	}
	return nil
}

func fieldMatchFilter(opData map[string]any, key string, value interface{}) bool {
	return reflect.DeepEqual(dict.Get(opData, key), value)
}

func versionMatchFilter(opData map[string]any, ver uint64) bool {
	version, exist := opData["version"]
	if version == nil || !exist {
		return false
	}

	switch version.(type) {
	case uint64:
		if version.(uint64) <= ver {
			return false
		}
	case uint32:
		if uint64(version.(uint64)) <= ver {
			return false
		}
	case int64:
		if uint64(version.(int64)) <= ver {
			return false
		}
	case int32:
		if uint64(version.(int32)) <= ver {
			return false
		}
	case float64:
		if uint64(version.(float64)) <= ver {
			return false
		}
	case float32:
		if uint64(version.(float32)) <= ver {
			return false
		}
	case string:
		i, err := strconv.ParseInt(version.(string), 10, 64)
		if err != nil {
			fmt.Printf("version data wrong %t\r\n", version)
			return false
		}
		uintVersion := uint64(i)
		if uintVersion <= ver {
			return false
		}
	default:
		fmt.Printf("unknow version type %t\r\n", version)
		return false

	}
	return true
}

func (m *MongoCli[R]) Watch(ctx context.Context, kind string, q map[string]any) (<-chan core.Event, <-chan error) {
	errC := make(chan error)
	query := parseQ(q)

	ns := fmt.Sprintf("%s.%s", query.DB, query.Coll)
	directReadFilter := func(op *Op) bool {
		if op.IsDelete() {
			return true
		}
		if !versionMatchFilter(op.Data, query.Version) {
			return false
		}
		for k, v := range query.Q {
			if pass := fieldMatchFilter(op.Data, k, v); !pass {
				return false
			}
		}
		return true
	}
	oplogTailCtx := Start(m.cli, &Options{
		DirectReadNs:     []string{ns},
		ChangeStreamNs:   []string{ns},
		MaxAwaitTime:     10,
		DirectReadFilter: directReadFilter,
	})

	decode := func(t any, s map[string]any) error {
		bsData, err := bson.Marshal(s)
		if err != nil {
			return err
		}
		return bson.Unmarshal(bsData, t)
	}
	eventC := make(chan core.Event, 0)
	go func() {
		for {
			select {
			case <-ctx.Done():
				close(eventC)
				oplogTailCtx.Stop()
				return
			case <-oplogTailCtx.ErrC:
				close(eventC)
				return
			case op, ok := <-oplogTailCtx.OpC:
				if !ok {
					return
				}
				if op.Data == nil {
					op.Data = make(map[string]interface{})
				}
				var evtOp core.EventType
				switch {
				case op.IsInsert():
					evtOp = core.ADDED
				case op.IsUpdate():
					evtOp = core.MODIFIED
				case op.IsDelete():
					evtOp = core.DELETED
					op.Data["version"] = time.Now().Unix()
				}
				// result need add support for multiple fields
				//op.Data["UID"] = op.Id
				//op.Data["_id"] = op.Id
				op.Data["kind"] = kind
				var r R
				if err := decode(&r, op.Data); err != nil {
					errC <- err
					return
				}
				evt := core.Event{
					Type:   evtOp,
					Object: r,
				}
				eventC <- evt
			}
		}
	}()

	return eventC, errC
}

func (m *MongoCli[R]) checkExistAndCreate(ctx context.Context, db, collection string, enableIndex bool, uniqeKeys ...string) error {
	names, err := m.cli.Database(db).
		ListCollectionNames(ctx, nilfilter)
	if err != nil {
		return err
	}
	exist := false

	for _, name := range names {
		if name == collection {
			exist = true
		}
	}

	if !exist {
		if err := m.cli.Database(db).CreateCollection(ctx, collection); err != nil {
			return err
		}
	}

	if enableIndex {
		keys := bson.D{}
		for _, k := range uniqeKeys {
			keys = append(keys, bson.E{Key: k, Value: 1})
		}
		indexModel := mongo.IndexModel{
			Keys:    keys,
			Options: options.Index().SetUnique(true),
		}
		if _, err = m.cli.Database(db).
			Collection(collection).
			Indexes().
			CreateOne(ctx, indexModel); err != nil {
			return err
		}
	}

	return nil
}
