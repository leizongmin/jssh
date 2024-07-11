use std::path::PathBuf;

use clap::{Parser, Subcommand};

/// Welcome to jssh, a lightweight JavaScript runtime.
#[derive(Parser)]
#[command(version, about, long_about = None)]
struct Cli {
    /// Script file to run
    script_file: Option<PathBuf>,
    #[command(subcommand)]
    command: Option<Commands>,
    /// Arguments for script
    args: Vec<String>,
}

#[derive(Subcommand)]
enum Commands {
    /// Create self-contained binary file
    Build {
        script_file: PathBuf,
        target_file: Option<PathBuf>,
    },
    /// Start REPL
    REPL {
        script_file: Option<PathBuf>,
    },
    /// Run script from argument and print the result
    Eval {
        script_code: String,
    },
    /// Run script from argument
    Exec {
        script_code: String,
    },
}

fn main() {
    let cli = Cli::parse();
    match cli.command {
        Some(Commands::Build { script_file, target_file }) => {
            println!("Build: {:?} {:?}", script_file, target_file);
        }
        Some(Commands::REPL { script_file }) => {
            println!("REPL: {:?} {:?}", script_file, cli.args);
        }
        Some(Commands::Eval { script_code }) => {
            println!("Eval: {:?} {:?}", script_code, cli.args);
        }
        Some(Commands::Exec { script_code }) => {
            println!("Exec: {:?} {:?}", script_code, cli.args);
        }
        None => {
            match cli.script_file {
                Some(script_file) => {
                    println!("Run: {:?} {:?}", script_file, cli.args);
                }
                None => {
                    println!("REPL");
                }
            }
        }
    }
}
