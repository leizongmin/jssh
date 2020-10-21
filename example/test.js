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
