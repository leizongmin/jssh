{
  jssh.http.get = function (url, headers = {}) {
    return jssh.http.request("GET", url, headers);
  };

  jssh.http.head = function (url, headers = {}) {
    return jssh.http.request("HEAD", url, headers);
  };

  jssh.http.options = function (url, headers = {}) {
    return jssh.http.request("OPTIONS", url, headers);
  };

  jssh.http.post = function (url, headers = {}, body = "") {
    return jssh.http.request("POST", url, headers, body);
  };

  jssh.http.put = function (url, headers = {}, body = "") {
    return jssh.http.request("PUT", url, headers, body);
  };

  jssh.http.delete = function (url, headers = {}, body = "") {
    return jssh.http.request("DELETE", url, headers, body);
  };

  Object.freeze(jssh.http);
}
