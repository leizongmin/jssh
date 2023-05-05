package jsshcmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/leizongmin/jssh/internal/jsexecutor"
	"github.com/leizongmin/jssh/internal/utils"
)

func fileInfoToMap(s os.FileInfo) utils.H {
	return utils.H{
		"name":    s.Name(),
		"isdir":   s.IsDir(),
		"mode":    uint32(s.Mode()),
		"modtime": s.ModTime().Unix(),
		"size":    s.Size(),
	}
}

func jsFnFsReaddir(global utils.H) jsexecutor.JSFunction {
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
		retList := make([]utils.H, 0)
		for _, item := range list {
			retList = append(retList, fileInfoToMap(item))
		}
		return jsexecutor.AnyToJSValue(ctx, retList)
	}
}

func jsFnFsReadfile(global utils.H) jsexecutor.JSFunction {
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

func jsFnFsReadfilebytes(global utils.H) jsexecutor.JSFunction {
	return func(ctx *jsexecutor.JSContext, this jsexecutor.JSValue, args []jsexecutor.JSValue) jsexecutor.JSValue {
		if len(args) < 1 {
			return ctx.ThrowSyntaxError("fs.readfile: missing path name")
		}
		if !args[0].IsString() {
			return ctx.ThrowTypeError("fs.readfilebytes: first argument expected string type")
		}
		file := args[0].String()

		b, err := ioutil.ReadFile(file)
		if err != nil {
			return ctx.ThrowError(err)
		}

		return jsexecutor.AnyToJSValue(ctx, b)
	}
}

func jsFnFsExist(global utils.H) jsexecutor.JSFunction {
	return func(ctx *jsexecutor.JSContext, this jsexecutor.JSValue, args []jsexecutor.JSValue) jsexecutor.JSValue {
		if len(args) < 1 {
			return ctx.ThrowSyntaxError("fs.exist: missing path name")
		}
		if !args[0].IsString() {
			return ctx.ThrowTypeError("fs.exist: first argument expected string type")
		}
		file := args[0].String()

		if _, err := os.Stat(file); err != nil {
			if strings.HasSuffix(err.Error(), "no such file or directory") {
				return ctx.Bool(false)
			}
			return ctx.ThrowError(err)
		}
		return ctx.Bool(true)
	}
}

func jsFnFsStat(global utils.H) jsexecutor.JSFunction {
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

func jsFnFsWritefile(global utils.H) jsexecutor.JSFunction {
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

		var data []byte
		if args[1].IsString() {
			data = []byte(args[1].String())
		} else if jsexecutor.JSValueIsUint8Array(args[1]) {
			if b, err := jsexecutor.JSValueUint8ArrayToByteSlice(args[1]); err != nil {
				return ctx.ThrowTypeError(fmt.Sprintf("fs.writefile: read Uint8Array data failed: %s", err))
			} else {
				data = b
			}
		} else {
			return ctx.ThrowTypeError("fs.writefile: second argument expected string or Uint8Array type")
		}

		perm := int64(0666)
		if len(args) >= 3 {
			if !args[2].IsNumber() {
				return ctx.ThrowTypeError("fs.writefile: third argument expected number type")
			}
			perm = args[2].Int64()
		}

		if err := ioutil.WriteFile(file, data, os.FileMode(perm)); err != nil {
			return ctx.ThrowError(err)
		}

		return ctx.Bool(true)
	}
}

func jsFnFsAppendfile(global utils.H) jsexecutor.JSFunction {
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

		var data []byte
		if args[1].IsString() {
			data = []byte(args[1].String())
		} else if jsexecutor.JSValueIsUint8Array(args[1]) {
			if b, err := jsexecutor.JSValueUint8ArrayToByteSlice(args[1]); err != nil {
				return ctx.ThrowTypeError(fmt.Sprintf("fs.appendfile: read Uint8Array data failed: %s", err))
			} else {
				data = b
			}
		} else {
			return ctx.ThrowTypeError("fs.appendfile: second argument expected string or Uint8Array type")
		}

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
				stdLog.Printf("fs.appendfile: %s", err)
			}
		}()
		if _, err := f.Write(data); err != nil {
			return ctx.ThrowError(err)
		}

		return ctx.Bool(true)
	}
}
