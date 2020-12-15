package jsshcmd

func jsFnShPty(global typeutil.H) jsexecutor.JSFunction {
	return func(ctx *jsexecutor.JSContext, this jsexecutor.JSValue, args []jsexecutor.JSValue) jsexecutor.JSValue {
		return ctx.ThrowInternalError("pty: does not supported current OS")
	}
}
