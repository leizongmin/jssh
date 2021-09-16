package jsbuiltin

import (
	_ "embed"
)

var modules []JsModule

//go:embed builtin_0.js
var codeBuiltin0 string

//go:embed builtin_assert.js
var codeBuiltinAssert string

//go:embed builtin_cli.js
var codeBuiltinCli string

//go:embed builtin_console.js
var codeBuiltinConsole string

//go:embed builtin_exec.js
var codeBuiltinExec string

//go:embed builtin_fs.js
var codeBuiltinFs string

//go:embed builtin_global.js
var codeBuiltinGlobal string

//go:embed builtin_log.js
var codeBuiltinLog string

//go:embed builtin_http.js
var codeBuiltinHttp string

func init() {
	modules = append(modules, JsModule{File: "builtin_0.js", Code: codeBuiltin0})
	modules = append(modules, JsModule{File: "builtin_assert.js", Code: codeBuiltinAssert})
	modules = append(modules, JsModule{File: "builtin_cli.js", Code: codeBuiltinCli})
	modules = append(modules, JsModule{File: "builtin_console.js", Code: codeBuiltinConsole})
	modules = append(modules, JsModule{File: "builtin_exec.js", Code: codeBuiltinExec})
	modules = append(modules, JsModule{File: "builtin_fs.js", Code: codeBuiltinFs})
	modules = append(modules, JsModule{File: "builtin_global.js", Code: codeBuiltinGlobal})
	modules = append(modules, JsModule{File: "builtin_log.js", Code: codeBuiltinLog})
	modules = append(modules, JsModule{File: "builtin_http.js", Code: codeBuiltinHttp})
}
