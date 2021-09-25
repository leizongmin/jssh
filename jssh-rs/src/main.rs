use anyhow::Result;
use clap::{App, Arg};
use quick_js::JsValue;

use crate::context::JsContext;
use crate::error::uri_error;
use simplelog::{ColorChoice, CombinedLogger, Config, LevelFilter, TermLogger, TerminalMode};

mod context;
mod error;

fn main() {
    CombinedLogger::init(vec![TermLogger::new(
        LevelFilter::Trace,
        Config::default(),
        TerminalMode::Mixed,
        ColorChoice::Auto,
    )])
    .unwrap();

    let mut app = App::new("jssh")
        .version("0.0.0-alpha")
        .author("Zongmin Lei <leizongmin@gmail.com>")
        .about("A tiny JavaScript runtime")
        .subcommand(
            App::new("run").about("Run script file").arg(
                Arg::new("file")
                    .required(true)
                    .about("local file path or URL, e.g. file:///path/to/file.js, http://example.com/file.js"),
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
                eval_js_file(file.to_string()).unwrap();
            } else {
                app.print_long_help().unwrap();
            }
        }
        Err(e) => {
            e.exit();
        }
    }
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
