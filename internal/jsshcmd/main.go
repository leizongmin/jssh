package jsshcmd

import (
	"fmt"
	"github.com/gookit/color"
	jsoniter "github.com/json-iterator/go"
	"github.com/leizongmin/go/typeutil"
	"github.com/leizongmin/jssh/internal/jsbuiltin"
	"github.com/leizongmin/jssh/internal/jsexecutor"
	"github.com/leizongmin/jssh/internal/pkginfo"
	"github.com/lithdew/quickjs"
	"github.com/peterh/liner"
	"io/ioutil"
	"os"
	"os/exec"
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
		run("", os.Args[1], true, nil)
		return
	}

	if first == "-c" || first == "-x" {
		if len(os.Args) < 3 {
			printUsage(codeFileError)
			return
		}
		if first == "-c" {
			run("", os.Args[2], false, nil)
		} else {
			run("", "return "+os.Args[2], false, func(ret jsexecutor.JSValue) {
				if !ret.IsUndefined() && !ret.IsNull() {
					fmt.Println(ret.String())
				}
			})
		}
		return
	}

	var file string
	var content string
	if isUrl(first) {
		s, err := httpGetFileContent(first)
		if err != nil {
			printExitMessage(err.Error(), codeFileError, false)
		}
		file = first
		content = s
	} else {
		f, err := filepath.Abs(first)
		if err != nil {
			printExitMessage(err.Error(), codeFileError, false)
		}
		file = f
		buf, err := ioutil.ReadFile(file)
		if err != nil {
			printExitMessage(err.Error(), codeFileError, false)
		}
		content = string(buf)
	}
	run(file, content, false, nil)
}

func run(file string, content string, interactive bool, onEnd func(ret jsexecutor.JSValue)) {
	global := getJsGlobal(file)
	jsRuntime := jsexecutor.NewJSRuntime()
	defer jsRuntime.Free()

	ctx := jsRuntime.NewContext()
	defer ctx.Free()
	jsexecutor.MergeMapToJSObject(ctx, ctx.Globals(), global)

	builtinModules := jsbuiltin.GetJs()
	for _, m := range builtinModules {
		ret, err := ctx.EvalFile(m.Code, fmt.Sprintf("internal/%s", m.File))
		if err != nil {
			fmt.Println(color.FgRed.Render(fmt.Sprintf("load builtin js modules fail: %s", formatJsError(err))))
			return
		}
		ret.Free()
	}

	globalFiles := getJsshGlobalFilePath()
	for _, f := range globalFiles {
		if b, err := ioutil.ReadFile(f); err != nil {
			fmt.Println(color.FgRed.Render(err))
		} else {
			content := string(b) + "\n;;"
			if ret, err := ctx.EvalFile(content, f); err != nil {
				fmt.Println(color.FgRed.Render(formatJsError(err)))
			} else {
				ret.Free()
			}
		}
	}
	ctx.Globals().Set("__globalfiles", jsexecutor.AnyToJSValue(ctx, globalFiles))

	if interactive {
		printAuthorInfo()
		fmt.Println("Press Ctrl+C to exit the REPL.")
		fmt.Println("Type ;; in the end of the line to eval.")
		fmt.Println()

		jsGlobals := ctx.Globals()
		repl := liner.NewLiner()
		defer repl.Close()
		repl.SetCtrlCAborts(true)
		repl.SetCompleter(func(line string) (c []string) {
			if names, err := jsGlobals.PropertyNames(); err != nil {
				fmt.Println(color.FgRed.Render(err))
			} else {
				line := strings.ToLower(line)
				for _, n := range names {
					a := jsGlobals.GetByAtom(n.Atom)
					s := n.String()
					if a.IsFunction() {
						s += "("
					}
					if strings.HasPrefix(s, line) {
						c = append(c, s)
					}
					a.Free()
				}
			}
			c = append(c, replCompleter(line)...)
			return c
		})
		prompt := fmt.Sprintf("%s> ", pkginfo.Name)

		historyFile, isEnableHistoryFile := getJsshHistoryFilePath()
		if isEnableHistoryFile {
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
		}

		bufLines := make([]string, 0)
		for {
			code, err := repl.Prompt(prompt)
			if err != nil {
				if err == liner.ErrPromptAborted {
					fmt.Println(color.FgRed.Render("Aborted"))
					break
				} else {
					fmt.Println(color.FgRed.Render(fmt.Sprintf("Error reading line: %s", err)))
				}
			}
			bufLines = append(bufLines, code)
			repl.AppendHistory(code)

			if strings.HasSuffix(code, ";;") {
				content := strings.Join(bufLines, "\n")
				bufLines = make([]string, 0)
				if ret, err := ctx.Eval(content); err != nil {
					fmt.Println(color.FgRed.Render(formatJsError(err)))
				} else {
					jsonPrint := false
					if ret.IsArray() {
						jsonPrint = true
					} else if ret.IsObject() {
						if ret.String() == "[object Object]" {
							jsonPrint = true
						}
					}
					if jsonPrint {
						a, err := jsexecutor.JSValueToAny(ret)
						if err != nil {
							fmt.Println(color.FgRed.Render(err))
							continue
						}
						s, err := jsoniter.MarshalToString(a)
						if err != nil {
							fmt.Println(color.FgRed.Render(err))
							continue
						}
						fmt.Println(color.FgLightBlue.Render(s))
					} else {
						fmt.Println(color.FgLightBlue.Render(ret.String()))
					}
					if onEnd != nil {
						onEnd(ret)
					}
					ret.Free()
				}
			}
		}

		if isEnableHistoryFile {
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
		}
	} else {

		if ret, err := ctx.EvalFile(wrapJsFile(content), file); err != nil {
			printExitMessage(formatJsError(err), codeScriptError, false)
		} else {
			if onEnd != nil {
				onEnd(ret)
			}
			ret.Free()
		}
	}
}

