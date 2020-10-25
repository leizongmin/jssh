const cli = {};

{
  const args = [];
  const opts = {};

  function getFlagName(s) {
    if (s.startsWith("--")) {
      return s.slice(2);
    }
    if (s.startsWith("-")) {
      return s.slice(1);
    }
  }

  for (let i = 2; i < __args.length; i++) {
    const v = __args[i];
    const v2 = __args[i + 1];
    if (v.startsWith("-")) {
      const r = v.match(/^--?([\w\-_]+)=(.*)$/);
      if (r) {
        opts[r[1]] = r[2];
      } else {
        if (v2 !== undefined) {
          if (v2.startsWith("-")) {
            opts[getFlagName(v)] = true;
          } else {
            opts[getFlagName(v)] = v2;
            i++;
          }
        } else {
          opts[getFlagName(v)] = true;
        }
      }
    } else {
      args.push(v);
    }
  }

  cli.get = function (n) {
    if (typeof n === "number") {
      return args[n];
    } else {
      return opts[n];
    }
  };

  cli.bool = function (n) {
    if (opts[n] === false || opts[n] === undefined) return false;
    if (opts[n] === true) return true;
    const s = opts[n].toLowerCase();
    return !(s === "0" || s === "f" || s === "false");
  };

  cli.args = function () {
    return [...args];
  };

  cli.opts = function () {
    return { ...opts };
  };

  cli.prompt = function (message) {
    if (message) print(message);
    return readline();
  };
}
