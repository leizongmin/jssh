{
  log.fatal = function (message, ...args) {
    log.error(message, ...args);
    exit(1);
  };
}
