use std::collections::HashMap;
use std::io::Write;
use std::ops::Deref;

use anyhow::Result;
use quick_js::console::LogConsole;
use quick_js::{Arguments, Context, JsValue};

use crate::error::{execution_error, generic_error};

pub struct JsContext {
  qjs_ctx: Context,
}

impl JsContext {
  pub fn new() -> Result<JsContext> {
    let context = Context::builder().console(LogConsole {}).build().map_err(|e| generic_error(e.to_string()))?;
    Ok(JsContext { qjs_ctx: context })
  }

  pub fn load_std(&self) -> Result<()> {
    self.qjs_ctx.add_callback("__builtin_op_stdout_write", builtin_op_stdout_write)?;
    self.qjs_ctx.add_callback("__builtin_op_stderr_write", builtin_op_stderr_write)?;
    self.qjs_ctx.add_callback("__builtin_op_exit", builtin_op_exit)?;
    self.qjs_ctx.add_callback("__builtin_op_env", builtin_op_env)?;

    self.qjs_ctx.eval(include_str!("runtime/js/00_jssh.js"))?;
    self.qjs_ctx.eval(include_str!("runtime/js/10_format.js"))?;
    self.qjs_ctx.eval(include_str!("runtime/js/20_log.js"))?;

    Ok(())
  }

  pub fn eval(&self, code: &str) -> Result<JsValue> {
    Ok(self.qjs_ctx.eval(code).map_err(|e| execution_error(e))?)
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

fn builtin_op_stderr_write(args: Arguments) {
  let mut stderr = std::io::stderr();
  for a in args.into_vec() {
    if let Some(s) = a.as_str() {
      write!(stderr, "{}", s);
    } else {
      write!(stderr, "{:?}", a);
    }
  }
}

fn builtin_op_exit(args: Arguments) {
  if let Some(a) = args.into_vec().get(0) {
    if let Some(code) = get_i32_from_js_value(a) {
      return std::process::exit(code);
    }
  }
  std::process::exit(0);
}

fn get_i32_from_js_value(v: &JsValue) -> Option<i32> {
  match v {
    JsValue::Int(v) => Some(*v),
    JsValue::Float(v) => Some(*v as i32),
    JsValue::BigInt(v) => v.as_i64().map(|v| v as i32).or_else(|| None),
    _ => None,
  }
}

fn get_string_from_js_value(v: &JsValue) -> Option<String> {
  match v {
    JsValue::String(s) => Some(s.clone()),
    JsValue::Int(v) => Some(format!("{}", v)),
    JsValue::Float(v) => Some(format!("{}", v)),
    JsValue::BigInt(v) => v.as_i64().map(|v| format!("{}", v)).or_else(|| None),
    _ => None,
  }
}

fn builtin_op_env(_args: Arguments) -> JsValue {
  let mut env = HashMap::new();
  for (k, v) in std::env::vars() {
    env.insert(k, JsValue::String(v));
  }
  JsValue::Object(env)
}
