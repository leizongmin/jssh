
{
  const cli = {};
  const _args = (cli._args = []);
  const _opts = (cli._opts = {});

  const getFlagName = (s) => {
    if (s.startsWith("--")) {
      return s.slice(2);
    }
    if (s.startsWith("-")) {
      return s.slice(1);
    }
  };

  for (let i = 2; i < jssh.__args.length; i++) {
    const v = jssh.__args[i];
    const v2 = jssh.__args[i + 1];
    if (v.startsWith("-")) {
      const r = v.match(/^--?([\w\-_]+)=(.*)$/);
      if (r) {
        _opts[r[1]] = r[2];
      } else {
        if (v2 !== undefined) {
          if (v2.startsWith("-")) {
            _opts[getFlagName(v)] = true;
          } else {
            _opts[getFlagName(v)] = v2;
            i++;
          }
        } else {
          _opts[getFlagName(v)] = true;
        }
      }
    } else {
      _args.push(v);
    }
  }

  cli.get = function get(n) {
    if (typeof n === "number") {
      return _args[n];
    } else {
      return _opts[n];
    }
  };

  cli.bool = function bool(n) {
    if (_opts[n] === false || _opts[n] === undefined) return false;
    if (_opts[n] === true) return true;
    const s = _opts[n].toLowerCase();
    return !(s === "0" || s === "f" || s === "false");
  };

  cli.args = function args() {
    return [..._args];
  };

  cli.opts = function opts() {
    return { ..._opts };
  };

  cli.prompt = function prompt(message) {
    if (message) jssh.print(message);
    return jssh.readline();
  };

  cli._subcommand = {};

  cli.subcommand = function subcommand(name, callback) {
    if (typeof callback !== `function`) {
      throw new TypeError(`callback expected a function`);
    }
    if (cli._subcommand[name]) {
      throw new Error(`subcommand ${name} is already registered`);
    }
    cli._subcommand[name] = callback;
  };

  cli.subcommandstart = function subcommandstart() {
    const name = cli.get(0);
    if (cli._subcommand[name]) {
      return cli._subcommand[name]();
    }
    if (cli._subcommand[`*`]) {
      return cli._subcommand[`*`]();
    }
    throw new Error(`unrecognized subcommand ${name}`);
  };

  jssh.cli = cli;
  Object.freeze(cli);
}
