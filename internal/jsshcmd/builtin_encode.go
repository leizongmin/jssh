package jsshcmd

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"github.com/leizongmin/go/typeutil"
	"github.com/leizongmin/jssh/internal/jsexecutor"
)

func JsFnBase64encode(global typeutil.H) jsexecutor.JSFunction {
	return func(ctx *jsexecutor.JSContext, this jsexecutor.JSValue, args []jsexecutor.JSValue) jsexecutor.JSValue {
		if len(args) < 1 {
			return ctx.ThrowSyntaxError("base64encode: missing data")
		}
		if !args[0].IsString() {
			return ctx.ThrowTypeError("base64encode: first argument expected string type")
		}
		data := args[0].String()

		ret := base64.StdEncoding.EncodeToString([]byte(data))
		return ctx.String(ret)
	}
}

func JsFnBase64decode(global typeutil.H) jsexecutor.JSFunction {
	return func(ctx *jsexecutor.JSContext, this jsexecutor.JSValue, args []jsexecutor.JSValue) jsexecutor.JSValue {
		if len(args) < 1 {
			return ctx.ThrowSyntaxError("base64decode: missing data")
		}
		if !args[0].IsString() {
			return ctx.ThrowTypeError("base64decode: first argument expected string type")
		}
		data := args[0].String()

		ret, err := base64.StdEncoding.DecodeString(data)
		if err != nil {
			return ctx.ThrowError(err)
		}
		return ctx.String(string(ret))
	}
}

func JsFnMd5(global typeutil.H) jsexecutor.JSFunction {
	return func(ctx *jsexecutor.JSContext, this jsexecutor.JSValue, args []jsexecutor.JSValue) jsexecutor.JSValue {
		if len(args) < 1 {
			return ctx.ThrowSyntaxError("md5: missing data")
		}
		if !args[0].IsString() {
			return ctx.ThrowTypeError("md5: first argument expected string type")
		}
		data := args[0].String()

		b := md5.Sum([]byte(data))
		ret := fmt.Sprintf("%x", b)
		return ctx.String(ret)
	}
}

func JsFnSha1(global typeutil.H) jsexecutor.JSFunction {
	return func(ctx *jsexecutor.JSContext, this jsexecutor.JSValue, args []jsexecutor.JSValue) jsexecutor.JSValue {
		if len(args) < 1 {
			return ctx.ThrowSyntaxError("sha1: missing data")
		}
		if !args[0].IsString() {
			return ctx.ThrowTypeError("sha1: first argument expected string type")
		}
		data := args[0].String()

		b := sha1.Sum([]byte(data))
		ret := fmt.Sprintf("%x", b)
		return ctx.String(ret)
	}
}

func JsFnSha256(global typeutil.H) jsexecutor.JSFunction {
	return func(ctx *jsexecutor.JSContext, this jsexecutor.JSValue, args []jsexecutor.JSValue) jsexecutor.JSValue {
		if len(args) < 1 {
			return ctx.ThrowSyntaxError("sha256: missing data")
		}
		if !args[0].IsString() {
			return ctx.ThrowTypeError("sha256: first argument expected string type")
		}
		data := args[0].String()

		b := sha256.Sum256([]byte(data))
		ret := fmt.Sprintf("%x", b)
		return ctx.String(ret)
	}
}
