package commonTool

import (
	"fmt"
	"reflect"
	"sort"
	"strings"
)

func StructToSignString(v interface{}, ignore ...string) string {
	val := reflect.ValueOf(v)
	typ := reflect.TypeOf(v)

	// 可忽略字段列表 → map
	ignoreMap := make(map[string]struct{})
	for _, k := range ignore {
		ignoreMap[k] = struct{}{}
	}

	kv := make([]string, 0)

	// 支持指针
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
		typ = typ.Elem()
	}

	for i := 0; i < val.NumField(); i++ {
		field := typ.Field(i)
		key := field.Tag.Get("json")
		if key == "" {
			key = field.Name
		}

		// 忽略字段
		if _, skip := ignoreMap[key]; skip {
			continue
		}

		value := val.Field(i)
		// 忽略零值
		if isZero(value) {
			continue
		}

		kv = append(kv, fmt.Sprintf("%s=%v", key, value.Interface()))
	}

	// 字典序排序
	sort.Strings(kv)

	// 拼接 a=1&b=2&c=3
	return strings.Join(kv, "&")
}

func isZero(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.String:
		return v.String() == ""
	case reflect.Bool:
		return !v.Bool()
	case reflect.Int, reflect.Int64, reflect.Int32, reflect.Int16, reflect.Int8:
		return v.Int() == 0
	case reflect.Uint, reflect.Uint32, reflect.Uint64:
		return v.Uint() == 0
	case reflect.Float64, reflect.Float32:
		return v.Float() == 0
	case reflect.Interface, reflect.Ptr:
		return v.IsNil()
	default:
		// 默认使用 DeepEqual zero value
		return reflect.DeepEqual(v.Interface(), reflect.Zero(v.Type()).Interface())
	}
}
