package jsbuiltin

import (
	"encoding/base64"
	"log"
	"strings"
)

type JsModule struct {
	File string
	Code string
}

func GetJs() []JsModule {
	retModules := make([]JsModule, 0)
	for _, m := range modules {
		// 仅返回内置的模块
		if strings.HasPrefix(m.File, "builtin_") {
			b, err := base64.StdEncoding.DecodeString(m.Code)
			if err != nil {
				log.Fatalf("jsbuiltin.GetJs: %s", err)
			}
			retModules = append(retModules, JsModule{File: m.File, Code: string(b)})
		}
	}
	return retModules
}
