// +build !windows

package jsshcmd

import (
	"fmt"
	"github.com/creack/pty"
	"github.com/leizongmin/go/typeutil"
	"github.com/leizongmin/jssh/internal/jsexecutor"
	"golang.org/x/term"
	"io"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
)

func jsFnShPty(global typeutil.H) jsexecutor.JSFunction {
	return func(ctx *jsexecutor.JSContext, this jsexecutor.JSValue, args []jsexecutor.JSValue) jsexecutor.JSValue {
		if len(args) < 1 {
			return ctx.ThrowSyntaxError("pty: missing exec command")
		}
		if !args[0].IsString() {
			return ctx.ThrowTypeError("pty: first argument expected string type")
		}
		cmd := args[0].String()

		env := cloneMap(global["__env"].(typeutil.H))
		if len(args) >= 2 {
			if args[1].IsNull() || args[1].IsUndefined() {
			} else {
				if !args[1].IsObject() {
					return ctx.ThrowTypeError("pty: second argument expected an object")
				}
				second, err := jsexecutor.JSValueToAny(args[1])
				if err != nil {
					return ctx.ThrowError(err)
				}
				env2, ok := second.(typeutil.H)
				if !ok {
					return ctx.ThrowTypeError("pty: second argument expected an object")
				}
				for n, v := range env2 {
					env[n] = v
				}
			}
		}

		sh := exec.Command("sh", "-c", cmd)
		for n, v := range env {
			sh.Env = append(sh.Env, fmt.Sprintf("%s=%s", n, v))
		}

		ptmx, err := pty.Start(sh)
		if err != nil {
			return ctx.ThrowError(err)
		}
		defer func() {
			if err := ptmx.Close(); err != nil {
				errLog.Printf("pty: %s", err)
			}
		}()

		// 处理resize
		ch := make(chan os.Signal, 1)
		signal.Notify(ch, syscall.SIGWINCH)
		go func() {
			for range ch {
				if err := pty.InheritSize(os.Stdin, ptmx); err != nil {
					errLog.Printf("pty: error resizing pty: %s", err)
				}
			}
		}()
		// 初始化到size
		ch <- syscall.SIGWINCH

		oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
		if err != nil {
			panic(err)
		}
		defer func() {
			if err := term.Restore(int(os.Stdin.Fd()), oldState); err != nil {
				errLog.Printf("pty: %s", err)
			}
		}()

		// 复制pty输出到stdout
		go func() {
			if _, err = io.Copy(ptmx, os.Stdin); err != nil {
				errLog.Printf("pty: %s", err)
			}
		}()
		_, _ = io.Copy(os.Stdout, ptmx)

		pid := 0
		if sh.Process != nil {
			pid = sh.Process.Pid
		}
		return jsexecutor.AnyToJSValue(ctx, typeutil.H{"pid": pid})
	}
}
