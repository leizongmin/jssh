package main

import (
	"fmt"
	"github.com/leizongmin/go/cliargs"
	"github.com/leizongmin/go/typeutil"
	"github.com/leizongmin/jssh/internal/scriptx"
	"io/ioutil"
	"os"
	"path/filepath"
)

const cmdName = "jssh"
const cmdVersion = "v1.0"

const (
	codeOK          = 0
	codeSystem      = 1
	codeFileError   = 2
	codeScriptError = 3
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
		printExitMessage(fmt.Sprintf("%s %s", cmdName, cmdVersion), codeOK, false)
		return
	}

	file := a.GetArg(0)
	if len(file) < 1 {
		printExitMessage("Missing input script file!", codeFileError, true)
	}
	file, err := filepath.Abs(file)
	if err != nil {
		printExitMessage(err.Error(), codeFileError, false)
	}
	dir := filepath.Dir(file)

	buf, err := ioutil.ReadFile(file)
	if err != nil {
		printExitMessage(err.Error(), codeFileError, false)
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
	global["exit"] = jsFunctionExit(global)

	fsModule := make(typeutil.H)
	fsModule["readdir"] = jsFunctionFsReaddir(global)
	fsModule["readfile"] = jsFunctionFsReadfile(global)
	fsModule["readstat"] = jsFunctionFsReadstat(global)
	fsModule["writefile"] = jsFunctionFsWritefile(global)
	fsModule["appendfile"] = jsFunctionFsAppendfile(global)
	global["fs"] = fsModule

	pathModule := make(typeutil.H)
	pathModule["join"] = jsFunctionPathJoin(global)
	pathModule["abs"] = jsFunctionPathAbs(global)
	pathModule["base"] = jsFunctionPathBase(global)
	pathModule["ext"] = jsFunctionPathExt(global)
	pathModule["dir"] = jsFunctionPathDir(global)
	global["path"] = pathModule

	cliModule := make(typeutil.H)
	cliModule["get"] = jsFunctionCliGet(global)
	cliModule["bool"] = jsFunctionCliBool(global)
	cliModule["args"] = jsFunctionCliArgs(global)
	cliModule["opts"] = jsFunctionCliOpts(global)
	global["cli"] = cliModule

	jsRuntime := scriptx.NewJSRuntime()
	defer jsRuntime.Free()
	ret, err := scriptx.EvalJS(jsRuntime, content, global)
	if err != nil {
		printExitMessage(err.Error(), codeScriptError, false)
	}
	defer ret.Free()
}

func printUsage(code int) {
	fmt.Printf("Example usage:\n")
	fmt.Printf("  %s <script.js> [arg1] [arg2] [...]\n", cmdName)
	os.Exit(code)
}

func printExitMessage(message string, code int, usage bool) {
	fmt.Println(message)
	if usage {
		fmt.Println()
		printUsage(code)
	} else {
		os.Exit(code)
	}
}
