package jsshcmd

import (
	"bytes"
	"fmt"
	"io"
	"math"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/leizongmin/jssh/internal/jsexecutor"
	"github.com/leizongmin/jssh/internal/pkginfo"
	"github.com/leizongmin/jssh/internal/utils"
	"github.com/leizongmin/jssh/internal/utils/httputil"
)

var httpGlobalHeaders map[string]string
var httpGlobalTimeout int64 = 60_000

func init() {
	httpGlobalHeaders = make(map[string]string)
	httpGlobalHeaders["user-agent"] = fmt.Sprintf("%s/%s", pkginfo.Name, pkginfo.LongVersion)
}

func jsFnHttpTimeout(global utils.H) jsexecutor.JSFunction {
	return func(ctx *jsexecutor.JSContext, this jsexecutor.JSValue, args []jsexecutor.JSValue) jsexecutor.JSValue {
		if len(args) < 1 {
			return ctx.ThrowSyntaxError("http.timeout: missing timeout millisecond")
		}
		if !args[0].IsNumber() {
			return ctx.ThrowTypeError("http.timeout: first argument expected number type")
		}

		httpGlobalTimeout = args[0].Int64()
		httputil.GetClient().Timeout = time.Millisecond * time.Duration(httpGlobalTimeout)

		return ctx.Int64(httpGlobalTimeout)
	}
}

func jsFnHttpRequest(global utils.H) jsexecutor.JSFunction {
	return func(ctx *jsexecutor.JSContext, this jsexecutor.JSValue, args []jsexecutor.JSValue) jsexecutor.JSValue {
		if len(args) < 1 {
			return ctx.ThrowSyntaxError("http.request: missing request method")
		}
		if !args[0].IsString() {
			return ctx.ThrowTypeError("http.request: first argument expected string type")
		}
		method := args[0].String()

		if len(args) < 2 {
			return ctx.ThrowSyntaxError("http.request: missing request url")
		}
		if !args[1].IsString() {
			return ctx.ThrowTypeError("http.request: second argument expected string type")
		}
		url := args[1].String()

		req := httputil.Request()
		req.Method = strings.ToUpper(method)
		req.URL = url
		for n, v := range httpGlobalHeaders {
			req.SetHeader(n, v)
		}

		if len(args) >= 3 {
			if !args[2].IsObject() {
				return ctx.ThrowTypeError("http.request: third argument expected object type")
			}
			third, err := jsexecutor.JSValueToAny(args[2])
			if err != nil {
				return ctx.ThrowError(err)
			}
			headers, ok := third.(utils.H)
			if !ok {
				return ctx.ThrowTypeError("http.request: third argument expected an object")
			}
			for n, v := range headers {
				if s, ok := v.(string); ok {
					req.SetHeader(n, s)
				}
			}
		}

		if len(args) >= 4 {
			if !args[3].IsString() {
				return ctx.ThrowTypeError("http.request: fourth argument expected string type")
			}
			body := args[3].String()
			req.WithBody(bytes.NewReader([]byte(body)))
		}

		res, err := req.Send()
		if err != nil {
			return ctx.ThrowError(err)
		}
		defer func() {
			if err := res.Close(); err != nil {
				stdLog.Printf("http.request: %s", err)
			}
		}()

		ret := make(utils.H)
		ret["status"] = res.Status()
		ret["url"] = res.URL()
		ret["headers"] = getHeaderMap(res.Header())
		b, err := res.Body()
		if err != nil {
			if err != io.EOF {
				return ctx.ThrowError(err)
			}
		}
		ret["body"] = string(b)
		return jsexecutor.AnyToJSValue(ctx, ret)
	}
}

func jsFnHttpDownload(global utils.H) jsexecutor.JSFunction {
	return func(ctx *jsexecutor.JSContext, this jsexecutor.JSValue, args []jsexecutor.JSValue) jsexecutor.JSValue {
		if len(args) < 1 {
			return ctx.ThrowSyntaxError("http.download: missing request url")
		}
		if !args[0].IsString() {
			return ctx.ThrowTypeError("http.download: first argument expected string type")
		}
		url := args[0].String()

		var filename string
		if len(args) >= 2 {
			if !args[1].IsString() {
				return ctx.ThrowTypeError("http.download: second argument expected string type")
			}
			filename = args[1].String()
		}
		if len(filename) < 1 {
			filename = filepath.Join(os.TempDir(), fmt.Sprintf("jssh-http-download-%d-%d", time.Now().Unix(), utils.Int63n(math.MaxInt64)))
		}

		req := httputil.Request()
		req.Method = "GET"
		req.URL = url
		for n, v := range httpGlobalHeaders {
			req.SetHeader(n, v)
		}

		res, err := req.Send()
		if err != nil {
			return ctx.ThrowError(err)
		}
		defer func() {
			if err := res.Close(); err != nil {
				stdLog.Printf("http.download: %s", err)
			}
		}()

		f, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
		if err != nil {
			return ctx.ThrowError(err)
		}
		if _, err := io.Copy(f, res.Origin().Body); err != nil {
			return ctx.ThrowError(err)
		}
		if err := res.Origin().Body.Close(); err != nil {
			stdLog.Printf("http.download: %s", err)
		}
		if err := f.Close(); err != nil {
			return ctx.ThrowError(err)
		}

		return ctx.String(filename)
	}
}
