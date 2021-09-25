const socket = {};
{
  let defaultTimeout = 60000;
  socket.timeout = function (ms = defaultTimeout) {
    if (ms > 0) {
      defaultTimeout = ms;
      return defaultTimeout;
    } else {
      throw new TypeError(`invalid argument: expected a number greater than 0`);
    }
  };

  socket.tcpsend = function (host, port, data, timeout = defaultTimeout) {
    return jssh.op.tcp_send(host, port, data, timeout);
  };

  socket.tcptest = function (host, port, timeout = defaultTimeout) {
    return jssh.op.tcp_test(host, port, timeout);
  };
}
Object.freeze(log);
