package jsbuiltin

import (
	"embed"
	"io/fs"
	"sort"
)

var modules []JsModule

//go:embed 00_module.js
//go:embed 01_assert.js
//go:embed 01_cli.js
//go:embed 01_console.js
//go:embed 01_exec.js
//go:embed 01_fs.js
//go:embed 01_global.js
//go:embed 01_http.js
//go:embed 01_log.js
//go:embed 99_bootstrap.js
var jsFs embed.FS

func init() {
	list, err := jsFs.ReadDir(".")
	if err != nil {
		panic(err)
	}
	sort.Sort(byName(list))
	for _, f := range list {
		modules = append(modules, JsModule{
			File: f.Name(),
			Code: string(mustReadFile(jsFs, f.Name())),
		})
	}
}

func mustReadFile(jsFs embed.FS, name string) []byte {
	b, err := fs.ReadFile(jsFs, name)
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
