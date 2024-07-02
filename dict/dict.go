package dict

import (
	"reflect"
	"strings"
)

// (src, dest,["spec.userId","spec.userName"])
func CompareMergeObject(src, dest map[string]any, paths ...string) (map[string]any, bool) {
	isUpdate := false
	updateMap := map[string]any{}
	for _, p := range paths {
		srcContent := Get(src, p)
		desContent := Get(dest, p)
		if reflect.DeepEqual(srcContent, desContent) {
			continue
		}
		updateMap[p] = desContent
		isUpdate = true
	}
	return updateMap, isUpdate
}

// Set "path":"a.b.c"
// data = {"a":{"b":{"c":123}}}
// Set(data,"a.b.c",123)
func Set(data map[string]any, path string, value any) {
	head, remain := shift(path)
	_, exist := data[head]
	if !exist {
		data[head] = make(map[string]any)
	}
	if remain == "" {
		data[head] = value
		return
	}
	Set(data[head].(map[string]any), remain, value)
}

// Get data = {"a":{"b":{"c":123}}}
// Get(data,"a.b.c") = 123
func Get(data map[string]any, path string) (value any) {
	head, remain := shift(path)
	_, exist := data[head]
	if exist {
		if remain == "" {
			return data[head]
		}
		switch data[head].(type) {
		case map[string]any:
			return Get(data[head].(map[string]any), remain)
		case map[any]any:
			_data := make(map[string]any)
			for k, v := range data[head].(map[any]any) {
				_data[k.(string)] = v
			}
			return Get(_data, remain)
		}
	}
	return nil
}

// Delete data = {"a":{"b":{"c":123}}}
// Delete(data,"a.b.c") = {"a":{"b":""}}
func Delete(data map[string]any, path string) {
	head, remain := shift(path)
	_, exist := data[head]
	if exist {
		if remain == "" {
			delete(data, head)
			return
		}
		switch data[head].(type) {
		case map[string]any:
			Delete(data[head].(map[string]any), remain)
			return
		}
	}
	return
}

func shift(path string) (head string, remain string) {
	slice := strings.Split(path, ".")
	if len(slice) < 1 {
		return "", ""
	}
	if len(slice) < 2 {
		remain = ""
		head = slice[0]
		return
	}
	return slice[0], strings.Join(slice[1:], ".")
}
