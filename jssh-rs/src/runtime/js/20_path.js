const path = {};
{
  path.join = function (...args) {
    return jssh.op.path_join(
      ...args.map((s) => (s.startsWith("./") ? s.slice(2) : s))
    );
  };

  path.abs = function (p) {
    return jssh.op.path_abs(p);
  };

  path.base = function (p) {
    return jssh.op.path_base(p);
  };

  path.ext = function (p) {
    return jssh.op.path_ext(p);
  };

  path.dir = function (p) {
    return jssh.op.path_dir(p);
  };
}
Object.freeze(path);
