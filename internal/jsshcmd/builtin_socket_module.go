package jsshcmd

import (
	"fmt"
	"github.com/leizongmin/go/typeutil"
	"github.com/leizongmin/jssh/internal/jsexecutor"
	"io/ioutil"
	"net"
	"time"
)

var socketGlobalTimeout int64 = 60_000

func jsFnSocketTimeout(global typeutil.H) jsexecutor.JSFunction {
	return func(ctx *jsexecutor.JSContext, this jsexecutor.JSValue, args []jsexecutor.JSValue) jsexecutor.JSValue {
		if len(args) < 1 {
			return ctx.ThrowSyntaxError("socket.timeout: missing timeout millisecond")
		}
		if !args[0].IsNumber() {
			return ctx.ThrowTypeError("socket.timeout: first argument expected number type")
		}

		socketGlobalTimeout = args[0].Int64()
		return ctx.Int64(socketGlobalTimeout)
	}
}

func jsFnSocketTcpsend(global typeutil.H) jsexecutor.JSFunction {
	return func(ctx *jsexecutor.JSContext, this jsexecutor.JSValue, args []jsexecutor.JSValue) jsexecutor.JSValue {
		if len(args) < 1 {
			return ctx.ThrowSyntaxError("socket.tcpsend: missing host")
		}
		if !args[0].IsString() {
			return ctx.ThrowTypeError("socket.tcpsend: first argument expected string type")
		}
		host := args[0].String()

		if len(args) < 2 {
			return ctx.ThrowSyntaxError("socket.tcpsend: missing port")
		}
		if !args[1].IsNumber() {
			return ctx.ThrowTypeError("socket.tcpsend: second argument expected number type")
		}
		port := int(args[1].Int32())

		if len(args) < 3 {
			return ctx.ThrowSyntaxError("socket.tcpsend: missing data")
		}
		if !args[2].IsString() {
			return ctx.ThrowTypeError("socket.tcpsend: third argument expected string type")
		}
		data := args[2].String()

		timeout := time.Millisecond * time.Duration(socketGlobalTimeout)
		conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%d", host, port), timeout)
		if err != nil {
			return ctx.ThrowError(err)
		}
		defer func() {
			if err := conn.Close(); err != nil {
				errLog.Printf("socket.tcpsend: %s", err)
			}
		}()
		if err := conn.SetDeadline(time.Now().Add(timeout)); err != nil {
			errLog.Printf("socket.tcpsend: %s", err)
		}

		if _, err := conn.Write([]byte(data)); err != nil {
			return ctx.ThrowError(err)
		}
		output, err := ioutil.ReadAll(conn)
		if err != nil {
			return ctx.ThrowError(err)
		}

		return ctx.String(string(output))
	}
}

func jsFnSocketTcptest(global typeutil.H) jsexecutor.JSFunction {
	return func(ctx *jsexecutor.JSContext, this jsexecutor.JSValue, args []jsexecutor.JSValue) jsexecutor.JSValue {
		if len(args) < 1 {
			return ctx.ThrowSyntaxError("socket.tcpsend: missing host")
		}
		if !args[0].IsString() {
			return ctx.ThrowTypeError("socket.tcpsend: first argument expected string type")
		}
		host := args[0].String()

		if len(args) < 2 {
			return ctx.ThrowSyntaxError("socket.tcpsend: missing port")
		}
		if !args[1].IsNumber() {
			return ctx.ThrowTypeError("socket.tcpsend: second argument expected number type")
		}
		port := int(args[1].Int32())

		timeout := time.Millisecond * time.Duration(socketGlobalTimeout)
		conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%d", host, port), timeout)
		if err != nil {
			return ctx.Bool(false)
		}
		defer func() {
			if err := conn.Close(); err != nil {
				errLog.Printf("socket.tcpsend: %s", err)
			}
		}()

		return ctx.Bool(true)
	}
}
