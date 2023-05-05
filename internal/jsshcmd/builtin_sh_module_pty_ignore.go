//go:build windows
// +build windows

package jsshcmd

import (
	"github.com/leizongmin/jssh/internal/jsexecutor"
	"github.com/leizongmin/jssh/internal/utils"
)

func jsFnShPty(global utils.H) jsexecutor.JSFunction {
	return func(ctx *jsexecutor.JSContext, this jsexecutor.JSValue, args []jsexecutor.JSValue) jsexecutor.JSValue {
		return ctx.ThrowInternalError("pty: does not supported current OS")
	}
}
