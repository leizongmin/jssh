package jsshcmd

import (
	"bytes"
	"fmt"
	"github.com/leizongmin/go/typeutil"
	"github.com/leizongmin/jssh/internal/jsexecutor"
	"os"
	"os/exec"
	"sync"
)

func jsFnShSetenv(global typeutil.H) jsexecutor.JSFunction {
	return func(ctx *jsexecutor.JSContext, this jsexecutor.JSValue, args []jsexecutor.JSValue) jsexecutor.JSValue {
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
		jsexecutor.MergeMapToJSObject(ctx, ctx.Globals(), global)

		return ctx.Bool(true)
	}
}

func jsFnShChdir(global typeutil.H) jsexecutor.JSFunction {
	return func(ctx *jsexecutor.JSContext, this jsexecutor.JSValue, args []jsexecutor.JSValue) jsexecutor.JSValue {
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

func jsFnShCwd(global typeutil.H) jsexecutor.JSFunction {
	return func(ctx *jsexecutor.JSContext, this jsexecutor.JSValue, args []jsexecutor.JSValue) jsexecutor.JSValue {
		dir, err := os.Getwd()
		if err != nil {
			return ctx.ThrowError(err)
		}
		return ctx.String(dir)
	}
}

func jsFnShExec(global typeutil.H) jsexecutor.JSFunction {
	return func(ctx *jsexecutor.JSContext, this jsexecutor.JSValue, args []jsexecutor.JSValue) jsexecutor.JSValue {
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
				second, err := jsexecutor.JSValueToAny(args[1])
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

		saveOutput := false
		pipeOutput := true
		if len(args) >= 3 {
			if !args[2].IsNumber() {
				return ctx.ThrowTypeError("exec: third argument expected number type")
			}
			mode := args[2].Int32()
			switch mode {
			case 0:
				saveOutput = false
				pipeOutput = true
			case 1:
				saveOutput = true
				pipeOutput = false
			case 2:
				saveOutput = true
				pipeOutput = true
			default:
				return ctx.ThrowTypeError("exec: mode expected one of 0,1,2")
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

			var saveBuffer *bytes.Buffer
			if saveOutput {
				saveBuffer = bytes.NewBuffer(nil)
			}

			go func() {
				if _, err := pipeBufferAndSave(os.Stdout, stdout, saveBuffer); err != nil {
					if err != os.ErrClosed {
						stdLog.Printf("exec: [stdout] %s", err)
					}
				}
				wg.Done()
			}()
			go func() {
				if _, err := pipeBufferAndSave(os.Stderr, stderr, saveBuffer); err != nil {
					if err != os.ErrClosed {
						stdLog.Printf("exec: [stderr] %s", err)
					}
				}
				wg.Done()
			}()

			if err := sh.Start(); err != nil {
				return ctx.ThrowError(err)
			}
			wg.Wait()

			if err := stdout.Close(); err != nil {
				stdLog.Printf("exec: [stdout] %s", err)
			}
			if err := stderr.Close(); err != nil {
				stdLog.Printf("exec: [stderr] %s", err)
			}
			if err := sh.Wait(); err != nil {
				stdLog.Printf("exec: %s", err)
			}

			var output []byte
			if saveBuffer != nil {
				output = saveBuffer.Bytes()
			}
			global["__output"] = string(output)
			global["__outputbytes"] = len(output)
		} else {

			out, err := sh.CombinedOutput()
			if err != nil {
				stdLog.Printf("exec: %s", err)
			}
			global["__output"] = string(out)
			global["__outputbytes"] = len(out)
		}

		pid := 0
		if sh.Process != nil {
			pid = sh.Process.Pid
		}

		code := sh.ProcessState.ExitCode()
		global["__code"] = code
		jsexecutor.MergeMapToJSObject(ctx, ctx.Globals(), global)

		return jsexecutor.AnyToJSValue(ctx, typeutil.H{
			"pid":         pid,
			"code":        code,
			"output":      global["__output"],
			"outputbytes": global["__outputbytes"],
		})
	}
}

func jsFnShBgexec(global typeutil.H) jsexecutor.JSFunction {
	return func(ctx *jsexecutor.JSContext, this jsexecutor.JSValue, args []jsexecutor.JSValue) jsexecutor.JSValue {
		if len(args) < 1 {
			return ctx.ThrowSyntaxError("bgexec: missing exec command")
		}
		if !args[0].IsString() {
			return ctx.ThrowTypeError("bgexec: first argument expected string type")
		}
		cmd := args[0].String()

		env := cloneMap(global["__env"].(typeutil.H))
		if len(args) >= 2 {
			if args[1].IsNull() || args[1].IsUndefined() {
			} else {
				if !args[1].IsObject() {
					return ctx.ThrowTypeError("bgexec: second argument expected an object")
				}
				second, err := jsexecutor.JSValueToAny(args[1])
				if err != nil {
					return ctx.ThrowError(err)
				}
				env2, ok := second.(typeutil.H)
				if !ok {
					return ctx.ThrowTypeError("bgexec: second argument expected an object")
				}
				for n, v := range env2 {
					env[n] = v
				}
			}
		}

		saveOutput := false
		pipeOutput := true
		if len(args) >= 3 {
			if !args[2].IsNumber() {
				return ctx.ThrowTypeError("bgexec: third argument expected number type")
			}
			mode := args[2].Int32()
			switch mode {
			case 0:
				saveOutput = false
				pipeOutput = true
			case 1:
				saveOutput = true
				pipeOutput = false
			case 2:
				saveOutput = true
				pipeOutput = true
			default:
				return ctx.ThrowTypeError("bgexec: mode expected one of 0,1,2")
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

			var saveBuffer *bytes.Buffer
			if saveOutput {
				saveBuffer = bytes.NewBuffer(nil)
			}

			go func() {
				if _, err := pipeBufferAndSave(os.Stdout, stdout, saveBuffer); err != nil {
					if err != os.ErrClosed {
						stdLog.Printf("bgexec: [stdout] %s", err)
					}
				}
				wg.Done()
			}()
			go func() {
				if _, err := pipeBufferAndSave(os.Stderr, stderr, saveBuffer); err != nil {
					if err != os.ErrClosed {
						stdLog.Printf("bgexec: [stderr] %s", err)
					}
				}
				wg.Done()
			}()

			if err := sh.Start(); err != nil {
				return ctx.ThrowError(err)
			}

			go func() {
				wg.Wait()

				if err := stdout.Close(); err != nil {
					stdLog.Printf("bgexec: [stdout] %s", err)
				}
				if err := stderr.Close(); err != nil {
					stdLog.Printf("bgexec: [stderr] %s", err)
				}
				if err := sh.Wait(); err != nil {
					stdLog.Printf("bgexec: %s", err)
				}
			}()
		} else {
			go func() {
				_, err := sh.CombinedOutput()
				if err != nil {
					stdLog.Printf("bgexec: %s", err)
				}
			}()
		}

		pid := 0
		if sh.Process != nil {
			pid = sh.Process.Pid
		}
		return jsexecutor.AnyToJSValue(ctx, typeutil.H{"pid": pid})
	}
}
