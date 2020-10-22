package jsshcmd

import (
	"fmt"
	"github.com/leizongmin/go/cliargs"
	"github.com/leizongmin/go/typeutil"
	"github.com/leizongmin/jssh/internal/jsexecutor"
	"os"
	"strings"
)

var parsedCliArgs *cliargs.CliArgs

func init() {
	if len(os.Args) > 2 {
		parsedCliArgs = cliargs.Parse(os.Args[2:])
	} else {
		parsedCliArgs = cliargs.Parse([]string{})
	}
}

func JsFnCliGet(global typeutil.H) jsexecutor.JSFunction {
	return func(ctx *jsexecutor.JSContext, this jsexecutor.JSValue, args []jsexecutor.JSValue) jsexecutor.JSValue {
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

func JsFnCliBool(global typeutil.H) jsexecutor.JSFunction {
	return func(ctx *jsexecutor.JSContext, this jsexecutor.JSValue, args []jsexecutor.JSValue) jsexecutor.JSValue {
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

func JsFnCliArgs(global typeutil.H) jsexecutor.JSFunction {
	return func(ctx *jsexecutor.JSContext, this jsexecutor.JSValue, args []jsexecutor.JSValue) jsexecutor.JSValue {
		return jsexecutor.AnyToJSValue(ctx, parsedCliArgs.Args)
	}
}

func JsFnCliOpts(global typeutil.H) jsexecutor.JSFunction {
	return func(ctx *jsexecutor.JSContext, this jsexecutor.JSValue, args []jsexecutor.JSValue) jsexecutor.JSValue {
		opts := make(typeutil.H)
		for n, v := range parsedCliArgs.Options {
			opts[n] = v.Value
		}
		return jsexecutor.AnyToJSValue(ctx, opts)
	}
}

func JsFnCliPrompt(global typeutil.H) jsexecutor.JSFunction {
	return func(ctx *jsexecutor.JSContext, this jsexecutor.JSValue, args []jsexecutor.JSValue) jsexecutor.JSValue {
		if len(args) >= 1 {
			if !args[0].IsString() {
				return ctx.ThrowTypeError("cli.bool: first argument expected string type")
			}
			msg := args[0].String()
			if len(msg) > 0 {
				fmt.Print(msg)
			}
		}

		var line string
		_, err := fmt.Scanln(&line)
		if err != nil {
			return ctx.ThrowError(err)
		}
		return ctx.String(line)
	}
}
