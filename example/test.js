log("aaa: %s %s %s", __bin, __dirname, __filename)
log(JSON.stringify(__args))
println(JSON.stringify(__env))
println("%f", Date.now())

sleep(500)
setenv("__xx__", new Date().toString())
log(JSON.stringify(__env))

log(pwd())
log(cd(__dirname))
log(cwd())

exec("pwd")
log("%f %f %s", __code, __outputBytes, __output)

exec("pwd", {}, true)
log("%f %f %s", __code, __outputBytes, __output)

// exec("node")

if (exec(`ls -al ${__homedir}`, {}, true) === 0) {
    __output.split("\n").forEach(line => log(line))
}

fs.readdir(__homedir).forEach(f => log(JSON.stringify(f)))
log(fs.readfile(`${__homedir}/.gitconfig`))
log(JSON.stringify(fs.readstat(`${__homedir}/.gitconfig`)))

set("xyz", 12345)
log("xyz = %f", xyz)

const file = `${__tmpdir}/${Date.now()}-${Math.random()}.txt`
log(fs.writefile(file, "hello"))
log(fs.appendfile(file, "world"))
log(fs.readfile(file))

log(path.abs("."))
log(path.base(file))
log(path.dir(file))
log(path.ext(file))
log(path.join("a", "b", "c"))
log(path.abs(path.join("a", "b", "c")))

log("%s, %s, %v, %v, %s, %s", cli.get(0), cli.get("n"), cli.bool("n"), cli.bool("x"), JSON.stringify(cli.args()), JSON.stringify(cli.opts()))

exit(123)