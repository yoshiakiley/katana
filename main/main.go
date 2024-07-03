package main

import (
	"context"
	"fmt"
	"github.com/yoshiakiley/katana"
	"github.com/yoshiakiley/katana/common"
	"github.com/yoshiakiley/katana/store/mongo"
)

func main() {
	ctx := context.TODO()
	err := katana.InitStore(ctx, katana.StoreTypeMongo, "mongodb://localhost:27017")
	if err != nil {
		panic(err)
	}

	query := map[string]any{
		common.Schema:     "algalon",
		common.Collection: "exchange_code",
	}

	data, err := mongo.List[map[string]any](ctx, query)
	if err != nil {
		panic(err)
	}
	fmt.Println(data)

}
