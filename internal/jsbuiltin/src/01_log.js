function println(...args) {
  print(...args);
  print("\n");
}

const log = {};
{
  const levels = { ERROR: 1, WARN: 2, INFO: 3, DEBUG: 4 };
  const logLevel = levels[(__env["JSSH_LOG"] || "INFO").toUpperCase()];

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
      stdoutlog(gray(format(message, ...args)));
    }
  };

  log.info = function info(message, ...args) {
    if (logLevel >= levels.INFO) {
      stdoutlog(green(format(message, ...args)));
    }
  };

  log.warn = function error(message, ...args) {
    if (logLevel >= levels.WARN) {
      stderrlog(yellow(format(message, ...args)));
    }
  };

  log.error = function error(message, ...args) {
    if (logLevel >= levels.ERROR) {
      stderrlog(red(format(message, ...args)));
    }
  };

  log.fatal = function fatal(message, ...args) {
    log.error(message, ...args);
    exit(1);
  };
}

Object.freeze(log);