func formatJsError(err error) string {
	if err2, ok := err.(*quickjs.Error); ok {
		return fmt.Sprintf("%s\n%s", err2.Cause, err2.Stack)
	}
	return err.Error()
}

func wrapJsFile(code string) string {
	code = removeShebangLine(code)
	return fmt.Sprintf("(function () { %s\n})();", code)
}

func removeShebangLine(code string) string {
	if !strings.HasPrefix(code, "#!") {
		return code
	}
	lines := strings.SplitN(code, "\n", 2)
	if len(lines) < 2 {
		return "\n"
	}
	return "\n" + lines[1]
}

func getJsshHistoryFilePath() (f string, enable bool) {
	e := os.Getenv("JSSH_HISTORY")
	if len(e) < 1 {
		return filepath.Join(mustGetHomeDir(), fmt.Sprintf(".%s_history", pkginfo.Name)), true
	}
	if e == "0" {
		return "", false
	}
	f, err := filepath.Abs(e)
	if err != nil {
		errLog.Println(color.FgRed.Render(fmt.Sprintf("cannot get file path from environment variable [JSSH_HISTORY]: %s", err)))
		return "", false
	}
	return f, true
}

func getJsshGlobalFilePath() (list []string) {
	e := os.Getenv("JSSH_GLOBAL")
	var dir string
	if len(e) < 1 {
		f := filepath.Join(mustGetHomeDir(), fmt.Sprintf(".%s_global.js", pkginfo.Name))
		s := tryGetFileStat(f)
		if s != nil && !s.IsDir() {
			list = append(list, f)
		}
		dir = filepath.Join(mustGetHomeDir(), fmt.Sprintf(".%s_global.d", pkginfo.Name))
	} else if e != "0" {
		f, err := filepath.Abs(e)
		if err != nil {
			errLog.Println(color.FgRed.Render(fmt.Sprintf("cannot get global file from environment variable [JSSH_GLOBAL]: %s", err)))
		} else {
			s := tryGetFileStat(f)
			if s == nil {
				errLog.Println(color.FgRed.Render(fmt.Sprintf("cannot get global file from environment variable [JSSH_GLOBAL]: %s", err)))
			} else {
				if s.IsDir() {
					dir = f
				} else {
					list = append(list, f)
				}
			}
		}
	}
	if len(dir) > 0 {
		files, err := ioutil.ReadDir(dir)
		if err != nil {
			if !strings.HasSuffix(err.Error(), "no such file or directory") {
				errLog.Println(color.FgRed.Render(fmt.Sprintf("cannot get global file from environment variable [JSSH_GLOBAL]: %s", err)))
			}
		} else {
			for _, f := range files {
				if strings.HasSuffix(f.Name(), ".js") && !f.IsDir() {
					list = append(list, filepath.Join(dir, f.Name()))
				}
			}
		}
	}
	return list
}

func tryGetFileStat(name string) os.FileInfo {
	s, err := os.Stat(name)
	if err != nil {
		return nil
	}
	return s
}

