#!/usr/bin/env go run github.com/leizongmin/jssh

log.info("aaa: %s %s %s", __bin, __dirname, __filename);
log.error(JSON.stringify(__args));
println(JSON.stringify(__env));
println("%f", Date.now());

sleep(500);
setenv("__xx__", new Date().toString());
log.info(JSON.stringify(__env));

log.info(pwd());
log.info(cd(__dirname));
log.info(cwd());
log.warn(cwd());

exec("pwd");
log.info("%f %f %s", __code, __outputbytes, __output);

exec2("pwd");
log.info("%f %f %s", __code, __outputbytes, __output);

if (exec1(`ls -al ${__homedir}`).code === 0) {
  __output.split("\n").forEach((line) => log.error(line));
}

fs.readdir(__homedir).forEach((f) => log.error(JSON.stringify(f)));
log.info(fs.readfile(`${__homedir}/.gitconfig`));
log.info(JSON.stringify(fs.stat(`${__homedir}/.gitconfig`)));

global.xyz = 12345;
log.info("xyz = %f", xyz);

const file = `${__tmpdir}/${Date.now()}-${Math.random()}.txt`;
log.info(fs.writefile(file, "hello"));
log.info(fs.appendfile(file, "world"));
log.info(fs.readfile(file));

log.info(path.abs("."));
log.info(path.base(file));
log.info(path.dir(file));
log.info(path.ext(file));
log.info(path.join("a", "b", "c"));
log.info(path.abs(path.join("a", "b", "c")));

log.info(
  "%s, %s, %v, %v, %s, %s",
  cli.get(0),
  cli.get("n"),
  cli.bool("n"),
  cli.bool("x"),
  JSON.stringify(cli.args()),
  JSON.stringify(cli.opts())
);

log.error(new Error().stack);

if (cli.bool("request")) {
  log.info(JSON.stringify(http.request("GET", "http://baidu.com")));
  log.info(format("%s-%s", "aaa", "bbb"));
}

if (cli.bool("bgexec")) {
  log.info("bgexec: pid=%v", bgexec("ping qq.com -c 60"));
  log.info("bgexec: pid=%v", bgexec("ping baidu.com -c 60"));
  log.info("tail: %s", JSON.stringify(exec1(`tail ${__filename}`)));
  sleep(3000);
}

log.info(JSON.stringify(loadconfig("config.json")));
log.info(JSON.stringify(loadconfig("config.toml")));
log.info(JSON.stringify(loadconfig("config.yaml")));
log.info(JSON.stringify(loadconfig("config.txt", "toml")));

if (cli.bool("ssh")) {
  ssh.set("auth", "password");
  ssh.set("user", "testjssh");
  ssh.set("password", "123456");
  ssh.open("192.168.2.200");
  ssh.setenv("a", "123");
  log.info(JSON.stringify(ssh.exec("echo $a,$b", { b: "456" })));
  log.info(JSON.stringify(ssh.exec1("echo $a,$b", { b: "456" })));
  log.info(JSON.stringify(ssh.exec("pwd")));
  ssh.close();
}

if (cli.bool("prompt")) {
  log.info("prompt: %s", cli.prompt());
  log.info("prompt: %s", cli.prompt("what's your name: "));
}

if (cli.bool("download")) {
  log.info(
    "download: %s",
    http.download("https://gitee.com/leizongmin/jssh/raw/main/main.go")
  );
  log.info(
    "download: %s",
    http.download(
      "https://gitee.com/leizongmin/jssh/raw/main/main.go",
      path.join(__tmpdir, "test-download")
    )
  );
}

if (cli.bool("tcp")) {
  socket.timeout(2_000);
  log.info("socket: %v", socket.tcptest("baidu.com", 80));
  log.info("socket: %v", socket.tcptest("baidu.com", 81));
  const rawReq = [
    "GET /search?q=test HTTP/1.1",
    "Host: baidu.com",
    "User-Agent: jssh",
    "Accept: */*",
    "Connection: close",
    "",
    "",
  ].join("\r\n");
  println(rawReq);
  log.info("socket: %v", socket.tcpsend("baidu.com", 80, rawReq));
  log.info("socket: %v", socket.tcpsend("baidu.com", 81, rawReq));
}

sql.set("connMaxLifetime", 10_000);
sql.open("mysql", "root:@tcp(localhost:3306)/mysql?interpolateParams=true");
println(JSON.stringify(sql.query("show tables")));
println(JSON.stringify(sql.query("show databases")));
const tableName = `jssh_test_${Date.now()}`;
println(JSON.stringify(sql.exec(`create table ${tableName}(id int)`)));
println(
  JSON.stringify(
    sql.exec(`insert into ${tableName}(id) values (?),(?)`, 123, 456)
  )
);
println(JSON.stringify(sql.query(`select * from  ${tableName}`)));
println(JSON.stringify(sql.exec(`drop table ${tableName}`)));
sql.close();

println(evalfile("./hello.js"));
println(evalfile("./hello.js"));
println(evalfile("./__file_does_not_exist__.js", "println('hello')"));
fnFromHello();

exit(123);
