package jsbuiltin

import (
	"strings"
)

// JsModule JS模块
type JsModule struct {
	File string
	Code string
}

// GetJs 获得JS模块列表
func GetJs() []JsModule {
	retModules := make([]JsModule, 0)
	for _, m := range modules {
		// 仅返回内置的模块
		if strings.HasPrefix(m.File, "builtin_") {
			retModules = append(retModules, JsModule{File: m.File, Code: m.Code})
		}
	}
	return retModules
}
