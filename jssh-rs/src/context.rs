use quick_js::console::LogConsole;
use quick_js::{Arguments, Context, JsValue};

use crate::error::{execution_error, generic_error, AnyError};

pub struct JsContext {
    quick_js_context: Context,
}

impl JsContext {
    pub fn new() -> Result<JsContext, AnyError> {
        let context = Context::builder()
            .console(LogConsole {})
            .build()
            .map_err(|e| generic_error(e.to_string()))?;
        Ok(JsContext {
            quick_js_context: context,
        })
    }

    pub fn load_std(&self) -> Result<(), AnyError> {
        self.quick_js_context
            .add_callback("__builtin_op_stdout_write", builtin_op_stdout_write);
        self.quick_js_context.eval(include_str!("runtime/js/core.js"));
        Ok(())
    }

    pub fn eval(&self, code: &str) -> Result<JsValue, AnyError> {
        Ok(self.quick_js_context.eval(code).map_err(|e| execution_error(e))?)
    }
}

fn builtin_op_stdout_write(args: Arguments) {
    for a in args.into_vec() {
        if let Some(s) = a.as_str() {
            print!("{}", s)
        } else {
            print!("{:?}", a)
        }
    }
}
