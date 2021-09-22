// 定义全局对象jssh，将所有内置Rust函数写入jssh.op对象
const jssh = {op: {}};
try {
  const keyPrefix = "__builtin_op_";
  Object.keys(globalThis).filter(k => k.startsWith(keyPrefix)).forEach(k => {
    const n = k.slice(keyPrefix.length);
    jssh.op[n] = globalThis[k];
    Object.defineProperty(jssh.op, n, {configurable: false})
    delete globalThis[k];
  });
  Object.defineProperty(jssh, "op", {configurable: false});
} catch (err) {
  console.error(err.message + "\n" + err.stack);
  throw err;
}
globalThis.jssh = jssh;
Object.defineProperty(globalThis, "jssh", {configurable: false});

function println(...args) {
  jssh.op.stdout_write(args.map(v => v.toString()).join(" ") + "\n");
}
