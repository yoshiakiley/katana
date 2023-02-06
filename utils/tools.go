package utils

import (
	"encoding/json"
	"github.com/google/uuid"
	"reflect"
	"strings"
	"time"
)

func UnstructuredObjectToInstanceObj(src interface{}, dst interface{}) error {
	data, err := json.Marshal(src)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, dst)
}

func GetVersion() uint64 {
	return uint64(time.Now().Unix())
}

func GetUUID() string {
	uuidStr := uuid.New().String()
	uuidList := strings.Split(uuidStr, "-")
	uuidStr = strings.Join(uuidList, "")
	return uuidStr
}

func ReverseSlice[T []float64](slice T) {
	size := len(slice)
	swap := reflect.Swapper(slice)
	for i, j := 0, size-1; i < j; i, j = i+1, j-1 {
		swap(i, j)
	}
}
