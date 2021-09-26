const http = {};
{
  let defaultTimeout = 60000;
  http.timeout = function (ms = defaultTimeout) {
    if (ms > 0) {
      defaultTimeout = ms;
      return defaultTimeout;
    } else {
      throw new TypeError(`invalid argument: expected a number greater than 0`);
    }
  };

  http.request = function (
    method,
    url,
    headers = null,
    body = null,
    timeout = defaultTimeout
  ) {
    return jssh.op.http_request(method, url, headers, body, timeout);
  };

  http.get = function (
    url,
    headers = null,
    body = null,
    timeout = defaultTimeout
  ) {
    return jssh.op.http_request("GET", url, headers, body, timeout);
  };

  http.head = function (
    url,
    headers = null,
    body = null,
    timeout = defaultTimeout
  ) {
    return jssh.op.http_request("HEAD", url, headers, body, timeout);
  };

  http.post = function (
    url,
    headers = null,
    body = null,
    timeout = defaultTimeout
  ) {
    return jssh.op.http_request("POST", url, headers, body, timeout);
  };

  http.put = function (
    url,
    headers = null,
    body = null,
    timeout = defaultTimeout
  ) {
    return jssh.op.http_request("PUT", url, headers, body, timeout);
  };

  http.delete = function (
    url,
    headers = null,
    body = null,
    timeout = defaultTimeout
  ) {
    return jssh.op.http_request("DELETE", url, headers, body, timeout);
  };

  http.options = function (
    url,
    headers = null,
    body = null,
    timeout = defaultTimeout
  ) {
    return jssh.op.http_request("OPTIONS", url, headers, body, timeout);
  };

  http.trace = function (
    url,
    headers = null,
    body = null,
    timeout = defaultTimeout
  ) {
    return jssh.op.http_request("TRACE", url, headers, body, timeout);
  };

  http.download = function (
    url,
    saveToPath = path.join(
      __downloaddir,
      `jssh-http-download-${Date.now()}-${randomstring(6)}`
    )
  ) {
    const res = http.get(url);
    fs.writefile(saveToPath, res.body);
    return saveToPath;
  };
}
Object.freeze(http);
