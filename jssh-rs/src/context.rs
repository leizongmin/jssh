use std::collections::HashMap;
use std::fs::{DirEntry, Metadata};
use std::io::{Read, Write};
use std::net::{SocketAddr, TcpStream};
use std::os::unix::fs::MetadataExt;
use std::str::FromStr;
use std::time::Duration;

use anyhow::Result;
use quick_js::console::LogConsole;
use quick_js::{Arguments, Context, JsValue};

use crate::error::{execution_error, generic_error, invalid_argument_error, system_error};

pub struct JsContext {
  qjs_ctx: Context,
}

impl JsContext {
  pub fn new() -> Result<JsContext> {
    let qjs_ctx = Context::builder().console(LogConsole {}).build().map_err(|e| generic_error(e.to_string()))?;
    Ok(JsContext { qjs_ctx })
  }

  pub fn load_std(&self) -> Result<()> {
    self.qjs_ctx.add_callback("__builtin_op_stdout_write", builtin_op_stdout_write)?;
    self.qjs_ctx.add_callback("__builtin_op_stderr_write", builtin_op_stderr_write)?;
    self.qjs_ctx.add_callback("__builtin_op_stdin_read_line", builtin_op_stdin_read_line)?;
    self.qjs_ctx.add_callback("__builtin_op_exit", builtin_op_exit)?;
    self.qjs_ctx.add_callback("__builtin_op_env", builtin_op_env)?;
    self.qjs_ctx.add_callback("__builtin_op_args", builtin_op_args)?;

    self.qjs_ctx.add_callback("__builtin_op_tcp_send", builtin_op_tcp_send)?;
    self.qjs_ctx.add_callback("__builtin_op_tcp_test", builtin_op_tcp_test)?;

    self.qjs_ctx.add_callback("__builtin_op_file_read", builtin_op_file_read)?;
    self.qjs_ctx.add_callback("__builtin_op_dir_read", builtin_op_dir_read)?;
    self.qjs_ctx.add_callback("__builtin_op_file_write", builtin_op_file_write)?;
    self.qjs_ctx.add_callback("__builtin_op_file_append", builtin_op_file_append)?;
    self.qjs_ctx.add_callback("__builtin_op_file_stat", builtin_op_file_stat)?;
    self.qjs_ctx.add_callback("__builtin_op_file_exist", builtin_op_file_exist)?;

    self.qjs_ctx.eval(include_str!("runtime/js/00_jssh.js"))?;
    self.qjs_ctx.eval(include_str!("runtime/js/10_format.js"))?;
    self.qjs_ctx.eval(include_str!("runtime/js/20_assert.js"))?;
    self.qjs_ctx.eval(include_str!("runtime/js/20_cli.js"))?;
    self.qjs_ctx.eval(include_str!("runtime/js/20_log.js"))?;
    self.qjs_ctx.eval(include_str!("runtime/js/20_socket.js"))?;
    self.qjs_ctx.eval(include_str!("runtime/js/20_fs.js"))?;

    Ok(())
  }

  pub fn eval(&self, code: &str) -> Result<JsValue> {
    Ok(self.qjs_ctx.eval(code).map_err(|e| execution_error(e))?)
  }
}

fn builtin_op_stdout_write(args: Arguments) -> Result<JsValue> {
  let mut stdout = std::io::stdout();
  for a in args.into_vec() {
    if let Some(s) = a.as_str() {
      write!(stdout, "{}", s)?;
    } else {
      write!(stdout, "{:?}", a)?;
    }
  }
  stdout.flush()?;
  Ok(JsValue::Undefined)
}

fn builtin_op_stderr_write(args: Arguments) -> Result<JsValue> {
  let mut stderr = std::io::stderr();
  for a in args.into_vec() {
    if let Some(s) = a.as_str() {
      write!(stderr, "{}", s)?;
    } else {
      write!(stderr, "{:?}", a)?;
    }
  }
  stderr.flush()?;
  Ok(JsValue::Undefined)
}

