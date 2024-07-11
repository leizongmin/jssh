package jsexecutor

import (
	"fmt"
	"reflect"

	"github.com/buke/quickjs-go"

	"github.com/leizongmin/jssh/internal/utils"
)

type JSRuntime = quickjs.Runtime                                             // JSRuntime类型
type JSContext = quickjs.Context                                             // JSContext类型
type JSValue = quickjs.Value                                                 // JSValue类型
type JSFunction = func(ctx *JSContext, this JSValue, args []JSValue) JSValue // JSFunction类型

// NewJSRuntime 创建新的JSRuntime实例
func NewJSRuntime() JSRuntime {
	return quickjs.NewRuntime()
}

// IsGoFunction 判断是否为Go的函数类型
func IsGoFunction(f interface{}) bool {
	return reflect.TypeOf(f).Kind().String() == "func"
}

// MergeMapToJSObject 将map类型的值合并到一个JSObject中
func MergeMapToJSObject(ctx *quickjs.Context, obj quickjs.Value, vars utils.H) quickjs.Value {
	for n, v := range vars {
		obj.Set(n, AnyToJSValue(ctx, v))
	}
	return obj
}

// EvalJS 执行JS代码并返回JSValue结果
func EvalJS(jsRuntime quickjs.Runtime, code string, vars utils.H) (quickjs.Value, error) {
	ctx := jsRuntime.NewContext()
	defer ctx.Close()
	MergeMapToJSObject(ctx, ctx.Globals(), vars)
	return ctx.Eval(code)
}

// EvalJSFile 执行JS文件并返回JSValue结果
func EvalJSFile(jsRuntime quickjs.Runtime, code string, filename string, vars utils.H) (quickjs.Value, error) {
	ctx := jsRuntime.NewContext()
	defer ctx.Close()
	MergeMapToJSObject(ctx, ctx.Globals(), vars)
	return ctx.Eval(code, quickjs.EvalFileName(filename))
}

// EvalJSAndGetResult 执行JS并返回并返回interface{}结果
func EvalJSAndGetResult(jsRuntime quickjs.Runtime, code string, vars utils.H) (interface{}, error) {
	ret, err := EvalJS(jsRuntime, code, vars)
	if err != nil {
		return nil, err
	}
	defer ret.Free()
	return JSValueToAny(ret)
}

// JSValueToAny 将JSValue转换为interface{}
func JSValueToAny(value quickjs.Value) (interface{}, error) {
	if value.IsString() {
		return value.String(), nil
	}
	if value.IsNumber() {
		return value.Float64(), nil
	}
	if value.IsBool() {
		return value.Bool(), nil
	}
	if value.IsNull() || value.IsUndefined() {
		return nil, nil
	}
	if value.IsBigInt() {
		return value.BigInt().Int64(), nil
	}
	if value.IsBigFloat() {
		v, _ := value.BigFloat().Float64()
		return v, nil
	}
	if value.IsError() || value.IsException() {
		return value.Error(), nil
	}
	if value.IsFunction() {
		return value.String(), nil
	}
	if value.IsArray() {
		size := int(value.Len())
		arr := make([]interface{}, 0)
		for i := 0; i < size; i++ {
			v, err := JSValueToAny(value.Get(utils.AnythingToString(i)))
			if err != nil {
				return nil, err
			}
			arr = append(arr, v)
		}
		return arr, nil
	}
	if value.IsObject() {
		if JSValueIsUint8Array(value) {
			return JSValueUint8ArrayToByteSlice(value)
		}
		props, err := value.PropertyNames()
		if err != nil {
			return nil, err
		}
		m := make(utils.H)
		for _, p := range props {
			m[p] = value.Get(p)
		}
		return m, nil
	}
	return nil, fmt.Errorf("unexpected JS value: %+v", value)
}

// JSValueIsUint8Array 判断是否为Uint8Array
func JSValueIsUint8Array(value JSValue) bool {
	constructor := value.Get("constructor")
	defer constructor.Free()
	if !constructor.IsFunction() {
		return false
	}
	name := constructor.Get("name")
	defer name.Free()
	return name.String() == "Uint8Array"
}

// JSValueUint8ArrayToByteSlice 将Uint8Array转换为[]byte
func JSValueUint8ArrayToByteSlice(value quickjs.Value) ([]byte, error) {
	size := int(value.Len())
	arr := make([]byte, 0)
	for i := 0; i < size; i++ {
		v := value.Get(utils.AnythingToString(i))
		arr = append(arr, byte(v.Uint32()))
	}
	return arr, nil
}

func mapToJSValue(ctx *quickjs.Context, m utils.H) quickjs.Value {
	obj := ctx.Object()
	for n, v := range m {
		obj.Set(n, AnyToJSValue(ctx, v))
	}
	return obj
}

// AnyToJSValue 将interface{}转换为JSValue
func AnyToJSValue(ctx *quickjs.Context, value interface{}) quickjs.Value {
	if value == nil {
		return ctx.Undefined()
	}
	v := reflect.ValueOf(value)
	vt := v.Type()
	switch vt.Kind().String() {
	case "map":
		if m, ok := value.(map[string]interface{}); ok {
			return mapToJSValue(ctx, m)
		}
		if m, ok := value.(map[interface{}]interface{}); ok {
			m2 := make(utils.H)
			for k, v := range m {
				m2[fmt.Sprintf("%v", k)] = v
			}
			return mapToJSValue(ctx, m2)
		}
		return ctx.ThrowTypeError("AnyToJSValue: unsupported map type: %+v", value)
	case "slice":
		return anySliceToJSValue(ctx, v, vt)
	case "array":
		return anySliceToJSValue(ctx, v, vt)
	case "string":
		return ctx.String(value.(string))
	case "bool":
		return ctx.Bool(value.(bool))
	case "byte":
		return ctx.Uint32(uint32(value.(byte)))
	case "int":
		return ctx.Int32(int32(value.(int)))
	case "int8":
		return ctx.Int32(int32(value.(int8)))
	case "int16":
		return ctx.Int32(int32(value.(int16)))
	case "int32":
		return ctx.Int32(value.(int32))
	case "int64":
		return ctx.Int64(value.(int64))
	case "uint":
		return ctx.Int64(int64(value.(uint)))
	case "uint8":
		return ctx.Int32(int32(value.(uint8)))
	case "uint16":
		return ctx.Int32(int32(value.(uint16)))
	case "uint32":
		return ctx.Int64(int64(value.(uint32)))
	case "uint64":
		return ctx.BigUint64(value.(uint64))
	case "float32":
		return ctx.Float64(float64(value.(float32)))
	case "float64":
		return ctx.Float64(value.(float64))
	case "func":
		return ctx.Function(value.(JSFunction))
	default:
		return ctx.Undefined()
	}
}

// 将slice转换为JSValue
func anySliceToJSValue(ctx *quickjs.Context, v reflect.Value, vt reflect.Type) quickjs.Value {
	arr := ctx.Array()
	size := v.Len()
	if size > 0 {
		if vt.String() == "[]byte" {
			bytes := v.Bytes()
			for i := 0; i < len(bytes); i++ {
				arr.Set(int64(i), ctx.Uint32(uint32(bytes[i])))
			}
		} else {
			for i := 0; i < size; i++ {
				arr.Set(int64(i), AnyToJSValue(ctx, v.Index(i).Interface()))
			}
		}
	}
	return arr.ToValue()
}
