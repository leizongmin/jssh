package main

import (
	"github.com/leizongmin/go/cliargs"
	"github.com/leizongmin/go/typeutil"
	"github.com/leizongmin/jssh/internal/scriptx"
	"os"
	"strings"
)

var parsedCliArgs *cliargs.CliArgs

func init() {
	parsedCliArgs = cliargs.Parse(os.Args[2:])
}

func JsFnCliGet(global typeutil.H) scriptx.JSFunction {
	return func(ctx *scriptx.JSContext, this scriptx.JSValue, args []scriptx.JSValue) scriptx.JSValue {
		if len(args) < 1 {
			return ctx.ThrowSyntaxError("cli.get: missing index or flag name")
		}
		if args[0].IsNumber() {
			index := int(args[0].Int32())
			if index >= parsedCliArgs.ArgsCount() {
				return ctx.Undefined()
			}
			return ctx.String(parsedCliArgs.GetArg(index))
		}

		if args[0].IsString() {
			name := args[0].String()
			if !parsedCliArgs.HasOption(name) {
				return ctx.Undefined()
			}
			return ctx.String(parsedCliArgs.GetOption(name).Value)
		}

		return ctx.ThrowTypeError("cli.get: first argument expected number or string type")
	}
}

func JsFnCliBool(global typeutil.H) scriptx.JSFunction {
	return func(ctx *scriptx.JSContext, this scriptx.JSValue, args []scriptx.JSValue) scriptx.JSValue {
		if len(args) < 1 {
			return ctx.ThrowSyntaxError("cli.bool: missing flag name")
		}
		if !args[0].IsString() {
			return ctx.ThrowTypeError("cli.bool: first argument expected string type")
		}
		name := args[0].String()

		if !parsedCliArgs.HasOption(name) {
			return ctx.Bool(false)
		}
		opt := strings.ToLower(parsedCliArgs.GetOption(name).Value)
		if opt == "0" || opt == "f" || opt == "false" {
			return ctx.Bool(false)
		}
		return ctx.Bool(true)
	}
}

func JsFnCliArgs(global typeutil.H) scriptx.JSFunction {
	return func(ctx *scriptx.JSContext, this scriptx.JSValue, args []scriptx.JSValue) scriptx.JSValue {
		return scriptx.AnyToJSValue(ctx, parsedCliArgs.Args)
	}
}

func JsFnCliOpts(global typeutil.H) scriptx.JSFunction {
	return func(ctx *scriptx.JSContext, this scriptx.JSValue, args []scriptx.JSValue) scriptx.JSValue {
		opts := make(typeutil.H)
		for n, v := range parsedCliArgs.Options {
			opts[n] = v.Value
		}
		return scriptx.AnyToJSValue(ctx, opts)
	}
}
