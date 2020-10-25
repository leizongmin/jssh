package jsbuiltin

import (
	"encoding/base64"
	"log"
)

func GetJS() string {
	b, err := base64.StdEncoding.DecodeString(code)
	if err != nil {
		log.Fatalf("jsbuiltin.GetJS: %s", err)
	}
	return string(b)
}
