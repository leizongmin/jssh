println(Object.keys(globalThis));
println(Object.keys(jssh.op));
println("hello world");

const d = new Date();
println(d.toString());
console.log(globalThis);
console.log(__env);

log.debug("this is debug message");
log.info("this is info message");
log.warn("this is warn message");
log.error("this is error message");
log.info("hello, %s, %d", "world", 123);

console.log("env=", __env);
console.log("args=", __args);

console.log(socket.timeout());
console.log(socket.timeout(5000));
// console.log(socket.tcptest("123.151.137.18", 81));
console.log(socket.tcptest("123.151.137.18", 80));
console.log(
  socket.tcpsend(
    "123.151.137.18",
    80,
    "GET / HTTP/1.1\r\nHost: qq.com\r\nConnection: Close\r\n\r\n"
  )
);

try {
  assert(true, "test aaa");
  assert(false, "test bbb");
} catch (err) {
  console.log(err.message, err.stack);
}
{
  console.log(cli.args(), cli.opts());
  // const name = cli.prompt("what's your name? ");
  // console.log("Your name is", name);
}

console.log(fs.readdir("."));
console.log(fs.readfile("build.sh"));
console.log(fs.stat("."), fs.stat("build.sh"));
console.log(fs.exist("."), fs.exist("aaaaa"));
// console.log(fs.writefile("tmp.txt", "123"), fs.readfile("tmp.txt"));
// console.log(fs.writefile("tmp.txt", "456"), fs.readfile("tmp.txt"));
// console.log(fs.appendfile("tmp.txt", "aaa"), fs.readfile("tmp.txt"));

console.log(
  path.join("a"),
  path.join("a", "b", "c"),
  path.join("a", "/b", "./c/d"),
  path.join("https://example.com", "./html", "./index.html")
);
console.log(path.dir("/a/b/c"), path.dir("a/c"));
console.log(path.base("/a/b/c"), path.base("a/c"));
console.log(path.ext("/a/b/c.j"), path.ext("a/c.xx"));
console.log(path.abs("build.sh"));

console.log(http.get("https://example.com/", { a: 123, b: ["x", "b"] }));

console.log(http.download("https://example.com"));
console.log(__tmpdir, __homedir, __downloaddir);

console.log(formatdate("Y-m-d H:i:s"));
console.log(randomstring(10, "0123456789"));
console.log(deepmerge({ a: 123 }, { b: 456 }));

console.log(global);
console.log(sleep(1000));

exit(1);
