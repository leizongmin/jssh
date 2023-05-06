http.get = function (url, headers = {}) {
  return http.request("GET", url, headers);
};

http.head = function (url, headers = {}) {
  return http.request("HEAD", url, headers);
};

http.options = function (url, headers = {}) {
  return http.request("OPTIONS", url, headers);
};

http.post = function (url, headers = {}, body = "") {
  return http.request("POST", url, headers, body);
};

http.put = function (url, headers = {}, body = "") {
  return http.request("PUT", url, headers, body);
};

http.delete = function (url, headers = {}, body = "") {
  return http.request("DELETE", url, headers, body);
};

Object.freeze(http);
