use rquickjs::{
    embed, loader::Bundle, qjs, CatchResultExt, CaughtError, CaughtResult, Context, Module, Runtime,
};
use std::ffi::{CStr, CString};

/// load the `my_module.js` file and name it myModule
static BUNDLE: Bundle = embed! {
    "myModule": "script_module.js",
};

fn main() {
    let rt = unsafe { qjs::JS_NewRuntime() };
    let ctx = unsafe { qjs::JS_NewContext(rt) };
    let name = CString::new("anonymous").unwrap();
    let source: Vec<u8> = "throw new Error(123)".into();
    let len = source.len();
    let source = CString::new(source).unwrap();
    let flag = qjs::JS_EVAL_TYPE_MODULE | qjs::JS_EVAL_FLAG_STRICT;
    let ret = unsafe { qjs::JS_Eval(ctx, source.as_ptr(), len as _, name.as_ptr(), flag as i32) };
    let is_exception = unsafe { qjs::JS_IsException(ret) };
    if is_exception {
        let err_val = unsafe { qjs::JS_GetException(ctx) };
        let err_raw = unsafe { qjs::JS_ToCString(ctx, err_val) };
        let err = unsafe { CStr::from_ptr(err_raw) };
        let err = err.to_string_lossy();
        unsafe { qjs::JS_FreeCString(ctx, err_raw) };
        unsafe { qjs::JS_FreeValue(ctx, err_val) };
        eprintln!("eval error: {}", err);
    } else {
        let state = unsafe { qjs::JS_PromiseState(ctx, ret) };
        match state {
            qjs::JSPromiseStateEnum_JS_PROMISE_FULFILLED => {
                println!("promise fulfilled");
            }
            qjs::JSPromiseStateEnum_JS_PROMISE_REJECTED => {
                let reason_val = unsafe { qjs::JS_PromiseResult(ctx, ret) };
                let reason_raw = unsafe { qjs::JS_ToCString(ctx, reason_val) };
                let reason = unsafe { CStr::from_ptr(reason_raw) };
                let reason = reason.to_string_lossy();
                println!("promise rejected: {}", reason);
                unsafe { qjs::JS_FreeCString(ctx, reason_raw) };
                unsafe { qjs::JS_FreeValue(ctx, reason_val) };
            }
            qjs::JSPromiseStateEnum_JS_PROMISE_PENDING => {
                println!("promise pending");
            }
            _ => {
                println!("unknown promise state: {}", state);
            }
        }

        let result_raw = unsafe { qjs::JS_ToCString(ctx, ret) };
        let result = unsafe { CStr::from_ptr(result_raw) };
        let result = result.to_string_lossy();
        println!("result: {:?}", result);
        unsafe { qjs::JS_FreeCString(ctx, result_raw) };
    }
}

fn main_1() {
    let rt = Runtime::new().unwrap();
    let ctx = Context::full(&rt).unwrap();

    rt.set_loader(BUNDLE, BUNDLE);
    ctx.with(|ctx| {
        match Module::evaluate(
            ctx.clone(),
            "testModule.js",
            r#"
            // await 123;
            // throw new Error("Some error");
            import { f, foo } from 'myModule';
            // import * as a from 'myModule';
            // console.log(f());
            // if(foo() !== 2){
            //     throw new Error("Function didn't return the correct value");
            // }
        "#,
        ) {
            Ok(p) => match p.finish::<()>().catch(&ctx) {
                Ok(_) => {
                    println!("Module loaded successfully");
                }
                Err(e) => match e {
                    CaughtError::Error(e) => {
                        eprintln!("Error: {}", e.to_string());
                    }
                    CaughtError::Exception(e) => {
                        eprintln!("Exception: {}", e);
                    }
                    CaughtError::Value(e) => {
                        eprintln!("Value: {:?}", e);
                    }
                },
            },
            Err(e) => {
                eprintln!("Error: {}", e);
            }
        }
    })
}
