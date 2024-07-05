package jsshcmd

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/chzyer/readline"
	"github.com/gookit/color"

	"github.com/leizongmin/jssh/internal/jsbuiltin"
	"github.com/leizongmin/jssh/internal/jsexecutor"
	"github.com/leizongmin/jssh/internal/pkginfo"
	"github.com/leizongmin/jssh/internal/utils"
	"github.com/leizongmin/jssh/quickjs"
)

const (
	codeOK          = 0 // 正常
	codeSystem      = 1 // 系统错误
	codeFileError   = 2 // 文件错误
	codeScriptError = 3 // 脚本错误
)

const fixedSelfContainedBoundary = "Ac7mr7vA4ih4PhELJtGP1bN9McJ849PvBHfXGUws2Mw8x63EMQIdn5GYdShUE2YT"

func getSelfContainedBoundary() string {
	return "\n----" + fixedSelfContainedBoundary + "-" + pkginfo.CommitHash + "-" + pkginfo.CommitDate + "----\n"
}

// Main 主入口函数
func Main() {
	runtime.LockOSThread()

	if tryRunSelfContainedBoundary() {
		return
	}

	if len(os.Args) < 2 {
		printUsage(codeOK)
		return
	}
	first := os.Args[1]

	if first == "help" || first == "--help" {
		printUsage(codeOK)
		return
	}
	if first == "version" || first == "--version" {
		fmt.Printf("%s %s\n", pkginfo.Name, pkginfo.LongVersion)
		return
	}

	if first == "repl" || first == "--repl" {
		run("", os.Args[1], true, nil, nil)
		return
	}

	if first == "build" || first == "--build" {
		if len(os.Args) < 3 {
			printExitMessage("missing script file", codeFileError, true)
			return
		}
		sourceFile := os.Args[2]

		var targetFile string
		if len(os.Args) >= 4 {
			targetFile = os.Args[3]
		} else {
			base := path.Base(sourceFile)
			ext := path.Ext(base)
			if len(ext) > 0 {
				targetFile = base[0 : len(base)-len(ext)]
			} else {
				targetFile = sourceFile + ".bin"
			}
		}
		if s, err := crossPlatformFilepathAbs(targetFile); err != nil {
			printExitMessage(err.Error(), codeFileError, false)
		} else {
			targetFile = s
		}

		var source string
		if isUrl(sourceFile) {
			s, err := httpGetFileContent(sourceFile)
			if err != nil {
				printExitMessage(err.Error(), codeFileError, false)
			}
			source = s
		} else {
			b, err := ioutil.ReadFile(sourceFile)
			if err != nil {
				printExitMessage(err.Error(), codeFileError, false)
			}
			source = string(b)
		}

		createSelfContainedBinary(source, targetFile)
		return
	}

	if first == "exec" || first == "--exec" || first == "eval" || first == "--eval" {
		if len(os.Args) < 3 {
			printUsage(codeFileError)
			return
		}
		if first == "-c" {
			run("", os.Args[2], false, nil, nil)
		} else {
			run("", "return "+os.Args[2], false, nil, func(ret jsexecutor.JSValue) {
				if !ret.IsUndefined() && !ret.IsNull() {
					printJsValue(ret, false)
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
		f, err := crossPlatformFilepathAbs(first)
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
	run(file, content, false, nil, nil)
}

func tryRunSelfContainedBoundary() bool {
	boundary := getSelfContainedBoundary()
	boundaryBytes := []byte(boundary)
	selfFile := getCurrentAbsoluteBinPath()

	f, err := os.OpenFile(selfFile, os.O_RDONLY, os.ModePerm)
	if err != nil {
		printExitMessage(err.Error(), codeFileError, false)
		return true
	}
	info, err := f.Stat()
	if err != nil {
		printExitMessage(err.Error(), codeFileError, false)
		return true
	}

	endPos := info.Size() - int64(len(boundaryBytes)) - 8
	var sourceLen int64
	{
		bytesLen := len(boundaryBytes) + 8
		b := make([]byte, bytesLen)
		if n, err := f.ReadAt(b, endPos); err != nil {
			printExitMessage(err.Error(), codeFileError, false)
			return true
		} else if n != bytesLen {
			printExitMessage(fmt.Sprintf("read boundary errored, expected %d but got %d bytes", bytesLen, n), codeFileError, false)
			return true
		}
		if string(b[0:len(b)-8]) != boundary {
			return false
		}
		sourceLen = int64(binary.BigEndian.Uint64(b[len(b)-8:]))
	}

	startPos := endPos - sourceLen
	content := make([]byte, sourceLen)
	if n, err := f.ReadAt(content, startPos); err != nil {
		printExitMessage(err.Error(), codeFileError, false)
		return true
	} else if int64(n) != sourceLen {
		printExitMessage(fmt.Sprintf("read source content errored, expected %d but got %d bytes", sourceLen, n), codeFileError, false)
		return true
	}

	run(selfFile, string(content), false, utils.H{
		"__selfcontained": true,
		"__args":          append([]string{selfFile}, os.Args...),
	}, nil)
	return true
}

func createSelfContainedBinary(source string, targetFile string) {
	boundary := getSelfContainedBoundary()
	selfFile := getCurrentAbsoluteBinPath()
	if _, err := copyFile(selfFile, targetFile); err != nil {
		printExitMessage(err.Error(), codeFileError, false)
	}

	lenBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(lenBytes, uint64(len([]byte(source))))

	perm := os.FileMode(0755)
	data := []byte(boundary + source + boundary)
	f, err := os.OpenFile(targetFile, os.O_APPEND|os.O_WRONLY, perm)
	if err != nil {
		printExitMessage(err.Error(), codeFileError, false)
		return
	}
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			log.Println(err)
		}
	}(f)
	if _, err := f.Write(data); err != nil {
		printExitMessage(err.Error(), codeFileError, false)
	}
	if _, err := f.Write(lenBytes); err != nil {
		printExitMessage(err.Error(), codeFileError, false)
	}

	if err := os.Chmod(targetFile, perm); err != nil {
		printExitMessage(err.Error(), codeFileError, false)
	}

	fmt.Printf("Created self-contained binary file: %s\n", targetFile)
}

func run(file string, content string, interactive bool, customGlobal utils.H, onEnd func(ret jsexecutor.JSValue)) {
	global := getJsGlobal(file)
	jsRuntime := jsexecutor.NewJSRuntime()
	// FIXME: 目前有 bug 导致在 free 的时候会断言失败，目前仅在进程退出时需要 free，因此可以暂时去掉
	// defer jsRuntime.Free()

	if customGlobal != nil {
		for n, v := range customGlobal {
			global[n] = v
		}
	}

	ctx := jsRuntime.NewContext()
	defer ctx.Free()
	jsexecutor.MergeMapToJSObject(ctx, ctx.Globals(), global)

	builtinModules := jsbuiltin.GetJs()
	for _, m := range builtinModules {
		ret, err := ctx.EvalFile(m.Code, fmt.Sprintf("internal:%s", m.File))
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
		historyFile, _ := getJsshHistoryFilePath()
		repl, err := readline.NewEx(&readline.Config{
			Prompt:          fmt.Sprintf("%s> ", pkginfo.Name),
			HistoryFile:     historyFile,
			AutoComplete:    replCompleter(jsGlobals),
			InterruptPrompt: "^C",
			EOFPrompt:       "exit",
		})
		if err != nil {
			fmt.Println(color.FgRed.Render(fmt.Sprintf("Error initializing REPL: %s", err)))
			return
		}
		defer repl.Close()

		bufLines := make([]string, 0)
		abortedCounter := 0
		for {
			code, err := repl.Readline()
			if err != nil {
				if err == readline.ErrInterrupt {
					abortedCounter++
					if abortedCounter < 2 {
						fmt.Println(color.FgYellow.Render("(To exit, press Ctrl+C again or Ctrl+D)"))
						continue
					} else {
						fmt.Println("Bye")
						break
					}
				} else if err == io.EOF {
					fmt.Println("\nBye")
					break
				} else {
					fmt.Println(color.FgRed.Render(fmt.Sprintf("Error reading line: %s", err)))
					continue
				}
			}
			bufLines = append(bufLines, code)
			repl.SaveHistory(code)

			if isCompleteCode(strings.Join(bufLines, "\n")) {
				content := strings.Join(bufLines, "\n")
				bufLines = make([]string, 0)
				if ret, err := ctx.Eval(content); err != nil {
					fmt.Println(color.FgRed.Render(formatJsError(err)))
				} else {
					printJsValue(ret, true)
					if onEnd != nil {
						onEnd(ret)
					}
					ret.Free()
				}
			}

			abortedCounter = 0
			repl.SetPrompt(fmt.Sprintf("%s> ", pkginfo.Name))
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
	f, err := crossPlatformFilepathAbs(e)
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
		f, err := crossPlatformFilepathAbs(e)
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
			if !strings.Contains(err.Error(), "no such file or directory") && !strings.Contains(err.Error(), "The system cannot find the file specified") {
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

func getJsGlobal(file string) utils.H {
	dir := filepath.Dir(file)
	global := make(utils.H)
	jsshNS := make(utils.H)

	jsshNS["__selfcontained"] = false
	jsshNS["__cpucount"] = runtime.NumCPU()
	jsshNS["__os"] = runtime.GOOS
	jsshNS["__arch"] = runtime.GOARCH
	jsshNS["__version"] = pkginfo.LongVersion
	jsshNS["__bin"] = getCurrentAbsoluteBinPath()
	jsshNS["__pid"] = os.Getpid()
	jsshNS["__user"] = mustGetCurrentUsername()
	jsshNS["__tmpdir"] = os.TempDir()
	jsshNS["__homedir"] = mustGetHomeDir()
	jsshNS["__hostname"], _ = os.Hostname()
	jsshNS["__dirname"] = dir
	jsshNS["__filename"] = file
	jsshNS["__args"] = os.Args[:]
	jsshNS["__env"] = getEnvMap()
	jsshNS["__output"] = ""
	jsshNS["__outputbytes"] = 0
	jsshNS["__code"] = 0

	jsshNS["format"] = jsFnFormat(jsshNS)
	jsshNS["print"] = jsFnPrint(jsshNS)
	jsshNS["eprint"] = jsFnEprint(jsshNS)
	jsshNS["readline"] = jsFnReadline(jsshNS)
	jsshNS["stdoutlog"] = jsFnStdoutlog(jsshNS)
	jsshNS["stderrlog"] = jsFnStderrlog(jsshNS)
	jsshNS["evalfile"] = jsFnEvalfile(jsshNS)
	jsshNS["bytesize"] = jsFnBytesize(jsshNS)
	jsshNS["stdin"] = jsFnStdin(jsshNS)
	jsshNS["stdinbytes"] = jsFnStdinbytes(jsshNS)

	jsshNS["sleep"] = jsFnSleep(jsshNS)
	jsshNS["exit"] = jsFnExit(jsshNS)
	jsshNS["loadconfig"] = jsFnLoadconfig(jsshNS)

	jsshNS["base64encode"] = jsFnBase64encode(jsshNS)
	jsshNS["base64decode"] = jsFnBase64decode(jsshNS)
	jsshNS["md5"] = jsFnMd5(jsshNS)
	jsshNS["sha1"] = jsFnSha1(jsshNS)
	jsshNS["sha256"] = jsFnSha256(jsshNS)

	jsshNS["networkinterfaces"] = jsFnNetworkinterfaces(jsshNS)

	jsshNS["setenv"] = jsFnShSetenv(jsshNS)
	jsshNS["chdir"] = jsFnShChdir(jsshNS)
	jsshNS["cd"] = jsFnShChdir(jsshNS)
	jsshNS["cwd"] = jsFnShCwd(jsshNS)
	jsshNS["pwd"] = jsFnShCwd(jsshNS)
	jsshNS["exec"] = jsFnShExec(jsshNS)
	jsshNS["bgexec"] = jsFnShBgexec(jsshNS)
	jsshNS["pty"] = jsFnShPty(jsshNS)

	sshModule := make(utils.H)
	sshModule["set"] = jsFnSshSet(jsshNS)
	sshModule["open"] = jsFnSshOpen(jsshNS)
	sshModule["close"] = jsFnSshClose(jsshNS)
	sshModule["setenv"] = jsFnSshSetenv(jsshNS)
	sshModule["exec"] = jsFnSshExec(jsshNS)
	jsshNS["ssh"] = sshModule

	fsModule := make(utils.H)
	fsModule["readdir"] = jsFnFsReaddir(jsshNS)
	fsModule["readfile"] = jsFnFsReadfile(jsshNS)
	fsModule["stat"] = jsFnFsStat(jsshNS)
	fsModule["exist"] = jsFnFsExist(jsshNS)
	fsModule["writefile"] = jsFnFsWritefile(jsshNS)
	fsModule["appendfile"] = jsFnFsAppendfile(jsshNS)
	fsModule["readfilebytes"] = jsFnFsReadfilebytes(jsshNS)
	fsModule["mkdir"] = jsFnFsMkdir(jsshNS)
	fsModule["mkdirp"] = jsFnFsMkdirp(jsshNS)
	jsshNS["fs"] = fsModule

	pathModule := make(utils.H)
	pathModule["join"] = jsFnPathJoin(jsshNS)
	pathModule["abs"] = jsFnPathAbs(jsshNS)
	pathModule["base"] = jsFnPathBase(jsshNS)
	pathModule["ext"] = jsFnPathExt(jsshNS)
	pathModule["dir"] = jsFnPathDir(jsshNS)
	jsshNS["path"] = pathModule

	httpModule := make(utils.H)
	httpModule["timeout"] = jsFnHttpTimeout(jsshNS)
	httpModule["request"] = jsFnHttpRequest(jsshNS)
	httpModule["download"] = jsFnHttpDownload(jsshNS)
	jsshNS["http"] = httpModule

	socketModule := make(utils.H)
	socketModule["timeout"] = jsFnSocketTimeout(jsshNS)
	socketModule["tcpsend"] = jsFnSocketTcpsend(jsshNS)
	socketModule["tcptest"] = jsFnSocketTcptest(jsshNS)
	jsshNS["socket"] = socketModule

	global["jssh"] = jsshNS
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
	ret2, err := crossPlatformFilepathAbs(ret)
	if err != nil {
		errLog.Println(color.FgRed.Render(err.Error()))
		return bin
	}
	return ret2
}

func printJsValue(ret quickjs.Value, coloured bool) {
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
			if coloured {
				fmt.Println(color.FgRed.Render(err))
			} else {
				fmt.Println(err)
			}
			return
		}
		b, err := json.Marshal(a)
		if err != nil {
			if coloured {
				fmt.Println(color.FgRed.Render(err))
			} else {
				fmt.Println(err)
			}
			return
		}
		if coloured {
			fmt.Println(color.FgLightBlue.Render(string(b)))
		} else {
			fmt.Println(string(b))
		}
	} else {
		if coloured {
			fmt.Println(color.FgLightBlue.Render(ret.String()))
		} else {
			fmt.Println(ret.String())
		}
	}
}

func isCompleteCode(code string) bool {
	openBrackets := 0
	for _, char := range code {
		switch char {
		case '{', '[', '(':
			openBrackets++
		case '}', ']', ')':
			openBrackets--
		}
	}
	return openBrackets == 0
}
