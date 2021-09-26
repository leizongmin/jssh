// 定义全局对象jssh，将所有内置Rust函数写入jssh.op对象
const jssh = { op: {} };
try {
  const keyPrefix = "__builtin_op_";
  Object.keys(globalThis)
    .filter((k) => k.startsWith(keyPrefix))
    .forEach((k) => {
      const n = k.slice(keyPrefix.length);
      jssh.op[n] = globalThis[k];
      Object.defineProperty(jssh.op, n, { configurable: false });
      delete globalThis[k];
    });
  Object.defineProperty(jssh, "op", { configurable: false });
} catch (err) {
  console.error(err.message + "\n" + err.stack);
  throw err;
}
globalThis.jssh = jssh;
Object.defineProperty(globalThis, "jssh", { configurable: false });
const global = globalThis;

const __env = jssh.op.env();
const __args = jssh.op.args();
const __tmpdir = jssh.op.dir_temp();
const __homedir = jssh.op.dir_home();
const __downloaddir = jssh.op.dir_download();

function exit(code = 0) {
  return jssh.op.exit(code);
}

function print(...args) {
  return jssh.op.stdout_write(args.map((v) => v.toString()).join(" "));
}

function println(...args) {
  return jssh.op.stdout_write(args.map((v) => v.toString()).join(" ") + "\n");
}

function readline() {
  return jssh.op.stdin_read_line();
}

function stdoutlog(...args) {
  return jssh.op.stdout_write(args.map((v) => v.toString()).join(" ") + "\n");
}

function stderrlog(...args) {
  return jssh.op.stderr_write(args.map((v) => v.toString()).join(" ") + "\n");
}

function sleep(ms) {
  return jssh.op.sleep(ms);
}
