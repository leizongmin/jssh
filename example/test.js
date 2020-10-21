log("aaa: %s %s %s", __bin, __dirname, __filename)
log(JSON.stringify(__args))
println(JSON.stringify(__env))
println("%f", Date.now())

sleep(2000)
setenv("__xx__", new Date().toString())
log(JSON.stringify(__env))

exec("node")
