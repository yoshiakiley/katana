package katana

import (
	"context"
	"testing"
)

type A struct {
	BrandName string `json:"brand_name" bson:"brand_name"`
}

func BenchmarkServer_single(b *testing.B) {
	ctx, cancel := context.WithCancel(context.Background())
	storeUrl := "mongodb://127.0.0.1:27017/admin"
	err := InitStore(ctx, storeTypeMongo, storeUrl)
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = NewService[A]("abc321_product", "commodity")
	}

	cancel()
}

func BenchmarkServer(b *testing.B) {
	ctx, cancel := context.WithCancel(context.Background())
	storeUrl := "mongodb://127.0.0.1:27017/admin"
	err := InitStore(ctx, storeTypeMongo, storeUrl)
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = GetService[A]("abc321_product", "commodity", "a")
	}

	cancel()
}
