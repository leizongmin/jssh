package main

import (
	"fmt"
	"github.com/leizongmin/go/cliargs"
	"github.com/leizongmin/go/typeutil"
	"github.com/leizongmin/jssh/internal/scriptx"
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
	global["__pid"] = os.Getpid()
	global["__tmpdir"] = os.TempDir()
	global["__homedir"], _ = os.UserHomeDir()
	global["__hostname"], _ = os.Hostname()
	global["__dirname"] = dir
	global["__filename"] = file
	global["__args"] = os.Args[2:]
	global["__env"] = getEnvMap()
	global["__output"] = ""
	global["__outputBytes"] = 0
	global["__code"] = 0

	global["set"] = jsFunctionSet(global)
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
	global["readdir"] = jsFunctionReaddir(global)
	global["readfile"] = jsFunctionReadfile(global)
	global["readstat"] = jsFunctionReadstat(global)
	global["writefile"] = jsFunctionWritefile(global)
	global["appendfile"] = jsFunctionAppendfile(global)
	global["exit"] = jsFunctionExit(global)

	//filepathModule := make(typeutil.H)
	//filepathModule["join"] = jsFunctionFilepathJoin(global)
	//filepathModule["abs"] = jsFunctionFilepathJoin(global)
	//filepathModule["basename"] = jsFunctionFilepathJoin(global)
	//filepathModule["extname"] = jsFunctionFilepathAbs(global)
	//global["filepath"] = filepathModule

	jsRuntime := scriptx.NewJSRuntime()
	defer jsRuntime.Free()
	ret, err := scriptx.EvalJS(jsRuntime, content, global)
	if err != nil {
		printExitMessage(err.Error(), codeScriptError)
	}
	defer ret.Free()
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

func cloneMap(a typeutil.H) typeutil.H {
	b := make(typeutil.H)
	for n, v := range a {
		b[n] = v
	}
	return b
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
			global["__outputBytes"] = 0
		} else {

			out, err := sh.CombinedOutput()
			if err != nil {
				log.Printf("exec: %s", err)
			}
			global["__output"] = string(out)
			global["__outputBytes"] = len(out)
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

func fileInfoToMap(s os.FileInfo) typeutil.H {
	return typeutil.H{
		"name":    s.Name(),
		"isDir":   s.IsDir(),
		"mode":    uint32(s.Mode()),
		"modTime": s.ModTime().Unix(),
		"size":    s.Size(),
	}
}

func jsFunctionReaddir(global typeutil.H) scriptx.JSFunction {
	return func(ctx *scriptx.JSContext, this scriptx.JSValue, args []scriptx.JSValue) scriptx.JSValue {
		if len(args) < 1 {
			return ctx.ThrowSyntaxError("readdir: missing dir name")
		}
		if !args[0].IsString() {
			return ctx.ThrowTypeError("readdir: first argument expected string type")
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
		return scriptx.AnyToJSValue(ctx, retList)
	}
}

func jsFunctionReadfile(global typeutil.H) scriptx.JSFunction {
	return func(ctx *scriptx.JSContext, this scriptx.JSValue, args []scriptx.JSValue) scriptx.JSValue {
		if len(args) < 1 {
			return ctx.ThrowSyntaxError("readfile: missing path name")
		}
		if !args[0].IsString() {
			return ctx.ThrowTypeError("readfile: first argument expected string type")
		}
		file := args[0].String()

		b, err := ioutil.ReadFile(file)
		if err != nil {
			return ctx.ThrowError(err)
		}

		return ctx.String(string(b))
	}
}

func jsFunctionReadstat(global typeutil.H) scriptx.JSFunction {
	return func(ctx *scriptx.JSContext, this scriptx.JSValue, args []scriptx.JSValue) scriptx.JSValue {
		if len(args) < 1 {
			return ctx.ThrowSyntaxError("readstat: missing path name")
		}
		if !args[0].IsString() {
			return ctx.ThrowTypeError("readstat: first argument expected string type")
		}
		file := args[0].String()

		info, err := os.Stat(file)
		if err != nil {
			return ctx.ThrowError(err)
		}

		return scriptx.AnyToJSValue(ctx, fileInfoToMap(info))
	}
}

func jsFunctionWritefile(global typeutil.H) scriptx.JSFunction {
	return func(ctx *scriptx.JSContext, this scriptx.JSValue, args []scriptx.JSValue) scriptx.JSValue {
		if len(args) < 1 {
			return ctx.ThrowSyntaxError("writefile: missing file name")
		}
		if !args[0].IsString() {
			return ctx.ThrowTypeError("writefile: first argument expected string type")
		}
		file := args[0].String()

		if len(args) < 2 {
			return ctx.ThrowSyntaxError("writefile: missing data")
		}
		if !args[1].IsString() {
			return ctx.ThrowTypeError("writefile: second argument expected string type")
		}
		data := args[1].String()

		perm := int64(0666)
		if len(args) >= 3 {
			if !args[2].IsNumber() {
				return ctx.ThrowTypeError("writefile: third argument expected number type")
			}
			perm = args[2].Int64()
		}

		if err := ioutil.WriteFile(file, []byte(data), os.FileMode(perm)); err != nil {
			return ctx.ThrowError(err)
		}

		return ctx.Bool(true)
	}
}

func jsFunctionAppendfile(global typeutil.H) scriptx.JSFunction {
	return func(ctx *scriptx.JSContext, this scriptx.JSValue, args []scriptx.JSValue) scriptx.JSValue {
		if len(args) < 1 {
			return ctx.ThrowSyntaxError("appendfile: missing file name")
		}
		if !args[0].IsString() {
			return ctx.ThrowTypeError("appendfile: first argument expected string type")
		}
		file := args[0].String()

		if len(args) < 2 {
			return ctx.ThrowSyntaxError("appendfile: missing data")
		}
		if !args[1].IsString() {
			return ctx.ThrowTypeError("appendfile: second argument expected string type")
		}
		data := args[1].String()

		perm := int64(0644)
		if len(args) >= 3 {
			if !args[2].IsNumber() {
				return ctx.ThrowTypeError("appendfile: third argument expected number type")
			}
			perm = args[2].Int64()
		}

		f, err := os.OpenFile(file, os.O_APPEND|os.O_WRONLY, os.FileMode(perm))
		if err != nil {
			return ctx.ThrowError(err)
		}
		defer func() {
			if err := f.Close(); err != nil {
				log.Printf("appendfile: %s", err)
			}
		}()
		if _, err := f.WriteString(data); err != nil {
			return ctx.ThrowError(err)
		}

		return ctx.Bool(true)
	}
}

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
