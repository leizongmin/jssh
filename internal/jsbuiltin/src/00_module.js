const global = globalThis || this;

{
  const removeShebangLine = (data) => {
    if (typeof data !== "string") throw new Error(`unexpected input: ${data}`);
    if (!data.startsWith("#!")) return data;
    return data.replace(/^#![^\n]*/, "");
  };

  const isHttpUrl = (s) => /^https?:\/\//gi.test(s);

  const resolveWithExtension = (name) => {
    const extension = [".json", ".js"];
    if (fs.exist(name)) {
      if (fs.stat(name).isdir) {
        // 如果是目录，尝试 ${name}/package.json
        const pkgFile = path.join(name, "package.json");
        if (fs.exist(pkgFile)) {
          const pkg = loadModuleFromJsonContent(pkgFile, fs.readfile(pkgFile));
          if (pkg.main) {
            return resolveWithExtension(path.join(name, pkg.main));
          }
        }
        // 再尝试 ${name}/index.js, ${name}/index.json
        const indexFile = path.join(name, "index");
        if (fs.exist(indexFile)) {
          return indexFile;
        }
        for (const ext of extension) {
          if (fs.exist(indexFile + ext)) {
            return indexFile + ext;
          }
        }
      } else {
        // 文件则直接返回
        return name;
      }
    } else {
      // 如果文件不存在，尝试 ${name}.js, ${name}.json
      for (const ext of extension) {
        if (fs.exist(name + ext)) {
          return name + ext;
        }
      }
      // 再尝试 ${name}/index.js, ${name}/index.json
      const indexFile = path.join(name, "index");
      if (fs.exist(indexFile)) {
        return indexFile;
      }
      for (const ext of extension) {
        if (fs.exist(indexFile + ext)) {
          return indexFile + ext;
        }
      }
    }
  };

  const resolveModulePath = (name, dir) => {
    if (name === "." || name.startsWith("/") || name.startsWith("./")) {
      if (isHttpUrl(dir)) {
        return path.abs(path.join(dir, name));
      }
      const ret = resolveWithExtension(path.join(dir, name));
      if (ret) {
        return path.abs(ret);
      }
      return path.abs(name);
    }

    if (isHttpUrl(name)) {
      return name;
    }

    const paths = [];
    let d = dir;
    while (true) {
      let p = path.abs(path.join(d, "node_modules"));
      paths.push(p);
      const d2 = path.dir(d);
      if (d2 === d) {
        break;
      } else {
        d = d2;
      }
    }
    for (const p of paths) {
      const ret = resolveWithExtension(path.join(p, name));
      if (ret) return path.abs(ret);
    }
  };

  const requiremodule = (name, dir = __dirname) => {
    if (typeof name !== "string") {
      throw new TypeError(`module name expected string type`);
    }
    if (!name) {
      throw new TypeError(`empty module name`);
    }
    if (!dir) {
      throw new TypeError(`empty module dir`);
    }

    const file = resolveModulePath(name, dir);
    if (!file) {
      throw new Error(`cannot resolve module "${name}" on path "${dir}"`);
    }

    if (require.cache[file]) {
      return require.cache[file];
    }

    try {
      let content = "";
      if (isHttpUrl(file)) {
        const res = http.get(file);
        if (res.status === 200) {
          content = res.body;
        } else {
          // FIXME: 尝试加上 .js 后缀，以后优化此方法
          const res2 = http.get(file + ".js");
          if (res2.status === 200) {
            content = res2.body;
          } else {
            throw new Error(`http get "${file}" status "${res.status}"`);
          }
        }
      } else {
        content = fs.readfile(file);
      }
      if (file.endsWith(".json")) {
        return loadModuleFromJsonContent(file);
      } else {
        return loadModuleFromJsContent(file, path.dir(file), content);
      }
    } catch (err) {
      const err2 = new Error(
        `cannot load module "${name}": ${err.message}\n${err.stack}`
      );
      err2.moduleName = name;
      err2.resolvedFilename = file;
      err2.originError = err;
      throw err2;
    }
  };

  const loadModuleFromJsonContent = (filename, content) => {
    return (require.cache[filename] = JSON.parse(content));
  };

  const loadModuleFromJsContent = (filename, dirname, content) => {
    content = removeShebangLine(content);
    const wrapped = `
(function (require, module, __dirname, __filename) { var exports = module.exports; ${content}
return module;
})(function require(name) {
  return requiremodule(name, "${dirname}");
}, {exports:{},parent:this}, "${dirname}", "${filename}")
`.trimStart();
    return (require.cache[__filename] = evalfile(__filename, wrapped).exports);
  };

  const require = (name) => {
    return requiremodule(name, __dirname);
  };
  require.cache = {};
  global.require = require;
  global.requiremodule = requiremodule;
}
