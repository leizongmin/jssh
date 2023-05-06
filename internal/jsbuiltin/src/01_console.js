const console = {};

{
  console.trace = function trace(...args) {
    println(format(...args));
  };

  console.log = function log(...args) {
    println(format(...args));
  };

  console.info = function info(...args) {
    println(format(...args));
  };

  console.warn = function warn(...args) {
    eprintln(format(...args));
  };

  console.error = function error(...args) {
    eprintln(format(...args));
  };
}

Object.freeze(console);
