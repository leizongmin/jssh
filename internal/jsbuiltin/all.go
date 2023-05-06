package jsbuiltin

import (
	"embed"
	"io/fs"
	"path/filepath"
	"sort"
)

var modules []JsModule

//go:embed dist/00_module.js
//go:embed dist/01_assert.js
//go:embed dist/01_cli.js
//go:embed dist/01_console.js
//go:embed dist/01_exec.js
//go:embed dist/01_fs.js
//go:embed dist/01_global.js
//go:embed dist/01_http.js
//go:embed dist/01_log.js
var jsFs embed.FS

func init() {
	dir := "dist"
	list, err := jsFs.ReadDir(dir)
	if err != nil {
		panic(err)
	}
	sort.Sort(byName(list))
	for _, f := range list {
		modules = append(modules, JsModule{
			File: f.Name(),
			Code: string(mustReadFile(jsFs, filepath.Join(dir, f.Name()))),
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
