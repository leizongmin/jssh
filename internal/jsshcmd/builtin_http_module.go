package jsshcmd

import (
	"bytes"
	"fmt"
	"github.com/leizongmin/go/httputil"
	"github.com/leizongmin/go/typeutil"
	"github.com/leizongmin/jssh/internal/jsexecutor"
	"io"
	"strings"
	"time"
)

var httpGlobalHeaders map[string]string
var httpGlobalTimeout int64 = 60_000

func init() {
	httpGlobalHeaders = make(map[string]string)
	httpGlobalHeaders["user-agent"] = fmt.Sprintf("%s/%s", cmdName, cmdVersion)
}

func JsFnHttpTimeout(global typeutil.H) jsexecutor.JSFunction {
	return func(ctx *jsexecutor.JSContext, this jsexecutor.JSValue, args []jsexecutor.JSValue) jsexecutor.JSValue {
		if len(args) < 1 {
			return ctx.ThrowSyntaxError("http.timeout: missing timeout millisecond")
		}
		if !args[0].IsNumber() {
			return ctx.ThrowTypeError("http.timeout: first argument expected number type")
		}

		httpGlobalTimeout = args[0].Int64()
		return ctx.Int64(httpGlobalTimeout)
	}
}

func JsFnHttpRequest(global typeutil.H) jsexecutor.JSFunction {
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
		req.Timeout = time.Millisecond * time.Duration(httpGlobalTimeout)

		if len(args) >= 3 {
			if !args[2].IsObject() {
				return ctx.ThrowTypeError("http.request: third argument expected object type")
			}
			third, err := jsexecutor.JSValueToAny(args[1])
			if err != nil {
				return ctx.ThrowError(err)
			}
			headers, ok := third.(typeutil.H)
			if !ok {
				return ctx.ThrowTypeError("http.request: second argument expected an object")
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

		ret := make(typeutil.H)
		ret["status"] = res.Status()
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
