package main

import (
	"fmt"
	"github.com/leizongmin/go/typeutil"
	"github.com/leizongmin/jssh/internal/scriptx"
	"io"
	"log"
	"os"
	"os/exec"
	"sync"
	"time"
)

func jsFunctionExit(global typeutil.H) scriptx.JSFunction {
	return func(ctx *scriptx.JSContext, this scriptx.JSValue, args []scriptx.JSValue) scriptx.JSValue {
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

func jsFunctionSet(global typeutil.H) scriptx.JSFunction {
	return func(ctx *scriptx.JSContext, this scriptx.JSValue, args []scriptx.JSValue) scriptx.JSValue {
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

func jsFunctionLog(global typeutil.H) scriptx.JSFunction {
	return func(ctx *scriptx.JSContext, this scriptx.JSValue, args []scriptx.JSValue) scriptx.JSValue {
		if len(args) > 0 {
			s, err := scriptx.JSValueToAny(args[0])
			if err != nil {
				return ctx.ThrowError(err)
			}
			format, ok := s.(string)
			a := make([]interface{}, 0)
			for _, v := range args[1:] {
				v2, err := scriptx.JSValueToAny(v)
				if err != nil {
					return ctx.ThrowError(err)
				}
				a = append(a, v2)
			}
			if ok {
				log.Printf(format, a...)
			} else {
				a = append([]interface{}{s}, a...)
				log.Println(a)
			}
		}
		return ctx.Bool(true)
	}
}

func jsFunctionPrint(global typeutil.H) scriptx.JSFunction {
	return func(ctx *scriptx.JSContext, this scriptx.JSValue, args []scriptx.JSValue) scriptx.JSValue {
		if len(args) > 0 {
			s, err := scriptx.JSValueToAny(args[0])
			if err != nil {
				return ctx.ThrowError(err)
			}
			format, ok := s.(string)
			a := make([]interface{}, 0)
			for _, v := range args[1:] {
				v2, err := scriptx.JSValueToAny(v)
				if err != nil {
					return ctx.ThrowError(err)
				}
				a = append(a, v2)
			}
			if ok {
				fmt.Printf(format, a...)
			} else {
				a = append([]interface{}{s}, a...)
				fmt.Println(a)
			}
		}
		return ctx.Bool(true)
	}
}

func jsFunctionPrintln(global typeutil.H) scriptx.JSFunction {
	return func(ctx *scriptx.JSContext, this scriptx.JSValue, args []scriptx.JSValue) scriptx.JSValue {
		ret := global["print"].(scriptx.JSFunction)(ctx, this, args)
		fmt.Println()
		return ret
	}
}

func jsFunctionSetenv(global typeutil.H) scriptx.JSFunction {
	return func(ctx *scriptx.JSContext, this scriptx.JSValue, args []scriptx.JSValue) scriptx.JSValue {
		if len(args) < 1 {
			return ctx.ThrowSyntaxError("setenv: missing env name")
		}
		if !args[0].IsString() {
			return ctx.ThrowTypeError("setenv: first argument expected string type")
		}
		name := args[0].String()

		if len(args) < 2 {
			return ctx.ThrowSyntaxError("setenv: missing env value")
		}
		if !args[1].IsString() {
			return ctx.ThrowTypeError("setenv: second argument expected string type")
		}
		value := args[1].String()

		env := global["__env"].(typeutil.H)
		env[name] = value
		global["__env"] = env
		scriptx.MergeMapToJSObject(ctx, ctx.Globals(), global)

		return ctx.Bool(true)
	}
}

func jsFunctionExec(global typeutil.H) scriptx.JSFunction {
	return func(ctx *scriptx.JSContext, this scriptx.JSValue, args []scriptx.JSValue) scriptx.JSValue {
		if len(args) < 1 {
			return ctx.ThrowSyntaxError("exec: missing exec command")
		}
		if !args[0].IsString() {
			return ctx.ThrowTypeError("exec: first argument expected string type")
		}
		cmd := args[0].String()

		env := cloneMap(global["__env"].(typeutil.H))
		if len(args) >= 2 {
			if args[1].IsNull() || args[1].IsUndefined() {
			} else {
				if !args[1].IsObject() {
					return ctx.ThrowTypeError("exec: second argument expected an object")
				}
				second, err := scriptx.JSValueToAny(args[1])
				if err != nil {
					return ctx.ThrowError(err)
				}
				env2, ok := second.(typeutil.H)
				if !ok {
					return ctx.ThrowTypeError("exec: second argument expected an object")
				}
				for n, v := range env2 {
					env[n] = v
				}
			}
		}

		pipeOutput := true
		if len(args) >= 3 {
			if !args[2].IsBool() {
				return ctx.ThrowTypeError("exec: third argument expected boolean type")
			}
			if args[2].Bool() {
				pipeOutput = false
			}
		}

		sh := exec.Command("sh", "-c", cmd)
		for n, v := range env {
			sh.Env = append(sh.Env, fmt.Sprintf("%s=%s", n, v))
		}

		if pipeOutput {
			sh.Stdin = os.Stdin
			stdout, err := sh.StdoutPipe()
			if err != nil {
				return ctx.ThrowError(err)
			}
			stderr, err := sh.StderrPipe()
			if err != nil {
				return ctx.ThrowError(err)
			}

			var wg sync.WaitGroup
			wg.Add(2)

			go func() {
				if _, err := io.Copy(os.Stdout, stdout); err != nil {
					if err != os.ErrClosed {
						log.Printf("exec: [stdout] %s", err)
					}
				}
				wg.Done()
			}()
			go func() {
				if _, err := io.Copy(os.Stderr, stderr); err != nil {
					if err != os.ErrClosed {
						log.Printf("exec: [stderr] %s", err)
					}
				}
				wg.Done()
			}()

			if err := sh.Start(); err != nil {
				return ctx.ThrowError(err)
			}
			wg.Wait()

			if err := stdout.Close(); err != nil {
				log.Printf("exec: [stdout] %s", err)
			}
			if err := stderr.Close(); err != nil {
				log.Printf("exec: [stderr] %s", err)
			}
			if err := sh.Wait(); err != nil {
				log.Printf("exec: %s", err)
			}
			global["__output"] = ""
			global["__outputbytes"] = 0
		} else {

			out, err := sh.CombinedOutput()
			if err != nil {
				log.Printf("exec: %s", err)
			}
			global["__output"] = string(out)
			global["__outputbytes"] = len(out)
		}

		code := sh.ProcessState.ExitCode()
		global["__code"] = code
		scriptx.MergeMapToJSObject(ctx, ctx.Globals(), global)

		return scriptx.AnyToJSValue(ctx, code)
	}
}

func jsFunctionSleep(global typeutil.H) scriptx.JSFunction {
	return func(ctx *scriptx.JSContext, this scriptx.JSValue, args []scriptx.JSValue) scriptx.JSValue {
		if len(args) < 1 {
			return ctx.ThrowSyntaxError("sleep: missing millisecond argument")
		}
		if !args[0].IsNumber() {
			return ctx.ThrowTypeError("sleep: first argument expected number type")
		}
		ret, err := scriptx.JSValueToAny(args[0])
		if err != nil {
			return ctx.ThrowError(err)
		}
		ms, ok := ret.(float64)
		if !ok {
			return ctx.ThrowTypeError("sleep: first argument expected number type")
		}

		time.Sleep(time.Millisecond * time.Duration(ms))
		return scriptx.AnyToJSValue(ctx, ms)
	}
}

func jsFunctionChdir(global typeutil.H) scriptx.JSFunction {
	return func(ctx *scriptx.JSContext, this scriptx.JSValue, args []scriptx.JSValue) scriptx.JSValue {
		if len(args) < 1 {
			return ctx.ThrowSyntaxError("chdir: missing dir name")
		}
		if !args[0].IsString() {
			return ctx.ThrowTypeError("chdir: first argument expected string type")
		}
		dir := args[0].String()

		if err := os.Chdir(dir); err != nil {
			return ctx.ThrowError(err)
		}
		return ctx.Bool(true)
	}
}

func jsFunctionCwd(global typeutil.H) scriptx.JSFunction {
	return func(ctx *scriptx.JSContext, this scriptx.JSValue, args []scriptx.JSValue) scriptx.JSValue {
		dir, err := os.Getwd()
		if err != nil {
			return ctx.ThrowError(err)
		}
		return ctx.String(dir)
	}
}
