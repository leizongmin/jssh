// 增强的 ES 模块实现，改进 top-level await 支持

// 扩展现有的 ES 模块功能
if (typeof jssh !== 'undefined' && jssh.loadESModule) {
  
  // 保存原始的 loadESModule 函数
  const originalLoadESModule = jssh.loadESModule;
  const originalLoadESModuleSync = jssh.loadESModuleSync;
  
  // 改进的 ES 模块加载函数，更好地支持 top-level await
  const enhancedLoadESModule = async (specifier, referrer = __dirname) => {
    jssh.log.debug("Enhanced ES module load: specifier=%s, referrer=%s", specifier, referrer);
    
    // 如果是 .mjs 文件，使用增强的处理逻辑
    if (specifier.endsWith('.mjs') || specifier.includes('.mjs')) {
      const resolveESModulePath = jssh.__internal?.resolveESModulePath || ((name, dir) => {
        if (name.startsWith('./') || name.startsWith('../') || name.startsWith('/')) {
          return jssh.path.abs(jssh.path.join(dir, name));
        }
        return name;
      });
      
      const resolvedPath = resolveESModulePath(specifier, referrer);
      
      // 读取文件内容
      let content;
      try {
        content = jssh.fs.readfile(resolvedPath);
      } catch (error) {
        throw new Error(`Cannot load module ${specifier}: ${error.message}`);
      }
      
      // 检查是否包含 top-level await
      const hasTopLevelAwait = /(?:^|\n|\r\n?)\s*await\s+/.test(content);
      
      if (hasTopLevelAwait) {
        jssh.log.debug("Detected top-level await in module: %s", resolvedPath);
        
        // 创建一个支持 top-level await 的模块包装器
        const moduleWrapper = `
(async function() {
  // 模块导出对象
  const exports = {};
  const module = { exports };
  
  // 导出函数
  const exportDefault = (value) => { 
    exports.default = value; 
  };
  const exportNamed = (name, value) => { 
    exports[name] = value; 
  };
  
  // 处理各种导出语法
  const originalExport = globalThis.export;
  globalThis.export = { default: exportDefault, named: exportNamed };
  
  try {
    // 转换代码以支持导出
    let transformedContent = \`${content.replace(/`/g, '\\`').replace(/\${/g, '\\${')}\`;
    
    // 转换 export default
    transformedContent = transformedContent.replace(
      /export\\s+default\\s+([^;\\n]+)/g, 
      'exportDefault($1)'
    );
    
    // 转换 named exports
    transformedContent = transformedContent.replace(
      /export\\s+(?:const|let|var)\\s+([a-zA-Z_$][a-zA-Z0-9_$]*)\\s*=\\s*([^;\\n]+)/g,
      'const $1 = $2; exportNamed("$1", $1)'
    );
    
    // 转换 export function
    transformedContent = transformedContent.replace(
      /export\\s+function\\s+([a-zA-Z_$][a-zA-Z0-9_$]*)/g,
      'function $1'
    );
    transformedContent = transformedContent.replace(
      /function\\s+([a-zA-Z_$][a-zA-Z0-9_$]*)\\s*\\(/g,
      (match, name) => {
        if (!transformedContent.includes(\`exportNamed("\${name}"\`)) {
          return match.replace('function', 'const temp_func = function') + 
                 \`; exportNamed("\${name}", temp_func); function \${name}\`;
        }
        return match;
      }
    );
    
    // 转换 export { ... }
    transformedContent = transformedContent.replace(
      /export\\s*\\{([^}]+)\\}/g,
      (match, exports) => {
        return exports.split(',').map(exp => {
          const [from, to] = exp.trim().split(' as ');
          const exportName = to ? to.trim() : from.trim();
          const varName = from.trim();
          return \`exportNamed("\${exportName}", \${varName});\`;
        }).join('\\n');
      }
    );
    
    // 执行转换后的代码
    const AsyncFunction = Object.getPrototypeOf(async function(){}).constructor;
    const moduleFunction = new AsyncFunction('exportDefault', 'exportNamed', 'exports', 'module', transformedContent);
    
    await moduleFunction(exportDefault, exportNamed, exports, module);
    
    return exports;
    
  } finally {
    // 恢复原始的 export
    globalThis.export = originalExport;
  }
})()`;

        try {
          const result = await eval(moduleWrapper);
          return result;
        } catch (error) {
          jssh.log.error("Failed to execute enhanced ES module: %s", error.message);
          throw error;
        }
      }
    }
    
    // 对于其他情况，使用原始的加载函数
    return await originalLoadESModule(specifier, referrer);
  };
  
  // 替换原始函数
  jssh.loadESModule = enhancedLoadESModule;
  jssh.import = enhancedLoadESModule;
  
  // 同样为动态 import 函数添加支持
  if (typeof globalThis.import !== 'undefined') {
    globalThis.import = enhancedLoadESModule;
  }
  
  console.log('✅ Enhanced ES Module with top-level await support loaded');
}