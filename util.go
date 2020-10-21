package main

import (
	"github.com/leizongmin/go/typeutil"
	"os"
	"strings"
)

func getEnvMap() typeutil.H {
	env := make(typeutil.H)
	for _, line := range os.Environ() {
		splits := strings.Split(line, "=")
		k := splits[0]
		v := strings.Join(splits[1:], "=")
		env[k] = v
	}
	return env
}

func cloneMap(a typeutil.H) typeutil.H {
	b := make(typeutil.H)
	for n, v := range a {
		b[n] = v
	}
	return b
}
