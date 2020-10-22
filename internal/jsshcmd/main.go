package jsshcmd

import (
	"fmt"
	"github.com/gookit/color"
	"github.com/leizongmin/go/cliargs"
	"github.com/leizongmin/go/typeutil"
	"github.com/leizongmin/jssh/internal/jsexecutor"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
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

func Main() {
	runtime.LockOSThread()
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
	global["__bin"], _ = filepath.Abs(os.Args[0])
	global["__pid"] = os.Getpid()
	global["__tmpdir"] = os.TempDir()
	global["__homedir"], _ = os.UserHomeDir()
	global["__hostname"], _ = os.Hostname()
	global["__dirname"] = dir
	global["__filename"] = file
	global["__args"] = os.Args[:]
	global["__env"] = getEnvMap()
	global["__output"] = ""
	global["__outputbytes"] = 0
	global["__code"] = 0

	global["set"] = JsFnSet(global)
	global["sleep"] = JsFnSleep(global)
	global["exit"] = JsFnExit(global)

	global["format"] = JsFnFormat(global)
	global["print"] = JsFnPrint(global)
	global["println"] = JsFnPrintln(global)

	shModule := make(typeutil.H)
	shModule["setenv"] = JsFnShSetenv(global)
	shModule["chdir"] = JsFnShChdir(global)
	shModule["cd"] = JsFnShChdir(global)
	shModule["cwd"] = JsFnShCwd(global)
	shModule["pwd"] = JsFnShCwd(global)
	shModule["exec"] = JsFnShExec(global)
	shModule["bgexec"] = JsFnShBgexec(global)
	global["sh"] = shModule

	logModule := make(typeutil.H)
	logModule["info"] = JsFnLogInfo(global)
	logModule["error"] = JsFnLogError(global)
	global["log"] = logModule

	fsModule := make(typeutil.H)
	fsModule["readdir"] = JsFnFsReaddir(global)
	fsModule["readfile"] = JsFnFsReadfile(global)
	fsModule["stat"] = JsFnFsStat(global)
	fsModule["writefile"] = JsFnFsWritefile(global)
	fsModule["appendfile"] = JsFnFsAppendfile(global)
	global["fs"] = fsModule

	pathModule := make(typeutil.H)
	pathModule["join"] = JsFnPathJoin(global)
	pathModule["abs"] = JsFnPathAbs(global)
	pathModule["base"] = JsFnPathBase(global)
	pathModule["ext"] = JsFnPathExt(global)
	pathModule["dir"] = JsFnPathDir(global)
	global["path"] = pathModule

	cliModule := make(typeutil.H)
	cliModule["get"] = JsFnCliGet(global)
	cliModule["bool"] = JsFnCliBool(global)
	cliModule["args"] = JsFnCliArgs(global)
	cliModule["opts"] = JsFnCliOpts(global)
	global["cli"] = cliModule

	httpModule := make(typeutil.H)
	httpModule["timeout"] = JsFnHttpTimeout(global)
	httpModule["request"] = JsFnHttpRequest(global)
	global["http"] = httpModule

	jsRuntime := jsexecutor.NewJSRuntime()
	defer jsRuntime.Free()
	ret, err := jsexecutor.EvalJSFile(jsRuntime, content, file, global)
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
	fmt.Println(color.FgRed.Render(message))
	if usage {
		fmt.Println()
		printUsage(code)
	} else {
		os.Exit(code)
	}
}
