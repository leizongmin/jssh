package jsshcmd

import (
	"github.com/leizongmin/go/typeutil"
	"net/http"
	"os"
	"os/user"
	"strings"
)

func mustGetHomeDir() string {
	dir, err := os.UserHomeDir()
	if err != nil {
		errLog.Fatalln(err)
	}
	return dir
}

func mustGetCurrentUsername() string {
	u, err := user.Current()
	if err != nil {
		errLog.Fatalln(err)
	}
	return u.Username
}

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

func getHeaderMap(header http.Header) typeutil.H {
	ret := make(typeutil.H)
	for name, values := range header {
		name = strings.ToLower(name)
		if len(values) > 1 {
			ret[name] = values
		} else {
			ret[name] = values[0]
		}
	}
	return ret
}
