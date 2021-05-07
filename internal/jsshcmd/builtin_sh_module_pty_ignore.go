// +build windows

package jsshcmd

import (
	"github.com/leizongmin/go/typeutil"

	"github.com/leizongmin/jssh/internal/jsexecutor"
)

func jsFnShPty(global typeutil.H) jsexecutor.JSFunction {
	return func(ctx *jsexecutor.JSContext, this jsexecutor.JSValue, args []jsexecutor.JSValue) jsexecutor.JSValue {
		return ctx.ThrowInternalError("pty: does not supported current OS")
	}
}
