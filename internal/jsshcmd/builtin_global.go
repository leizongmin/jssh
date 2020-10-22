package jsshcmd

import (
	"github.com/leizongmin/go/typeutil"
	"github.com/leizongmin/jssh/internal/jsexecutor"
	"os"
	"time"
)

func JsFnExit(global typeutil.H) jsexecutor.JSFunction {
	return func(ctx *jsexecutor.JSContext, this jsexecutor.JSValue, args []jsexecutor.JSValue) jsexecutor.JSValue {
		if len(args) < 1 {
			os.Exit(0)
			return ctx.Int32(0)
		}
		if !args[0].IsNumber() {
			return ctx.ThrowTypeError("exit: first argument expected number type")
		}
		code := args[0].Int32()
		os.Exit(int(code))
		return ctx.Int32(code)
	}
}

func JsFnGet(global typeutil.H) jsexecutor.JSFunction {
	return func(ctx *jsexecutor.JSContext, this jsexecutor.JSValue, args []jsexecutor.JSValue) jsexecutor.JSValue {
		if len(args) < 1 {
			return ctx.ThrowSyntaxError("get: missing name")
		}
		if !args[0].IsString() {
			return ctx.ThrowTypeError("get: first argument expected string type")
		}
		name := args[0].String()

		return ctx.Globals().Get(name)
	}
}

func JsFnSet(global typeutil.H) jsexecutor.JSFunction {
	return func(ctx *jsexecutor.JSContext, this jsexecutor.JSValue, args []jsexecutor.JSValue) jsexecutor.JSValue {
		if len(args) < 1 {
			return ctx.ThrowSyntaxError("set: missing name")
		}
		if !args[0].IsString() {
			return ctx.ThrowTypeError("set: first argument expected string type")
		}
		name := args[0].String()

		if len(args) < 2 {
			return ctx.ThrowSyntaxError("set: missing value")
		}
		value := args[1]

		ctx.Globals().Set(name, value)

		return ctx.Bool(true)
	}
}

func JsFnSleep(global typeutil.H) jsexecutor.JSFunction {
	return func(ctx *jsexecutor.JSContext, this jsexecutor.JSValue, args []jsexecutor.JSValue) jsexecutor.JSValue {
		if len(args) < 1 {
			return ctx.ThrowSyntaxError("sleep: missing millisecond argument")
		}
		if !args[0].IsNumber() {
			return ctx.ThrowTypeError("sleep: first argument expected number type")
		}
		ret, err := jsexecutor.JSValueToAny(args[0])
		if err != nil {
			return ctx.ThrowError(err)
		}
		ms, ok := ret.(float64)
		if !ok {
			return ctx.ThrowTypeError("sleep: first argument expected number type")
		}

		time.Sleep(time.Millisecond * time.Duration(ms))
		return jsexecutor.AnyToJSValue(ctx, ms)
	}
}
