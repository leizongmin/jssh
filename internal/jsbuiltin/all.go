package jsbuiltin

import (
	"embed"
	_ "embed"
	"io/fs"
	"sort"
)

var modules []JsModule

//go:embed builtin_0.js
//go:embed builtin_assert.js
//go:embed builtin_cli.js
//go:embed builtin_console.js
//go:embed builtin_exec.js
//go:embed builtin_fs.js
//go:embed builtin_global.js
//go:embed builtin_log.js
//go:embed builtin_http.js
var files embed.FS

func init() {
	list, err := files.ReadDir(".")
	if err != nil {
		panic(err)
	}
	sort.Sort(byName(list))
	for _, f := range list {
		modules = append(modules, JsModule{
			File: f.Name(),
			Code: string(mustReadFile(files, f.Name())),
		})
	}
}

func mustReadFile(f embed.FS, name string) []byte {
	b, err := fs.ReadFile(f, name)
	if err != nil {
		panic(err)
	}
	return b
}

type byName []fs.DirEntry

func (b byName) Len() int {
	return len(b)
}

func (b byName) Less(i, j int) bool {
	return b[i].Name() < b[j].Name()
}

func (b byName) Swap(i, j int) {
	b[i], b[j] = b[j], b[i]
}
