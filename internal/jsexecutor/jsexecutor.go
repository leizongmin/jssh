package jsexecutor

import (
	"fmt"
	"github.com/leizongmin/go/textutil"
	"github.com/leizongmin/go/typeutil"
	"github.com/lithdew/quickjs"
	"reflect"
)

type JSRuntime = quickjs.Runtime
type JSContext = quickjs.Context
type JSValue = quickjs.Value
type JSFunction = func(ctx *JSContext, this JSValue, args []JSValue) JSValue

func NewJSRuntime() JSRuntime {
	return quickjs.NewRuntime()
}

func IsGoFunction(f interface{}) bool {
	if reflect.TypeOf(f).Kind().String() == "func" {
		return true
	}
	return false
}

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

func EvalJS(jsRuntime quickjs.Runtime, code string, vars typeutil.H) (quickjs.Value, error) {
	ctx := jsRuntime.NewContext()
	defer ctx.Free()
	MergeMapToJSObject(ctx, ctx.Globals(), vars)
	return ctx.Eval(code)
}

func EvalJSFile(jsRuntime quickjs.Runtime, code string, filename string, vars typeutil.H) (quickjs.Value, error) {
	ctx := jsRuntime.NewContext()
	defer ctx.Free()
	MergeMapToJSObject(ctx, ctx.Globals(), vars)
	return ctx.EvalFile(code, filename)
}

func EvalJSAndGetResult(jsRuntime quickjs.Runtime, code string, vars typeutil.H) (interface{}, error) {
	ret, err := EvalJS(jsRuntime, code, vars)
	if err != nil {
		return nil, err
	}
	defer ret.Free()
	return JSValueToAny(ret)
}

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

func MapToJSValue(ctx *quickjs.Context, m typeutil.H) quickjs.Value {
	obj := ctx.Object()
	for n, v := range m {
		obj.Set(n, AnyToJSValue(ctx, v))
	}
	return obj
}

func AnyToJSValue(ctx *quickjs.Context, value interface{}) quickjs.Value {
	v := reflect.ValueOf(value)
	vt := v.Type()
	switch vt.Kind().String() {
	case "map":
		return MapToJSValue(ctx, value.(map[string]interface{}))
	case "slice":
		{
			arr := ctx.Array()
			for i := 0; i < v.Len(); i++ {
				arr.SetByUint32(uint32(i), AnyToJSValue(ctx, v.Index(i).Interface()))
			}
			return arr
		}
	case "array":
		{
			arr := ctx.Array()
			for i := 0; i < v.Len(); i++ {
				arr.SetByUint32(uint32(i), AnyToJSValue(ctx, v.Index(i).Interface()))
			}
			return arr
		}
	case "string":
		return ctx.String(value.(string))
	case "bool":
		return ctx.Bool(value.(bool))
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
