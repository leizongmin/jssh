#![feature(path_try_exists)]

use anyhow::Result;
use clap::{App, Arg};
use quick_js::JsValue;

use crate::context::JsContext;
use crate::error::uri_error;

mod context;
mod error;

fn main() -> Result<()> {
  env_logger::builder().filter_level(log::LevelFilter::Trace).init();

  let mut app = App::new("jssh")
    .version("0.0.0-alpha")
    .author("Zongmin Lei <leizongmin@gmail.com>")
    .about("A tiny JavaScript runtime")
    .subcommand(
      App::new("run").about("Run script file").arg(
        Arg::new("file")
          .required(true)
          .about("local file path or URL, e.g. file:///path/to/file.js, https://example.com/file.js"),
      ),
    )
    .subcommand(App::new("exec").about("Run script from argument"))
    .subcommand(App::new("eval").about("Run script from argument and print the result"))
    .subcommand(App::new("build").about("Create self-contained binary file"))
    .subcommand(App::new("repl").about("Start REPL"))
    .subcommand(App::new("version").about("Show version"))
    .subcommand(App::new("help").about("Show usage"));

  match app.clone().try_get_matches() {
    Ok(matches) => {
      if let Some(_) = matches.subcommand_matches("version") {
        println!("{}", app.render_long_version());
      } else if let Some(args) = matches.subcommand_matches("run") {
        let file = args.value_of("file").expect("missing script file");
        println!("Run script file: {}", file);
        eval_js_file(file.to_string())?;
      } else {
        app.print_long_help()?;
      }
    }
    Err(e) => {
      e.exit();
    }
  };

  // let module_wat = r#"
  //   (module
  //   (type $t0 (func (param i32) (result i32)))
  //   (func $add_one (export "add_one") (type $t0) (param $p0 i32) (result i32)
  //       get_local $p0
  //       i32.const 1
  //       i32.add))
  //   "#;
  //
  // let store = wasmer::Store::default();
  // let module = wasmer::Module::new(&store, &module_wat)?;
  // // The module doesn't import anything, so we create an empty import object.
  // let import_object = wasmer::imports! {};
  // let instance = wasmer::Instance::new(&module, &import_object)?;
  //
  // let add_one = instance.exports.get_function("add_one")?;
  // let result = add_one.call(&[wasmer::Value::I32(42)])?;
  // assert_eq!(result[0], wasmer::Value::I32(43));

  Ok(())
}

fn eval_js(code: String) -> Result<JsValue> {
  let context = JsContext::new()?;
  context.load_std()?;
  context.eval(&code)
}

fn eval_js_file(file: String) -> Result<JsValue> {
  let code = std::fs::read_to_string(&file).map_err(|e| uri_error(format!("{}: {}", e, &file)))?;
  eval_js(code)
}
