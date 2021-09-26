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

exit(1);
