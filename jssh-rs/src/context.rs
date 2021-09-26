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

    self.qjs_ctx.add_callback("__builtin_op_path_join", builtin_op_path_join)?;
    self.qjs_ctx.add_callback("__builtin_op_path_abs", builtin_op_path_abs)?;
    self.qjs_ctx.add_callback("__builtin_op_path_base", builtin_op_path_base)?;
    self.qjs_ctx.add_callback("__builtin_op_path_ext", builtin_op_path_ext)?;
    self.qjs_ctx.add_callback("__builtin_op_path_dir", builtin_op_path_dir)?;

    self.qjs_ctx.add_callback("__builtin_op_http_request", builtin_op_http_request)?;

    self.qjs_ctx.eval(include_str!("runtime/js/00_jssh.js"))?;
    self.qjs_ctx.eval(include_str!("runtime/js/10_format.js"))?;
    self.qjs_ctx.eval(include_str!("runtime/js/20_assert.js"))?;
    self.qjs_ctx.eval(include_str!("runtime/js/20_cli.js"))?;
    self.qjs_ctx.eval(include_str!("runtime/js/20_log.js"))?;
    self.qjs_ctx.eval(include_str!("runtime/js/20_socket.js"))?;
    self.qjs_ctx.eval(include_str!("runtime/js/20_fs.js"))?;
    self.qjs_ctx.eval(include_str!("runtime/js/20_path.js"))?;
    self.qjs_ctx.eval(include_str!("runtime/js/20_http.js"))?;

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
  let stdin = std::io::stdin();
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

fn builtin_op_path_join(args: Arguments) -> Result<JsValue> {
  let args: Vec<_> = args.into_vec();
  let first = args.get(0).ok_or(invalid_argument_error("missing argument: path"))?;
  let first = get_string_from_js_value(first).ok_or(invalid_argument_error("invalid argument: path expected a string"))?;
  let mut path = std::path::Path::new(&first).to_path_buf();
  for item in &args[1..] {
    let item = get_string_from_js_value(item).ok_or(invalid_argument_error("invalid argument: path expected a string"))?;
    path = path.join(&item);
  }
  Ok(JsValue::String(path.to_string_lossy().to_string()))
}

fn builtin_op_path_ext(args: Arguments) -> Result<JsValue> {
  let args: Vec<_> = args.into_vec();
  let path = args.get(0).ok_or(invalid_argument_error("missing argument: path"))?;
  let path = get_string_from_js_value(path).ok_or(invalid_argument_error("invalid argument: path expected a string"))?;
  let ext = std::path::Path::new(&path)
    .extension()
    .ok_or(invalid_argument_error("invalid argument: path is not a valid file path"))?;
  Ok(JsValue::String(ext.to_string_lossy().to_string()))
}

fn builtin_op_path_base(args: Arguments) -> Result<JsValue> {
  let args: Vec<_> = args.into_vec();
  let path = args.get(0).ok_or(invalid_argument_error("missing argument: path"))?;
  let path = get_string_from_js_value(path).ok_or(invalid_argument_error("invalid argument: path expected a string"))?;
  let base = std::path::Path::new(&path)
    .file_name()
    .ok_or(invalid_argument_error("invalid argument: path is not a valid file path"))?;
  Ok(JsValue::String(base.to_string_lossy().to_string()))
}

fn builtin_op_path_dir(args: Arguments) -> Result<JsValue> {
  let args: Vec<_> = args.into_vec();
  let path = args.get(0).ok_or(invalid_argument_error("missing argument: path"))?;
  let path = get_string_from_js_value(path).ok_or(invalid_argument_error("invalid argument: path expected a string"))?;
  let dir = std::path::Path::new(&path)
    .parent()
    .ok_or(invalid_argument_error("invalid argument: path is not a valid file path"))?;
  Ok(JsValue::String(dir.to_string_lossy().to_string()))
}

fn builtin_op_path_abs(args: Arguments) -> Result<JsValue> {
  let args: Vec<_> = args.into_vec();
  let path = args.get(0).ok_or(invalid_argument_error("missing argument: path"))?;
  let path = get_string_from_js_value(path).ok_or(invalid_argument_error("invalid argument: path expected a string"))?;
  let abs = std::fs::canonicalize(&path)?;
  Ok(JsValue::String(abs.to_string_lossy().to_string()))
}

