package jsshcmd

import (
	"fmt"
	"github.com/gookit/color"
	jsoniter "github.com/json-iterator/go"
	"github.com/leizongmin/go/cliargs"
	"github.com/leizongmin/go/typeutil"
	"github.com/leizongmin/jssh/internal/jsexecutor"
	"github.com/leizongmin/jssh/internal/pkginfo"
	"github.com/peterh/liner"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

const (
	codeOK          = 0
	codeSystem      = 1
	codeFileError   = 2
	codeScriptError = 3
)

var parsedCliArgs *cliargs.CliArgs

func Main() {
	runtime.LockOSThread()

	if len(os.Args) < 2 {
		printUsage(codeOK)
		return
	}
	first := os.Args[1]

	if first == "-h" || first == "--help" {
		printUsage(codeOK)
		return
	}
	if first == "-v" || first == "--version" {
		fmt.Printf("%s %s\n", pkginfo.Name, pkginfo.LongVersion)
		return
	}

	if first == "-i" {
		parsedCliArgs = cliargs.Parse(os.Args[2:])
		run("", os.Args[1], true)
		return
	}

	if first == "-c" {
		if len(os.Args) < 3 {
			printUsage(codeFileError)
			return
		}
		parsedCliArgs = cliargs.Parse(os.Args[3:])
		run("", os.Args[2], false)
		return
	}

	file, err := filepath.Abs(first)
	if err != nil {
		printExitMessage(err.Error(), codeFileError, false)
	}

	buf, err := ioutil.ReadFile(file)
	if err != nil {
		printExitMessage(err.Error(), codeFileError, false)
	}
	content := string(buf)
	parsedCliArgs = cliargs.Parse(os.Args[2:])
	run(file, content, false)
}

func run(file string, content string, interactive bool) {
	global := getJsGlobal(file)
	jsRuntime := jsexecutor.NewJSRuntime()
	defer jsRuntime.Free()

	ctx := jsRuntime.NewContext()
	defer ctx.Free()
	jsexecutor.MergeMapToJSObject(ctx, ctx.Globals(), global)

	commonFile := filepath.Join(global["__homedir"].(string), fmt.Sprintf(".%src.js", pkginfo.Name))
	if b, err := ioutil.ReadFile(commonFile); err != nil {
		if !strings.HasSuffix(err.Error(), "no such file or directory") {
			fmt.Println(color.FgRed.Render(err))
		}
	} else {
		if _, err := ctx.EvalFile(string(b), commonFile); err != nil {
			fmt.Println(color.FgRed.Render(err))
		}
	}

	if interactive {
		printAuthorInfo()
		fmt.Println("Press Ctrl+C to exit the REPL.")
		fmt.Println("Type ;; in the end of the line to eval.")
		fmt.Println()

		repl := liner.NewLiner()
		defer repl.Close()
		repl.SetCtrlCAborts(true)
		prompt := fmt.Sprintf("%s> ", pkginfo.Name)

		historyFile := filepath.Join(global["__homedir"].(string), fmt.Sprintf(".%s_history", pkginfo.Name))

		if f, err := os.Open(historyFile); err != nil {
			if !strings.HasSuffix(err.Error(), "no such file or directory") {
				fmt.Println(color.FgRed.Render(err))
			}
		} else {
			if _, err := repl.ReadHistory(f); err != nil {
				fmt.Println(color.FgRed.Render(err))
			}
			if err := f.Close(); err != nil {
				fmt.Println(color.FgRed.Render(err))
			}
		}

		bufLines := make([]string, 0)
		for {
			code, err := repl.Prompt(prompt)
			if err != nil {
				if err == liner.ErrPromptAborted {
					fmt.Println(color.FgRed.Render("Aborted"))
					break
				} else {
					fmt.Println(color.FgRed.Render("Error reading line: %s", err))
				}
			}
			bufLines = append(bufLines, code)
			repl.AppendHistory(code)

			if strings.HasSuffix(code, ";;") {
				content := strings.Join(bufLines, "\n")
				bufLines = make([]string, 0)
				if ret, err := ctx.Eval(content); err != nil {
					fmt.Println(color.FgRed.Render(err))
				} else {
					if ret.IsObject() {
						s := ret.String()
						if s != "[object Object]" {
							fmt.Println(color.FgLightBlue.Render(s))
						} else {
							a, err := jsexecutor.JSValueToAny(ret)
							if err != nil {
								fmt.Println(color.FgRed.Render(err))
								continue
							}
							s, err = jsoniter.MarshalToString(a)
							if err != nil {
								fmt.Println(color.FgRed.Render(err))
								continue
							}
							fmt.Println(color.FgLightBlue.Render(s))
						}
					} else {
						fmt.Println(color.FgLightBlue.Render(ret.String()))
					}
					ret.Free()
				}
			}
		}

		if f, err := os.Create(historyFile); err != nil {
			fmt.Println(color.FgRed.Render(err))
		} else {
			if _, err := repl.WriteHistory(f); err != nil {
				fmt.Println(color.FgRed.Render(err))
			} else {
				if err := f.Close(); err != nil {
					fmt.Println(color.FgRed.Render(err))
				}
			}
		}
	} else {

		if ret, err := ctx.EvalFile(content, file); err != nil {
			printExitMessage(err.Error(), codeScriptError, false)
		} else {
			ret.Free()
		}
	}
}

func printAuthorInfo() {
	fmt.Printf("Welcome to %s %s\n", pkginfo.Name, pkginfo.LongVersion)
	fmt.Println("Author:  leizongmin@gmail.com")
	fmt.Println("Project: https://github.com/leizongmin/jssh")
	fmt.Println()
}

func printUsage(code int) {
	printAuthorInfo()
	fmt.Println("Example usage:")
	fmt.Printf("  %s script_file.js [arg1] [arg2] [...]     Run script file\n", pkginfo.Name)
	fmt.Printf("  %s -c \"script\" [arg1] [arg2] [...]        Run script from argument\n", pkginfo.Name)
	fmt.Printf("  %s -i                                     Start REPL\n", pkginfo.Name)
	fmt.Printf("  %s -h                                     Show usage\n", pkginfo.Name)
	fmt.Printf("  %s -v                                     Show version\n", pkginfo.Name)
	fmt.Println()
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

func getJsGlobal(file string) typeutil.H {
	dir := filepath.Dir(file)

	global := make(typeutil.H)

	global["__version"] = pkginfo.LongVersion
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
	global["get"] = JsFnGet(global)
	global["format"] = JsFnFormat(global)
	global["print"] = JsFnPrint(global)
	global["println"] = JsFnPrintln(global)

	global["sleep"] = JsFnSleep(global)
	global["exit"] = JsFnExit(global)
	global["loadconfig"] = JsFnLoadconfig(global)

	shModule := make(typeutil.H)
	shModule["setenv"] = JsFnShSetenv(global)
	shModule["chdir"] = JsFnShChdir(global)
	shModule["cd"] = JsFnShChdir(global)
	shModule["cwd"] = JsFnShCwd(global)
	shModule["pwd"] = JsFnShCwd(global)
	shModule["exec"] = JsFnShExec(global)
	shModule["bgexec"] = JsFnShBgexec(global)
	global["sh"] = shModule

	sshModule := make(typeutil.H)
	sshModule["set"] = JsFnSshSet(global)
	sshModule["open"] = JsFnSshOpen(global)
	sshModule["close"] = JsFnSshClose(global)
	sshModule["setenv"] = JsFnSshSetenv(global)
	sshModule["exec"] = JsFnSshExec(global)
	global["ssh"] = sshModule
	global["__ssh_config"] = typeutil.H{
		"user":    "root",
		"auth":    "key",
		"key":     filepath.Join(global["__homedir"].(string), ".ssh/id_rsa"),
		"keypass": "",
		"port":    22,
		"timeout": 60_000,
	}

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
	cliModule["prompt"] = JsFnCliPrompt(global)
	global["cli"] = cliModule

	httpModule := make(typeutil.H)
	httpModule["timeout"] = JsFnHttpTimeout(global)
	httpModule["request"] = JsFnHttpRequest(global)
	httpModule["download"] = JsFnHttpDownload(global)
	global["http"] = httpModule

	socketModule := make(typeutil.H)
	socketModule["timeout"] = JsFnSocketTimeout(global)
	socketModule["tcpsend"] = JsFnSocketTcpsend(global)
	socketModule["tcptest"] = JsFnSocketTcptest(global)
	global["socket"] = socketModule

	return global
}
