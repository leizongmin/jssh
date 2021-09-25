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
log.info("hello, %s", __env, __env);

exit(1);
