{
  function println(...args) {
    jssh.print(...args);
    jssh.print("\n");
  }

  function eprintln(...args) {
    jssh.eprint(...args);
    jssh.eprint("\n");
  }

  const log = {};
  {
    const levels = { ERROR: 1, WARN: 2, INFO: 3, DEBUG: 4 };
    const logLevel = levels[(jssh.__env["JSSH_LOG"] || "INFO").toUpperCase()];

    const reset = `\u001b[0m`;

    const red = (line) => {
      return `\u001b[31;1m${line}${reset}`;
    };

    const green = (line) => {
      return `\u001b[32;1m${line}${reset}`;
    };

    const yellow = (line) => {
      return `\u001b[33;1m${line}${reset}`;
    };

    const gray = (line) => {
      return `\u001b[2;1m${line}${reset}`;
    };

    log.debug = function debug(message, ...args) {
      if (logLevel >= levels.DEBUG) {
        jssh.stdoutlog(gray(jssh.format(message, ...args)));
      }
    };

    log.info = function info(message, ...args) {
      if (logLevel >= levels.INFO) {
        jssh.stdoutlog(green(jssh.format(message, ...args)));
      }
    };

    log.warn = function error(message, ...args) {
      if (logLevel >= levels.WARN) {
        jssh.stderrlog(yellow(jssh.format(message, ...args)));
      }
    };

    log.error = function error(message, ...args) {
      if (logLevel >= levels.ERROR) {
        jssh.stderrlog(red(jssh.format(message, ...args)));
      }
    };

    log.fatal = function fatal(message, ...args) {
      log.error(message, ...args);
      jssh.exit(1);
    };
  }

  jssh.println = println;
  jssh.eprintln = eprintln;
  jssh.log = log;
  Object.freeze(log);
}
