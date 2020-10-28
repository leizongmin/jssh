function println(...args) {
  print(...args);
  print("\n");
}

const log = {};
{
  const levels = { ERROR: 1, INFO: 2, DEBUG: 3 };
  const logLevel = levels[(__env["JSSH_LOG"] || "INFO").toUpperCase()];

  const reset = `\u001b[0m`;

  function red(line) {
    return `\u001b[31;1m${line}${reset}`;
  }

  function green(line) {
    return `\u001b[32;1m${line}${reset}`;
  }

  function gray(line) {
    return `\u001b[2;1m${line}${reset}`;
  }

  log.debug = function (message, ...args) {
    if (logLevel >= levels.DEBUG) {
      stdoutlog(gray(format(message, ...args)));
    }
  };

  log.info = function (message, ...args) {
    if (logLevel >= levels.INFO) {
      stdoutlog(green(format(message, ...args)));
    }
  };

  log.error = function (message, ...args) {
    if (logLevel >= levels.ERROR) {
      stderrlog(red(format(message, ...args)));
    }
  };

  log.fatal = function (message, ...args) {
    log.error(message, ...args);
    exit(1);
  };
}
