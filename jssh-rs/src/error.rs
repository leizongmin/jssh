use std::borrow::Cow;
use std::error::Error;
use std::fmt;
use std::fmt::{Display, Formatter};

use quick_js::ExecutionError;

/// A generic wrapper that can encapsulate any concrete error type.
pub type AnyError = anyhow::Error;

/// Creates a new error with a caller-specified error class name and message.
pub fn custom_error(class: &'static str, message: impl Into<Cow<'static, str>>) -> AnyError {
  CustomError {
    class,
    message: message.into(),
  }
  .into()
}

pub fn generic_error(message: impl Into<Cow<'static, str>>) -> AnyError {
  custom_error("Error", message)
}

pub fn execution_error(err: ExecutionError) -> AnyError {
  custom_error("ExecutionError", err.to_string())
}

pub fn uri_error(message: impl Into<Cow<'static, str>>) -> AnyError {
  custom_error("URIError", message)
}

/// A simple error type that lets the creator specify both the error message and
/// the error class name. This type is private; externally it only ever appears
/// wrapped in an `AnyError`. To retrieve the error class name from a wrapped
/// `CustomError`, use the function `get_custom_error_class()`.
#[derive(Debug)]
struct CustomError {
  class: &'static str,
  message: Cow<'static, str>,
}

impl Display for CustomError {
  fn fmt(&self, f: &mut Formatter) -> fmt::Result {
    f.write_str((&self.message).as_ref())
  }
}

impl Error for CustomError {}
