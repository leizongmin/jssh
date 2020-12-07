package jsbuiltin

import (
	jsoniter "github.com/json-iterator/go"
	"log"
	"strings"
)

// JS模块
type JsModule struct {
	File string
	Code string
}

// 获得JS模块列表
func GetJs() []JsModule {
	retModules := make([]JsModule, 0)
	for _, m := range modules {
		// 仅返回内置的模块
		if strings.HasPrefix(m.File, "builtin_") {
			var code string
			if err := jsoniter.UnmarshalFromString(m.Code, &code); err != nil {
				log.Fatalf("jsbuiltin.GetJs: %s", err)
			}
			retModules = append(retModules, JsModule{File: m.File, Code: code})
		}
	}
	return retModules
}
