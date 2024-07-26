use crate::bundle::{NativeModule, SCRIPT_MODULE};
use anyhow::{anyhow, Result};
use rquickjs::{
    loader::{
        BuiltinLoader, BuiltinResolver, FileResolver, ModuleLoader, NativeLoader, ScriptLoader,
    },
    CatchResultExt, Context, Ctx, Error, Function, Module, Promise, Runtime, Value,
    qjs
};

struct RuntimeCore {
    rt: Runtime,
    ctx: Context,
}

impl RuntimeCore {
    pub fn new() -> Result<Self> {
        let rt = Runtime::new()?;
        let ctx = Context::full(&rt)?;

        let resolver = (
            BuiltinResolver::default()
                .with_module("bundle/script_module")
                .with_module("bundle/native_module"),
            FileResolver::default()
                .with_path("./")
                .with_path("../../target/debug")
                .with_native(),
        );
        let loader = (
            BuiltinLoader::default().with_module("bundle/script_module", SCRIPT_MODULE),
            ModuleLoader::default().with_module("bundle/native_module", NativeModule),
            ScriptLoader::default(),
            NativeLoader::default(),
        );
        rt.set_loader(resolver, loader);

        Ok(RuntimeCore { rt, ctx })
    }

    pub fn eval<F>(&self, name: &str, source: &str, f: F) -> Result<()>
    where
        F: FnOnce(Ctx, Promise) -> Result<()> + Send,
    {
        self.ctx.with(|ctx| {
            let ret = Module::evaluate(ctx.clone(), name, source)?;
            f(ctx, ret)
        })
    }

    pub fn eval_anonymous<F>(&self, source: &str, f: F) -> Result<()>
    where
        F: FnOnce(Ctx, Promise) -> Result<()> + Send,
    {
        self.eval("anonymous", source, f)
    }
}

#[cfg(test)]
mod tests {
    use std::ffi::{CStr, CString};
    use crate::runtime::*;
    use anyhow::Result;
    use rquickjs::qjs;

    #[test]
    fn it_works() -> Result<()> {
        let rt = RuntimeCore::new()?;
        match rt.eval_anonymous("throw new Error(123)", |ctx, p| {
            match p.finish::<Value>() {
                Ok(result) => {
                    println!("result: {:?}", result);
                }
                Err(err) => {
                    println!("error: {}", err);
                }
            }
            Ok(())
        }) {
            Ok(_) => {
                println!("done")
            }
            Err(e) => {
                eprintln!("eval error: {}", e)
            }
        };
        Ok(())
    }

    #[test]
    fn raw_qjs() -> Result<()> {
        let rt = unsafe{ qjs::JS_NewRuntime() };
        let ctx = unsafe { qjs::JS_NewContext(rt) };
        let name = CString::new("anonymous")?;
        let source:Vec<u8> = "throw new Error(123)".into();
        let len = source.len();
        let source = CString::new(source)?;
        let flag = qjs::JS_EVAL_TYPE_MODULE | qjs::JS_EVAL_FLAG_STRICT;
        let ret = unsafe { qjs::JS_Eval(ctx, source.as_ptr(), len as _, name.as_ptr(), flag as i32) };
        let is_exception = unsafe {qjs::JS_IsException(ret)};
        if is_exception {
            let err_val = unsafe { qjs::JS_GetException(ctx) };
            let err_raw = unsafe { qjs::JS_ToCString(ctx, err_val) };
            let err = unsafe { CStr::from_ptr(err_raw) };
            let err = err.to_string_lossy();
            unsafe { qjs::JS_FreeCString(ctx, err_raw) };
            unsafe { qjs::JS_FreeValue(ctx, err_val)};
            eprintln!("eval error: {}", err);
        } else {
            let state = unsafe { qjs::JS_PromiseState(ctx, ret) };
            match state {
                qjs::JSPromiseStateEnum_JS_PROMISE_FULFILLED => { println!("promise fulfilled"); },
                qjs::JSPromiseStateEnum_JS_PROMISE_REJECTED => {
                    let reason_val = unsafe { qjs::JS_PromiseResult(ctx, ret) };
                    let reason_raw = unsafe { qjs::JS_ToCString(ctx, reason_val) };
                    let reason = unsafe { CStr::from_ptr(reason_raw) };
                    let reason = reason.to_string_lossy();
                    println!("promise rejected: {}", reason);
                    unsafe { qjs::JS_FreeCString(ctx, reason_raw) };
                    unsafe { qjs::JS_FreeValue(ctx, reason_val)};
                },
                qjs::JSPromiseStateEnum_JS_PROMISE_PENDING => { println!("promise pending"); },
                _ => { println!("unknown promise state: {}", state); }
            }

            let result_raw = unsafe { qjs::JS_ToCString(ctx, ret) };
            let result = unsafe { CStr::from_ptr(result_raw) };
            let result = result.to_string_lossy();
            println!("result: {:?}", result);
            unsafe { qjs::JS_FreeCString(ctx, result_raw) };
        }
        Ok(())
    }
}
