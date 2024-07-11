{
  const removeShebangLine = (data) => {
    if (typeof data !== "string") throw new Error(`unexpected input: ${data}`);
    if (!data.startsWith("#!")) return data;
    return data.replace(/^#![^\n]*/, "");
  };

  const isHttpUrl = (s) => /^https?:\/\//gi.test(s);

  const readUrlContent = (url) => {
    jssh.log.debug("require: read content from %s", url);
    return jssh.http.get(url);
  };

  const resolveWithExtension = (name) => {
    const extension = [".json", ".js"];
    if (jssh.fs.exist(name)) {
      if (jssh.fs.stat(name).isdir) {
        // 如果是目录，尝试 ${name}/package.json
        const pkgFile = jssh.path.join(name, "package.json");
        if (jssh.fs.exist(pkgFile)) {
          const pkg = loadModuleFromJsonContent(pkgFile, jssh.fs.readfile(pkgFile));
          if (pkg.main) {
            return resolveWithExtension(jssh.path.join(name, pkg.main));
          }
        }
        // 再尝试 ${name}/index.js, ${name}/index.json
        const indexFile = jssh.path.join(name, "index");
        if (jssh.fs.exist(indexFile)) {
          return indexFile;
        }
        for (const ext of extension) {
          if (jssh.fs.exist(indexFile + ext)) {
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
        if (jssh.fs.exist(name + ext)) {
          return name + ext;
        }
      }
      // 再尝试 ${name}/index.js, ${name}/index.json
      const indexFile = jssh.path.join(name, "index");
      if (jssh.fs.exist(indexFile)) {
        return indexFile;
      }
      for (const ext of extension) {
        if (jssh.fs.exist(indexFile + ext)) {
          return indexFile + ext;
        }
      }
    }
  };

  const resolveNpmModulePath = (name, dir) => {
    const paths = [];
    let d = dir;
    while (true) {
      let p = jssh.path.abs(jssh.path.join(d, "node_modules"));
      paths.push(p);
      const d2 = jssh.path.dir(d);
      if (d2 === d) {
        break;
      } else {
        d = d2;
      }
    }
    for (const p of paths) {
      const ret = resolveWithExtension(jssh.path.join(p, name));
      if (ret) return jssh.path.abs(ret);
    }
  };

  const httpPkgSites = [
    { prefix: 'unpkg:', getPath: (name) => `https://unpkg.com/${name}` },
    // { prefix: 'jsdelivr:', getPath: (name) => `https://cdn.jsdelivr.net/npm/${name}` },
    // { prefix: 'cdn:', getPath: (name) => `https://cdn.jsdelivr.net/npm/${name}` },
  ];

  const resolveModulePath = (name, dir) => {
    if (name === "." || name.startsWith("/") || name.startsWith("./")) {
      if (isHttpUrl(dir)) {
        return jssh.path.abs(jssh.path.join(dir, name));
      }
      const ret = resolveWithExtension(jssh.path.join(dir, name));
      if (ret) {
        return jssh.path.abs(ret);
      }
      return jssh.path.abs(name);
    }

    if (isHttpUrl(name)) {
      return name;
    }

    if (name.startsWith('npm:')) {
      return resolveNpmModulePath(name.slice('npm:'.length), dir);
    }

    for (const site of httpPkgSites) {
      if (name.startsWith(site.prefix)) {
        const pkgName = name.slice(site.prefix.length);
        return site.getPath(pkgName);
      }
    }

    return resolveNpmModulePath(name, dir);
  };

  const requiremodule = (name, dir = __dirname) => {
    jssh.log.debug("require: name=%s, dir=%s", name, dir);

    if (typeof name !== "string") {
      throw new TypeError(`module name expected string type`);
    }
    if (!name) {
      throw new TypeError(`empty module name`);
    }
    if (!dir) {
      throw new TypeError(`empty module dir`);
    }

    let file = resolveModulePath(name, dir);
    if (!file) {
      throw new Error(`cannot resolve module "${name}" on path "${dir}"`);
    }

    if (require.cache[file]) {
      return require.cache[file];
    }

    try {
      let content = "";
      if (isHttpUrl(file)) {
        const res = readUrlContent(file);
        if (res.status === 200) {
          content = res.body;
          file = res.url;
        } else {
          // FIXME: 尝试加上 .js 后缀，以后优化此方法
          const res2 = readUrlContent(file + ".js");
          if (res2.status === 200) {
            content = res2.body;
            file = res.url;
          } else {
            throw new Error(`http get "${file}" status "${res.status}"`);
          }
        }
      } else {
        content = jssh.fs.readfile(file);
      }
      if (file.endsWith(".json")) {
        return loadModuleFromJsonContent(file);
      } else {
        return loadModuleFromJsContent(file, jssh.path.dir(file), content);
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
    return (require.cache[__filename] = jssh.evalfile(__filename, wrapped).exports);
  };

  const require = (name) => {
    return requiremodule(name, __dirname);
  };
  require.cache = {};
  jssh.require = require;
  jssh.requiremodule = requiremodule;

  const importModuleCallbacks = [];

  const registerImportModuleCallback = (callback) => {
    if (typeof callback !== "function") {
      throw new TypeError("callback must be a function");
    }
    importModuleCallbacks.push(callback);
  };

  const callImportModuleCallbacks = (name, dir) => {
    for (const callback of importModuleCallbacks) {
      const result = callback(name, dir);
      if (result) {
        return result;
      }
    }
  };

  jssh.registerImportModuleCallback = registerImportModuleCallback;
  jssh.callImportModuleCallbacks = callImportModuleCallbacks;
}
