log.info("aaa: %s %s %s", __bin, __dirname, __filename)
log.error(JSON.stringify(__args))
println(JSON.stringify(__env))
println("%f", Date.now())

sleep(500)
setenv("__xx__", new Date().toString())
log.info(JSON.stringify(__env))

log.info(pwd())
log.info(cd(__dirname))
log.info(cwd())

exec("pwd")
log.info("%f %f %s", __code, __outputbytes, __output)

exec("pwd", {}, true)
log.info("%f %f %s", __code, __outputbytes, __output)

// exec("node")

if (exec(`ls -al ${__homedir}`, {}, true) === 0) {
    __output.split("\n").forEach(line => log.error(line))
}

fs.readdir(__homedir).forEach(f => log.error(JSON.stringify(f)))
log.info(fs.readfile(`${__homedir}/.gitconfig`))
log.info(JSON.stringify(fs.stat(`${__homedir}/.gitconfig`)))

set("xyz", 12345)
log.info("xyz = %f", xyz)

const file = `${__tmpdir}/${Date.now()}-${Math.random()}.txt`
log.info(fs.writefile(file, "hello"))
log.info(fs.appendfile(file, "world"))
log.info(fs.readfile(file))

log.info(path.abs("."))
log.info(path.base(file))
log.info(path.dir(file))
log.info(path.ext(file))
log.info(path.join("a", "b", "c"))
log.info(path.abs(path.join("a", "b", "c")))

log.info("%s, %s, %v, %v, %s, %s", cli.get(0), cli.get("n"), cli.bool("n"), cli.bool("x"), JSON.stringify(cli.args()), JSON.stringify(cli.opts()))

log.error(new Error().stack)

log.info(JSON.stringify(http.request("GET", "http://baidu.com")))
log.info(format("%s-%s", "aaa", "bbb"))

exit(123)