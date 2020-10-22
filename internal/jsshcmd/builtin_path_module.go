package jsshcmd

import (
	"github.com/leizongmin/go/typeutil"
	"github.com/leizongmin/jssh/internal/jsexecutor"
	"path/filepath"
)

func JsFnPathJoin(global typeutil.H) jsexecutor.JSFunction {
	return func(ctx *jsexecutor.JSContext, this jsexecutor.JSValue, args []jsexecutor.JSValue) jsexecutor.JSValue {
		if len(args) < 1 {
			return ctx.String("")
		}
		list := make([]string, 0)
		for _, item := range args {
			list = append(list, item.String())
		}

		return ctx.String(filepath.Join(list...))
	}
}

func JsFnPathAbs(global typeutil.H) jsexecutor.JSFunction {
	return func(ctx *jsexecutor.JSContext, this jsexecutor.JSValue, args []jsexecutor.JSValue) jsexecutor.JSValue {
		if len(args) < 1 {
			return ctx.ThrowSyntaxError("path.abs: missing path name")
		}
		if !args[0].IsString() {
			return ctx.ThrowTypeError("path.abs: first argument expected string type")
		}
		file := args[0].String()

		ret, err := filepath.Abs(file)
		if err != nil {
			return ctx.ThrowError(err)
		}
		return ctx.String(ret)
	}
}

func JsFnPathBase(global typeutil.H) jsexecutor.JSFunction {
	return func(ctx *jsexecutor.JSContext, this jsexecutor.JSValue, args []jsexecutor.JSValue) jsexecutor.JSValue {
		if len(args) < 1 {
			return ctx.ThrowSyntaxError("path.base: missing path name")
		}
		if !args[0].IsString() {
			return ctx.ThrowTypeError("path.base: first argument expected string type")
		}
		file := args[0].String()

		return ctx.String(filepath.Base(file))
	}
}

func JsFnPathExt(global typeutil.H) jsexecutor.JSFunction {
	return func(ctx *jsexecutor.JSContext, this jsexecutor.JSValue, args []jsexecutor.JSValue) jsexecutor.JSValue {
		if len(args) < 1 {
			return ctx.ThrowSyntaxError("path.ext: missing path name")
		}
		if !args[0].IsString() {
			return ctx.ThrowTypeError("path.ext: first argument expected string type")
		}
		file := args[0].String()

		return ctx.String(filepath.Ext(file))
	}
}

func JsFnPathDir(global typeutil.H) jsexecutor.JSFunction {
	return func(ctx *jsexecutor.JSContext, this jsexecutor.JSValue, args []jsexecutor.JSValue) jsexecutor.JSValue {
		if len(args) < 1 {
			return ctx.ThrowSyntaxError("path.dir: missing path name")
		}
		if !args[0].IsString() {
			return ctx.ThrowTypeError("path.dir: first argument expected string type")
		}
		file := args[0].String()

		return ctx.String(filepath.Dir(file))
	}
}
