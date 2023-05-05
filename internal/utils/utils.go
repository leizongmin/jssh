package utils

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"reflect"
	"time"
)

type H = map[string]interface{}

var source rand.Source
var seed int64
var r *rand.Rand

func init() {
	SetSeed(time.Now().UnixNano())
}

// 将任意内容转换为字符串
func AnythingToString(v interface{}) string {
	if v == nil {
		return ""
	}
	if v2, ok := v.(string); ok {
		return v2
	}
	if v2, ok := v.(map[string]interface{}); ok {
		ret, _ := json.Marshal(&v2)
		return string(ret)
	}
	if v2, ok := ToInterfaceArray(v); ok {
		ret, _ := json.Marshal(&v2)
		return string(ret)
	}
	return fmt.Sprintf("%+v", v)
}

func ToInterfaceArray(value interface{}) ([]interface{}, bool) {
	ret, ok := value.([]interface{})
	if ok {
		return ret, ok
	}
	if IsArrayOrSlice(value) {
		list := reflect.ValueOf(value)
		size := list.Len()
		newList := make([]interface{}, size)
		i := 0
		for i < size {
			newList[i] = list.Index(i).Interface()
			i++
		}
		return newList, true
	}
	return nil, false
}

func IsArrayOrSlice(value interface{}) bool {
	rt := reflect.TypeOf(value)
	switch rt.Kind() {
	case reflect.Slice:
		return true
	case reflect.Array:
		return true
	default:
		return false
	}
}

func IsArray(value interface{}) bool {
	rt := reflect.TypeOf(value)
	switch rt.Kind() {
	case reflect.Array:
		return true
	default:
		return false
	}
}

func IsSlice(value interface{}) bool {
	rt := reflect.TypeOf(value)
	switch rt.Kind() {
	case reflect.Slice:
		return true
	default:
		return false
	}
}

// 设置随机种子
func SetSeed(s int64) {
	seed = s
	source = rand.NewSource(seed)
	r = rand.New(source)
	rand.Seed(seed)
}

// 获取int63随机数
func Int63n(n int64) int64 {
	return r.Int63n(n)
}
