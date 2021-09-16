package jsshcmd

import (
	"github.com/leizongmin/go/typeutil"
	"github.com/leizongmin/jssh/internal/jsexecutor"
	"path/filepath"
	"strings"
)

func jsFnPathJoin(global typeutil.H) jsexecutor.JSFunction {
	return func(ctx *jsexecutor.JSContext, this jsexecutor.JSValue, args []jsexecutor.JSValue) jsexecutor.JSValue {
		if len(args) < 1 {
			return ctx.String("")
		}
		list := make([]string, 0)
		for _, item := range args {
			list = append(list, item.String())
		}

		if isUrl(list[0]) {
			ret := strings.Join(list, "/")
			ret = strings.ReplaceAll(ret, "/./", "/")
			return ctx.String(ret)
		}

		return ctx.String(filepath.Join(list...))
	}
}

func jsFnPathAbs(global typeutil.H) jsexecutor.JSFunction {
	return func(ctx *jsexecutor.JSContext, this jsexecutor.JSValue, args []jsexecutor.JSValue) jsexecutor.JSValue {
		if len(args) < 1 {
			return ctx.ThrowSyntaxError("path.abs: missing path name")
		}
		if !args[0].IsString() {
			return ctx.ThrowTypeError("path.abs: first argument expected string type")
		}
		file := args[0].String()

		if isUrl(file) {
			return ctx.String(file)
		}

		ret, err := crossPlatformFilepathAbs(file)
		if err != nil {
			return ctx.ThrowError(err)
		}
		return ctx.String(ret)
	}
}

func jsFnPathBase(global typeutil.H) jsexecutor.JSFunction {
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

func jsFnPathExt(global typeutil.H) jsexecutor.JSFunction {
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

func jsFnPathDir(global typeutil.H) jsexecutor.JSFunction {
	return func(ctx *jsexecutor.JSContext, this jsexecutor.JSValue, args []jsexecutor.JSValue) jsexecutor.JSValue {
		if len(args) < 1 {
			return ctx.ThrowSyntaxError("path.dir: missing path name")
		}
		if !args[0].IsString() {
			return ctx.ThrowTypeError("path.dir: first argument expected string type")
		}
		file := args[0].String()

		if isUrl(file) {
			i := strings.LastIndex(file, "/")
			if i == -1 {
				return ctx.String(file)
			}
			return ctx.String(file[:i])
		}

		return ctx.String(filepath.Dir(file))
	}
}
