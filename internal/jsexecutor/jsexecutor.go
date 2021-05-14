package jsexecutor

import (
	"fmt"
	"reflect"

	"github.com/leizongmin/go/textutil"
	"github.com/leizongmin/go/typeutil"

	"github.com/leizongmin/jssh/quickjs"
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
	if reflect.TypeOf(f).Kind().String() == "func" {
		return true
	}
	return false
}

// MergeMapToJSObject 将map类型的值合并到一个JSObject中
func MergeMapToJSObject(ctx *quickjs.Context, obj quickjs.Value, vars typeutil.H) quickjs.Value {
	for n, v := range vars {
		if IsGoFunction(v) {
			obj.SetFunction(n, v.(JSFunction))
		} else {
			obj.Set(n, AnyToJSValue(ctx, v))
		}
	}
	return obj
}

// EvalJS 执行JS代码并返回JSValue结果
func EvalJS(jsRuntime quickjs.Runtime, code string, vars typeutil.H) (quickjs.Value, error) {
	ctx := jsRuntime.NewContext()
	defer ctx.Free()
	MergeMapToJSObject(ctx, ctx.Globals(), vars)
	return ctx.Eval(code)
}

// EvalJSFile 执行JS文件并返回JSValue结果
func EvalJSFile(jsRuntime quickjs.Runtime, code string, filename string, vars typeutil.H) (quickjs.Value, error) {
	ctx := jsRuntime.NewContext()
	defer ctx.Free()
	MergeMapToJSObject(ctx, ctx.Globals(), vars)
	return ctx.EvalFile(code, filename)
}

// EvalJSAndGetResult 执行JS并返回并返回interface{}结果
func EvalJSAndGetResult(jsRuntime quickjs.Runtime, code string, vars typeutil.H) (interface{}, error) {
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
			v, err := JSValueToAny(value.Get(textutil.AnythingToString(i)))
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
		m := make(typeutil.H)
		for _, p := range props {
			if p.IsEnumerable {
				k := p.Atom.Value()
				k2 := p.Atom.Value().String()
				k.Free()
				v := value.Get(k2)
				v2, err := JSValueToAny(v)
				v.Free()
				if err != nil {
					return nil, err
				}
				m[p.Atom.String()] = v2
			}
		}
		return m, nil
	}
	return nil, fmt.Errorf("unexpected JS value: %+v", value)
}

// JSValueIsUint8Array 判断是否为Uint8Array
func JSValueIsUint8Array(value JSValue) bool {
	constructor := value.Get("constructor")
	defer constructor.Free()
	if !constructor.IsConstructor() {
		return false
	}
	name := constructor.Get("name")
	defer name.Free()
	if name.String() != "Uint8Array" {
		return false
	}
	return true
}

// JSValueUint8ArrayToByteSlice 将Uint8Array转换为[]byte
func JSValueUint8ArrayToByteSlice(value quickjs.Value) ([]byte, error) {
	size := int(value.Len())
	arr := make([]byte, 0)
	for i := 0; i < size; i++ {
		v := value.Get(textutil.AnythingToString(i))
		arr = append(arr, byte(v.Uint32()))
	}
	return arr, nil
}

func mapToJSValue(ctx *quickjs.Context, m typeutil.H) quickjs.Value {
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
			m2 := make(typeutil.H)
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
				arr.SetByUint32(uint32(i), ctx.Uint32(uint32(bytes[i])))
			}
		} else {
			for i := 0; i < size; i++ {
				arr.SetByUint32(uint32(i), AnyToJSValue(ctx, v.Index(i).Interface()))
			}
		}
	}
	return arr
}
