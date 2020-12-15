package jsshcmd

import (
	"bytes"
	"fmt"
	"github.com/leizongmin/go/httputil"
	"github.com/leizongmin/go/typeutil"
	"io"
	"net/http"
	"os"
	"os/user"
	"path/filepath"
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

func pipeBufferAndSave(dst io.Writer, src io.Reader, saveWriter *bytes.Buffer) (written int64, err error) {
	var buf []byte
	{
		size := 32 * 1024
		if l, ok := src.(*io.LimitedReader); ok && int64(size) > l.N {
			if l.N < 1 {
				size = 1
			} else {
				size = int(l.N)
			}
		}
		buf = make([]byte, size)
	}

	for {
		nr, er := src.Read(buf)
		if nr > 0 {
			nw, ew := dst.Write(buf[0:nr])
			if nw > 0 {
				written += int64(nw)
			}
			if ew != nil {
				err = ew
				break
			}
			if nr != nw {
				err = io.ErrShortWrite
				break
			}

			if saveWriter != nil {
				_, ew := saveWriter.Write(buf[0:nr])
				if ew != nil {
					err = ew
					break
				}
			}
		}
		if er != nil {
			if er != io.EOF {
				err = er
			}
			break
		}
	}

	return written, err
}

func isUrl(s string) bool {
	s = strings.ToLower(s)
	return strings.HasPrefix(s, "http://") || strings.HasPrefix(s, "https://")
}

func httpGetFileContent(url string) (string, error) {
	res, err := httputil.Request().GET(url).Send()
	if err != nil {
		return "", err
	}
	b, err := res.Body()
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func copyFile(src string, dst string) (int64, error) {
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return 0, err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return 0, fmt.Errorf("%s is not a regular file", src)
	}

	source, err := os.Open(src)
	if err != nil {
		return 0, err
	}
	defer source.Close()

	destination, err := os.Create(dst)
	if err != nil {
		return 0, err
	}
	defer destination.Close()
	nBytes, err := io.Copy(destination, source)
	return nBytes, err
}

func crossPlatformFilepathAbs(p string) (string, error) {
	f, err := filepath.Abs(p)
	if err != nil {
		// Windows会报错，需要忽略
		if strings.Contains(err.Error(), "The system cannot find the file specified") {
			return p, nil
		}
		return f, err
	}
	return f, nil
}