func getJsGlobal(file string) typeutil.H {
	dir := filepath.Dir(file)
	global := make(typeutil.H)

	global["__cpucount"] = runtime.NumCPU()
	global["__os"] = runtime.GOOS
	global["__arch"] = runtime.GOARCH
	global["__version"] = pkginfo.LongVersion
	global["__bin"] = getCurrentAbsoluteBinPath()
	global["__pid"] = os.Getpid()
	global["__user"] = mustGetCurrentUsername()
	global["__tmpdir"] = os.TempDir()
	global["__homedir"] = mustGetHomeDir()
	global["__hostname"], _ = os.Hostname()
	global["__dirname"] = dir
	global["__filename"] = file
	global["__args"] = os.Args[:]
	global["__env"] = getEnvMap()
	global["__output"] = ""
	global["__outputbytes"] = 0
	global["__code"] = 0

	global["format"] = jsFnFormat(global)
	global["print"] = jsFnPrint(global)
	global["readline"] = jsFnReadline(global)
	global["stdoutlog"] = jsFnStdoutlog(global)
	global["stderrlog"] = jsFnStderrlog(global)
	global["evalfile"] = jsFnEvalfile(global)
	global["bytesize"] = jsFnBytesize(global)

	global["sleep"] = jsFnSleep(global)
	global["exit"] = jsFnExit(global)
	global["loadconfig"] = jsFnLoadconfig(global)

	global["base64encode"] = jsFnBase64encode(global)
	global["base64decode"] = jsFnBase64decode(global)
	global["md5"] = jsFnMd5(global)
	global["sha1"] = jsFnSha1(global)
	global["sha256"] = jsFnSha256(global)

	global["networkinterfaces"] = jsFnNetworkinterfaces(global)

	global["setenv"] = jsFnShSetenv(global)
	global["chdir"] = jsFnShChdir(global)
	global["cd"] = jsFnShChdir(global)
	global["cwd"] = jsFnShCwd(global)
	global["pwd"] = jsFnShCwd(global)
	global["exec"] = jsFnShExec(global)
	global["bgexec"] = jsFnShBgexec(global)
	global["pty"] = jsFnShPty(global)

	sshModule := make(typeutil.H)
	sshModule["set"] = jsFnSshSet(global)
	sshModule["open"] = jsFnSshOpen(global)
	sshModule["close"] = jsFnSshClose(global)
	sshModule["setenv"] = jsFnSshSetenv(global)
	sshModule["exec"] = jsFnSshExec(global)
	global["ssh"] = sshModule

	fsModule := make(typeutil.H)
	fsModule["readdir"] = jsFnFsReaddir(global)
	fsModule["readfile"] = jsFnFsReadfile(global)
	fsModule["stat"] = jsFnFsStat(global)
	fsModule["exist"] = jsFnFsExist(global)
	fsModule["writefile"] = jsFnFsWritefile(global)
	fsModule["appendfile"] = jsFnFsAppendfile(global)
	global["fs"] = fsModule

	pathModule := make(typeutil.H)
	pathModule["join"] = jsFnPathJoin(global)
	pathModule["abs"] = jsFnPathAbs(global)
	pathModule["base"] = jsFnPathBase(global)
	pathModule["ext"] = jsFnPathExt(global)
	pathModule["dir"] = jsFnPathDir(global)
	global["path"] = pathModule

	httpModule := make(typeutil.H)
	httpModule["timeout"] = jsFnHttpTimeout(global)
	httpModule["request"] = jsFnHttpRequest(global)
	httpModule["download"] = jsFnHttpDownload(global)
	global["http"] = httpModule

	socketModule := make(typeutil.H)
	socketModule["timeout"] = jsFnSocketTimeout(global)
	socketModule["tcpsend"] = jsFnSocketTcpsend(global)
	socketModule["tcptest"] = jsFnSocketTcptest(global)
	global["socket"] = socketModule

	for n, v := range registeredGlobal {
		global[n] = v
	}

	return global
}

func getCurrentAbsoluteBinPath() string {
	bin := os.Args[0]
	if filepath.IsAbs(bin) {
		return bin
	}
	ret, err := exec.LookPath(bin)
	if err != nil {
		errLog.Println(color.FgRed.Render(err.Error()))
		return bin
	}
	ret2, err := filepath.Abs(ret)
	if err != nil {
		errLog.Println(color.FgRed.Render(err.Error()))
		return bin
	}
	return ret2
}
