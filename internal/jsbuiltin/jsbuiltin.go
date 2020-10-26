package jsbuiltin

import (
	"encoding/base64"
	"log"
	"strings"
)

func GetJS() string {
	rawModules := make([]string, 0)
	for _, code := range modules {
		b, err := base64.StdEncoding.DecodeString(code)
		if err != nil {
			log.Fatalf("jsbuiltin.GetJS: %s", err)
		}
		rawModules = append(rawModules, string(b))
	}
	return strings.Join(rawModules, "\n")
}
