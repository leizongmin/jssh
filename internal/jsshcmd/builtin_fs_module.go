package jsshcmd

import (
	"github.com/leizongmin/go/typeutil"
	"github.com/leizongmin/jssh/internal/jsexecutor"
	"io/ioutil"
	"log"
	"os"
)

func fileInfoToMap(s os.FileInfo) typeutil.H {
	return typeutil.H{
		"name":    s.Name(),
		"isdir":   s.IsDir(),
		"mode":    uint32(s.Mode()),
		"modtime": s.ModTime().Unix(),
		"size":    s.Size(),
	}
}

func JsFnFsReaddir(global typeutil.H) jsexecutor.JSFunction {
	return func(ctx *jsexecutor.JSContext, this jsexecutor.JSValue, args []jsexecutor.JSValue) jsexecutor.JSValue {
		if len(args) < 1 {
			return ctx.ThrowSyntaxError("fs.readdir: missing dir name")
		}
		if !args[0].IsString() {
			return ctx.ThrowTypeError("fs.readdir: first argument expected string type")
		}
		dir := args[0].String()

		list, err := ioutil.ReadDir(dir)
		if err != nil {
			return ctx.ThrowError(err)
		}
		retList := make([]typeutil.H, 0)
		for _, item := range list {
			retList = append(retList, fileInfoToMap(item))
		}
		return jsexecutor.AnyToJSValue(ctx, retList)
	}
}

func JsFnFsReadfile(global typeutil.H) jsexecutor.JSFunction {
	return func(ctx *jsexecutor.JSContext, this jsexecutor.JSValue, args []jsexecutor.JSValue) jsexecutor.JSValue {
		if len(args) < 1 {
			return ctx.ThrowSyntaxError("fs.readfile: missing path name")
		}
		if !args[0].IsString() {
			return ctx.ThrowTypeError("fs.readfile: first argument expected string type")
		}
		file := args[0].String()

		b, err := ioutil.ReadFile(file)
		if err != nil {
			return ctx.ThrowError(err)
		}

		return ctx.String(string(b))
	}
}

func JsFnFsStat(global typeutil.H) jsexecutor.JSFunction {
	return func(ctx *jsexecutor.JSContext, this jsexecutor.JSValue, args []jsexecutor.JSValue) jsexecutor.JSValue {
		if len(args) < 1 {
			return ctx.ThrowSyntaxError("fs.stat: missing path name")
		}
		if !args[0].IsString() {
			return ctx.ThrowTypeError("fs.stat: first argument expected string type")
		}
		file := args[0].String()

		info, err := os.Stat(file)
		if err != nil {
			return ctx.ThrowError(err)
		}

		return jsexecutor.AnyToJSValue(ctx, fileInfoToMap(info))
	}
}

func JsFnFsWritefile(global typeutil.H) jsexecutor.JSFunction {
	return func(ctx *jsexecutor.JSContext, this jsexecutor.JSValue, args []jsexecutor.JSValue) jsexecutor.JSValue {
		if len(args) < 1 {
			return ctx.ThrowSyntaxError("fs.writefile: missing file name")
		}
		if !args[0].IsString() {
			return ctx.ThrowTypeError("fs.writefile: first argument expected string type")
		}
		file := args[0].String()

		if len(args) < 2 {
			return ctx.ThrowSyntaxError("fs.writefile: missing data")
		}
		if !args[1].IsString() {
			return ctx.ThrowTypeError("fs.writefile: second argument expected string type")
		}
		data := args[1].String()

		perm := int64(0666)
		if len(args) >= 3 {
			if !args[2].IsNumber() {
				return ctx.ThrowTypeError("fs.writefile: third argument expected number type")
			}
			perm = args[2].Int64()
		}

		if err := ioutil.WriteFile(file, []byte(data), os.FileMode(perm)); err != nil {
			return ctx.ThrowError(err)
		}

		return ctx.Bool(true)
	}
}

func JsFnFsAppendfile(global typeutil.H) jsexecutor.JSFunction {
	return func(ctx *jsexecutor.JSContext, this jsexecutor.JSValue, args []jsexecutor.JSValue) jsexecutor.JSValue {
		if len(args) < 1 {
			return ctx.ThrowSyntaxError("fs.appendfile: missing file name")
		}
		if !args[0].IsString() {
			return ctx.ThrowTypeError("fs.appendfile: first argument expected string type")
		}
		file := args[0].String()

		if len(args) < 2 {
			return ctx.ThrowSyntaxError("fs.appendfile: missing data")
		}
		if !args[1].IsString() {
			return ctx.ThrowTypeError("fs.appendfile: second argument expected string type")
		}
		data := args[1].String()

		perm := int64(0644)
		if len(args) >= 3 {
			if !args[2].IsNumber() {
				return ctx.ThrowTypeError("fs.appendfile: third argument expected number type")
			}
			perm = args[2].Int64()
		}

		f, err := os.OpenFile(file, os.O_APPEND|os.O_WRONLY, os.FileMode(perm))
		if err != nil {
			return ctx.ThrowError(err)
		}
		defer func() {
			if err := f.Close(); err != nil {
				log.Printf("fs.appendfile: %s", err)
			}
		}()
		if _, err := f.WriteString(data); err != nil {
			return ctx.ThrowError(err)
		}

		return ctx.Bool(true)
	}
}
