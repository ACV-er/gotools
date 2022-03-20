package lib

import (
	"encoding/json"
	"errors"
	"reflect"
	"sort"
	"strconv"
)

type Stringer interface {
	String() string
}

// 将任意类型转换为string，对象会先尝试运行String(),如果没有，会尝试变成json
func AnyToString(src interface{}) (ret string) {
	if src == nil {
		return
	}

	switch v := src.(type) {
	case string:
		ret = v
	case int:
		ret = strconv.Itoa(v)
	case int64:
		ret = strconv.FormatInt(v, 10)
	case float64:
		ret = strconv.FormatFloat(v, 'f', -1, 64)
	case bool:
		ret = strconv.FormatBool(v)
	case Stringer:
		ret = v.String()
	default:
		if reflect.TypeOf(src).Kind() == reflect.Ptr {
			ret = AnyToString(reflect.ValueOf(src).Elem().Interface())
		} else {
			ret_byte, _ := json.Marshal(src)
			ret = string(ret_byte)
		}
	}
	return
}

// 获取结构体tag，json部分
func getJsonTag(field reflect.StructField) (tag string) {
	tag = field.Tag.Get("json")
	if tag == "" {
		tag = field.Name
	}
	return
}

// 将对象或map生成map[string]interface{}
func GenMapFromObject(src interface{}) (ret map[string]interface{}, err error) {
	ret = make(map[string]interface{})
	v := reflect.ValueOf(src)
	for v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	if v.Kind() == reflect.Struct {
		t := v.Type()
		for i := 0; i < v.NumField(); i++ {
			field := t.Field(i)
			// 避免访问私有变量
			if field.PkgPath != "" {
				continue
			}

			value_tmp := v.Field(i)
			for value_tmp.Kind() == reflect.Ptr {
				value_tmp = value_tmp.Elem()
			}

			ret[getJsonTag(field)] = value_tmp.Interface()
		}
	} else if v.Kind() == reflect.Map {
		for _, k := range v.MapKeys() {
			ret[k.String()] = v.MapIndex(k).Interface()
		}
	} else {
		err = errors.New("genMapFromObject: src is not a struct or map")
	}

	return
}

type SortMapNode struct {
	Key   string
	Value interface{}
}

type SortMap []SortMapNode

func GenArrFromObjectSortByFieldName(src interface{}, key_cmp func(string, string) bool) (ret SortMap, err error) {
	ret = make(SortMap, 0)
	ret_map, err := GenMapFromObject(src)
	if err != nil {
		return
	}

	for k, v := range ret_map {
		ret = append(ret, SortMapNode{k, v})
	}

	sort.SliceStable(ret, func(i, j int) bool {
		return key_cmp(ret[i].Key, ret[j].Key)
	})

	return
}

func GenArrFromObjectOrderByFieldNameAsc(src interface{}) (ret SortMap, err error) {
	ret, err = GenArrFromObjectSortByFieldName(src, func(a, b string) bool {
		return a < b
	})

	return
}