fn http_request(
  method: &str,
  url: &str,
  headers: Option<HashMap<String, Vec<String>>>,
  body: Option<String>,
  timeout: Duration,
) -> Result<reqwest::blocking::Response> {
  let client = reqwest::blocking::Client::builder().build()?;
  let mut req = client.request(reqwest::Method::from_bytes(method.as_bytes())?, url).timeout(timeout);
  if let Some(headers) = headers {
    let mut h = reqwest::header::HeaderMap::new();
    for (key, values) in headers {
      for value in values {
        h.append(reqwest::header::HeaderName::from_str(&key)?, reqwest::header::HeaderValue::from_str(&value)?);
      }
    }
    req = req.headers(h);
  }
  if let Some(body) = body {
    req = req.body(body);
  }
  log::debug!("http_request: method={} url={} timeout={:?}", method, url, timeout);
  let res = req.send()?;
  Ok(res)
}

fn builtin_op_http_request(args: Arguments) -> Result<JsValue> {
  let args: Vec<_> = args.into_vec();

  let method = args.get(0).ok_or(invalid_argument_error("missing argument: method"))?;
  let url = args.get(1).ok_or(invalid_argument_error("missing argument: url"))?;
  let headers = args.get(2).ok_or(invalid_argument_error("missing argument: headers"))?;
  let body = args.get(3).ok_or(invalid_argument_error("missing argument: body"))?;
  let timeout_ms = args.get(4).ok_or(invalid_argument_error("missing argument: timeout"))?;

  let method = get_string_from_js_value(method).ok_or(invalid_argument_error("invalid argument: method expected a string"))?;
  let url = get_string_from_js_value(url).ok_or(invalid_argument_error("invalid argument: url expected a string"))?;
  let headers = match headers {
    JsValue::Null | JsValue::Undefined => None,
    JsValue::Object(obj) => {
      let mut headers = HashMap::new();
      for (k, v) in obj {
        match v {
          JsValue::Array(arr) => {
            let mut values = Vec::new();
            for item in arr {
              let v = get_string_from_js_value(item).ok_or(invalid_argument_error("invalid argument: headers expected string key and value"))?;
              values.push(v);
            }
            headers.insert(k.clone(), values);
          }
          _ => {
            let v = get_string_from_js_value(v).ok_or(invalid_argument_error("invalid argument: headers expected string key and value"))?;
            headers.insert(k.clone(), vec![v]);
          }
        }
      }
      Some(headers)
    }
    _ => return Err(invalid_argument_error("invalid argument: headers expected an object or null")),
  };
  let body = get_string_from_js_value(body);
  let timeout_ms = get_i32_from_js_value(timeout_ms).ok_or(invalid_argument_error("invalid argument: timeout expected a number"))?;

  let res = http_request(&method, &url, headers, body, Duration::from_millis(timeout_ms as u64))?;
  let mut result = HashMap::new();
  result.insert("status".to_string(), JsValue::Int(res.status().as_u16() as i32));
  let mut headers = HashMap::new();
  for (k, v) in res.headers() {
    let name = k.to_string();
    match headers.get(&name) {
      Some(old_value) => match old_value {
        JsValue::String(s) => {
          headers.insert(name, JsValue::Array(vec![JsValue::String(s.clone()), JsValue::String(v.to_str()?.to_string())]));
        }
        JsValue::Array(arr) => {
          let mut arr = arr.clone();
          arr.push(JsValue::String(v.to_str()?.to_string()));
          headers.insert(name, JsValue::Array(arr));
        }
        _ => {
          headers.insert(name, JsValue::String(v.to_str()?.to_string()));
        }
      },
      None => {
        headers.insert(name, JsValue::String(v.to_str()?.to_string()));
      }
    }
  }
  result.insert("headers".to_string(), JsValue::Object(headers));
  result.insert("body".to_string(), JsValue::String(String::from_utf8_lossy(res.bytes()?.as_ref()).to_string()));

  Ok(JsValue::Object(result))
}
