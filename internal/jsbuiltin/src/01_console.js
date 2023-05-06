{
  const console = {};

  console.trace = function trace(...args) {
    jssh.println(jssh.format(...args));
  };

  console.log = function log(...args) {
    jssh.println(jssh.format(...args));
  };

  console.info = function info(...args) {
    jssh.println(jssh.format(...args));
  };

  console.warn = function warn(...args) {
    jssh.eprintln(jssh.format(...args));
  };

  console.error = function error(...args) {
    jssh.eprintln(jssh.format(...args));
  };

  jssh.console = console;
  Object.freeze(console);
}