fn builtin_op_stdin_read_line(_args: Arguments) -> Result<JsValue> {
  let mut line = String::new();
  let mut stdin = std::io::stdin();
  stdin.read_line(&mut line)?;
  Ok(JsValue::String(line))
}

fn builtin_op_exit(args: Arguments) {
  if let Some(a) = args.into_vec().get(0) {
    if let Some(code) = get_i32_from_js_value(a) {
      std::process::exit(code);
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

fn builtin_op_args(_args: Arguments) -> JsValue {
  JsValue::Array(std::env::args().map(JsValue::String).collect())
}

fn tcp_send(addr: String, data: Option<&[u8]>, timeout: Duration) -> Result<Vec<u8>> {
  let addr = SocketAddr::from_str(&addr)?;
  let mut stream = TcpStream::connect_timeout(&addr, timeout)?;
  stream.set_nodelay(true)?;
  stream.set_read_timeout(Some(timeout))?;
  stream.set_write_timeout(Some(timeout))?;
  if let Some(data) = data {
    stream.write(data)?;
  }
  let mut output = Vec::new();
  stream.read_to_end(&mut output)?;
  Ok(output)
}

fn builtin_op_tcp_send(args: Arguments) -> Result<JsValue> {
  let args = args.into_vec();
  let host = args.get(0).ok_or(invalid_argument_error("missing argument: host"))?;
  let port = args.get(1).ok_or(invalid_argument_error("missing argument: port"))?;
  let data = args.get(2).ok_or(invalid_argument_error("missing argument: data"))?;
  let timeout_ms = args.get(3).ok_or(invalid_argument_error("missing argument: timeout"))?;

  let host = get_string_from_js_value(host).ok_or(invalid_argument_error("invalid argument: host expected a string"))?;
  let port = get_i32_from_js_value(port).ok_or(invalid_argument_error("invalid argument: host expected a number"))?;
  let data = get_string_from_js_value(data).ok_or(invalid_argument_error("invalid argument: data expected a string"))?;
  let timeout_ms = get_i32_from_js_value(timeout_ms).ok_or(invalid_argument_error("invalid argument: timeout expected a number"))?;

  let output = tcp_send(format!("{}:{}", host, port), Some(data.as_bytes()), Duration::from_millis(timeout_ms as u64))?;
  let output = String::from_utf8_lossy(output.as_slice());
  Ok(JsValue::String(output.to_string()))
}

fn tcp_test(addr: String, timeout: Duration) -> Result<()> {
  let addr = SocketAddr::from_str(&addr)?;
  TcpStream::connect_timeout(&addr, timeout)?;
  Ok(())
}

fn builtin_op_tcp_test(args: Arguments) -> Result<JsValue> {
  let args = args.into_vec();
  let host = args.get(0).ok_or(invalid_argument_error("missing argument: host"))?;
  let port = args.get(1).ok_or(invalid_argument_error("missing argument: port"))?;
  let timeout_ms = args.get(2).ok_or(invalid_argument_error("missing argument: timeout"))?;

  let host = get_string_from_js_value(host).ok_or(invalid_argument_error("invalid argument: host expected a string"))?;
  let port = get_i32_from_js_value(port).ok_or(invalid_argument_error("invalid argument: host expected a number"))?;
  let timeout_ms = get_i32_from_js_value(timeout_ms).ok_or(invalid_argument_error("invalid argument: timeout expected a number"))?;

  let ok = tcp_test(format!("{}:{}", host, port), Duration::from_millis(timeout_ms as u64))
    .map(|_| true)
    .unwrap_or(false);
  Ok(JsValue::Bool(ok))
}

fn builtin_op_file_read(args: Arguments) -> Result<JsValue> {
  let args = args.into_vec();
  let path = args.get(0).ok_or(invalid_argument_error("missing argument: path"))?;
  let path = get_string_from_js_value(path).ok_or(invalid_argument_error("invalid argument: path expected a string"))?;
  let data = std::fs::read_to_string(&path)?;
  Ok(JsValue::String(data))
}

fn file_metadata_to_js_value(meta: &Metadata, file_name: &str) -> Result<JsValue> {
  let mut item = HashMap::new();
  let file_type = meta.file_type();
  item.insert("file_name".to_string(), JsValue::String(file_name.to_string()));
  item.insert("is_file".to_string(), JsValue::Bool(file_type.is_file()));
  item.insert("is_dir".to_string(), JsValue::Bool(file_type.is_dir()));
  item.insert("is_symlink".to_string(), JsValue::Bool(file_type.is_symlink()));
  item.insert("size".to_string(), JsValue::Int(meta.size() as i32));
  item.insert("atime".to_string(), JsValue::Int(meta.atime() as i32));
  item.insert("ctime".to_string(), JsValue::Int(meta.ctime() as i32));
  item.insert("mtime".to_string(), JsValue::Int(meta.mtime() as i32));
  Ok(JsValue::Object(item))
}

fn builtin_op_dir_read(args: Arguments) -> Result<JsValue> {
  let args = args.into_vec();
  let path = args.get(0).ok_or(invalid_argument_error("missing argument: path"))?;
  let path = get_string_from_js_value(path).ok_or(invalid_argument_error("invalid argument: path expected a string"))?;
  let entries: Vec<_> = std::fs::read_dir(&path)?.collect();
  let mut list = Vec::new();
  for entry in entries {
    let entry = entry?;
    let meta = entry.metadata()?;
    let item = file_metadata_to_js_value(&meta, entry.file_name().to_str().ok_or(system_error("cannot get file_name"))?)?;
    list.push(item);
  }
  Ok(JsValue::Array(list))
}

fn builtin_op_file_write(args: Arguments) -> Result<JsValue> {
  let args = args.into_vec();
  let path = args.get(0).ok_or(invalid_argument_error("missing argument: path"))?;
  let data = args.get(1).ok_or(invalid_argument_error("missing argument: data"))?;

  let path = get_string_from_js_value(path).ok_or(invalid_argument_error("invalid argument: path expected a string"))?;
  let data = get_string_from_js_value(data).ok_or(invalid_argument_error("invalid argument: data expected a string"))?;

  std::fs::write(&path, data.as_bytes())?;
  Ok(JsValue::Undefined)
}

fn builtin_op_file_append(args: Arguments) -> Result<JsValue> {
  let args = args.into_vec();
  let path = args.get(0).ok_or(invalid_argument_error("missing argument: path"))?;
  let data = args.get(1).ok_or(invalid_argument_error("missing argument: data"))?;

  let path = get_string_from_js_value(path).ok_or(invalid_argument_error("invalid argument: path expected a string"))?;
  let data = get_string_from_js_value(data).ok_or(invalid_argument_error("invalid argument: data expected a string"))?;

  // std::fs::(&path, data.as_bytes())?;
  let mut file = std::fs::OpenOptions::new().write(true).append(true).open(&path)?;
  file.write_all(data.as_bytes())?;
  Ok(JsValue::Undefined)
}

fn builtin_op_file_stat(args: Arguments) -> Result<JsValue> {
  let args = args.into_vec();
  let path = args.get(0).ok_or(invalid_argument_error("missing argument: path"))?;
  let path = get_string_from_js_value(path).ok_or(invalid_argument_error("invalid argument: path expected a string"))?;
  let meta = std::fs::metadata(&path)?;
  let data = file_metadata_to_js_value(&meta, &path)?;
  Ok(data)
}

fn builtin_op_file_exist(args: Arguments) -> Result<JsValue> {
  let args = args.into_vec();
  let path = args.get(0).ok_or(invalid_argument_error("missing argument: path"))?;
  let path = get_string_from_js_value(path).ok_or(invalid_argument_error("invalid argument: path expected a string"))?;
  let exist = std::fs::try_exists(&path)?;
  Ok(JsValue::Bool(exist))
}
