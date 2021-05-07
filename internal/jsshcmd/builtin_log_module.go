package jsshcmd

import (
	"fmt"
	"log"
	"os"

	"github.com/leizongmin/go/typeutil"

	"github.com/leizongmin/jssh/internal/jsexecutor"
	"github.com/leizongmin/jssh/internal/pkginfo"
)

var logPrefix = fmt.Sprintf("[%s] ", pkginfo.Name)
var stdLog = log.New(os.Stdout, logPrefix, log.LstdFlags)
var errLog = log.New(os.Stderr, logPrefix, log.LstdFlags)

func jsFnFormat(global typeutil.H) jsexecutor.JSFunction {
	return func(ctx *jsexecutor.JSContext, this jsexecutor.JSValue, args []jsexecutor.JSValue) jsexecutor.JSValue {
		if len(args) > 0 {
			s, err := jsexecutor.JSValueToAny(args[0])
			if err != nil {
				return ctx.ThrowError(err)
			}
			format, firstIsString := s.(string)
			a := make([]interface{}, 0)
			for _, v := range args[1:] {
				v2, err := jsexecutor.JSValueToAny(v)
				if err != nil {
					return ctx.ThrowError(err)
				}
				a = append(a, v2)
			}
			if firstIsString {
				if len(a) > 0 {
					return ctx.String(fmt.Sprintf(format, a...))
				}
				return ctx.String(format)
			}
			a = append([]interface{}{s}, a...)
			return ctx.String(fmt.Sprint(a...))
		}
		return ctx.String("")
	}
}

func jsFnPrint(global typeutil.H) jsexecutor.JSFunction {
	return func(ctx *jsexecutor.JSContext, this jsexecutor.JSValue, args []jsexecutor.JSValue) jsexecutor.JSValue {
		if len(args) > 0 {
			s, err := jsexecutor.JSValueToAny(args[0])
			if err != nil {
				return ctx.ThrowError(err)
			}
			format, ok := s.(string)
			a := make([]interface{}, 0)
			for _, v := range args[1:] {
				v2, err := jsexecutor.JSValueToAny(v)
				if err != nil {
					return ctx.ThrowError(err)
				}
				a = append(a, v2)
			}
			if ok {
				fmt.Printf(format, a...)
			} else {
				a = append([]interface{}{s}, a...)
				fmt.Print(a...)
			}
		}
		return ctx.Bool(true)
	}
}

func jsFnStdoutlog(global typeutil.H) jsexecutor.JSFunction {
	return func(ctx *jsexecutor.JSContext, this jsexecutor.JSValue, args []jsexecutor.JSValue) jsexecutor.JSValue {
		if len(args) < 1 {
			return ctx.ThrowTypeError("stdoutlog: missing log line argument")
		}
		if !args[0].IsString() {
			return ctx.ThrowTypeError("stdoutlog: first argument expected string type")
		}
		line := args[0].String()
		stdLog.Println(line)
		return ctx.Bool(true)
	}
}

func jsFnStderrlog(global typeutil.H) jsexecutor.JSFunction {
	return func(ctx *jsexecutor.JSContext, this jsexecutor.JSValue, args []jsexecutor.JSValue) jsexecutor.JSValue {
		if len(args) < 1 {
			return ctx.ThrowTypeError("stderrlog: missing log line argument")
		}
		if !args[0].IsString() {
			return ctx.ThrowTypeError("stderrlog: first argument expected string type")
		}
		line := args[0].String()
		errLog.Println(line)
		return ctx.Bool(true)
	}
}
