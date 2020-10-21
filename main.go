package main

import (
	"fmt"
	"github.com/leizongmin/go/cliargs"
	"github.com/leizongmin/go/typeutil"
	"github.com/leizongmin/jssh/scriptx"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

const cmdName = "jssh"
const cmdVersion = "v1.0"

const (
	codeOK = 0
	codeSystem
	codeFileError
	codeScriptError
)

func haveCliOption(a *cliargs.CliArgs, names ...string) bool {
	for _, n := range names {
		if a.GetOptionOrDefault(n, "false").Value != "false" {
			return true
		}
	}
	return false
}

func main() {
	a := cliargs.Parse(os.Args[1:])

	if haveCliOption(a, "h", "help") {
		printUsage(codeOK)
		return
	}
	if haveCliOption(a, "v", "version") {
		printExitMessage(fmt.Sprintf("%s %s", cmdName, cmdVersion), codeOK)
		return
	}

	file := a.GetArg(0)
	if len(file) < 1 {
		printExitMessage("Missing input script file!", codeFileError)
	}
	file, err := filepath.Abs(file)
	if err != nil {
		printExitMessage(err.Error(), codeFileError)
	}
	dir := filepath.Dir(file)

	buf, err := ioutil.ReadFile(file)
	if err != nil {
		printExitMessage(err.Error(), codeFileError)
	}
	content := string(buf)

	global := make(typeutil.H)
	global["__version"] = cmdVersion
	global["__bin"] = os.Args[0]
	global["__dirname"] = dir
	global["__filename"] = file
	global["__args"] = os.Args[2:]
	global["__env"] = getEnvMap()
	global["log"] = jsFunctionLog(global)
	global["print"] = jsFunctionPrint(global)
	global["println"] = jsFunctionPrintln(global)
	global["setenv"] = jsFunctionSetenv(global)
	global["exec"] = jsFunctionExec(global)
	global["sleep"] = jsFunctionSleep(global)
	global["chdir"] = jsFunctionChdir(global)
	global["cd"] = jsFunctionChdir(global)
	global["cwd"] = jsFunctionCwd(global)
	global["pwd"] = jsFunctionCwd(global)

	jsRuntime := scriptx.NewJSRuntime()
	defer jsRuntime.Free()
	ret, err := scriptx.EvalJS(jsRuntime, content, global)
	if err != nil {
		printExitMessage(err.Error(), codeScriptError)
	}
	defer ret.Free()
	fmt.Printf("%+v", ret)
}

func printUsage(code int) {
	fmt.Printf("Example usage:\n")
	fmt.Printf("  %s <script.js> [arg1] [arg2] [...]\n", cmdName)
	os.Exit(code)
}

func printExitMessage(message string, code int) {
	fmt.Println(message)
	fmt.Println()
	printUsage(code)
}

func getEnvMap() typeutil.H {
	env := make(typeutil.H)
	for _, line := range os.Environ() {
		splits := strings.Split(line, "=")
		k := splits[0]
		v := strings.Join(splits[1:], "=")
		env[k] = v
	}
	return env
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

		g := ctx.Globals()
		g.Set("__env", scriptx.AnyToJSValue(ctx, env))
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
		if len(args) > 1 {
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

		sh := exec.Command("sh", "-c", cmd)
		for n, v := range env {
			sh.Env = append(sh.Env, fmt.Sprintf("%s=%s", n, v))
		}
		stdin, err := sh.StdinPipe()
		if err != nil {
			return ctx.ThrowError(err)
		}
		stdout, err := sh.StdoutPipe()
		if err != nil {
			return ctx.ThrowError(err)
		}
		stderr, err := sh.StderrPipe()
		if err != nil {
			return ctx.ThrowError(err)
		}

		var wg sync.WaitGroup
		wg.Add(3)

		go func() {
			if _, err := io.Copy(os.Stdout, stdout); err != nil {
				if err != os.ErrClosed {
					log.Printf("exec: %s", err)
				}
			}
			wg.Done()
		}()
		go func() {
			if _, err := io.Copy(os.Stderr, stderr); err != nil {
				if err != os.ErrClosed {
					log.Printf("exec: %s", err)
				}
			}
			wg.Done()
		}()
		go func() {
			if _, err := io.Copy(stdin, os.Stdin); err != nil {
				if err != os.ErrClosed {
					log.Printf("exec: %s", err)
				}
			}
			wg.Done()
		}()

		if err := sh.Start(); err != nil {
			return ctx.ThrowError(err)
		}
		wg.Wait()

		if err := stdin.Close(); err != nil {
			log.Printf("exec: %s", err)
		}
		if err := stdout.Close(); err != nil {
			log.Printf("exec: %s", err)
		}
		if err := stderr.Close(); err != nil {
			log.Printf("exec: %s", err)
		}
		if err := sh.Wait(); err != nil {
			log.Printf("exec: %s", err)
		}

		code := sh.ProcessState.ExitCode()
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

func cloneMap(a typeutil.H) typeutil.H {
	b := make(typeutil.H)
	for n, v := range a {
		b[n] = v
	}
	return b
}
