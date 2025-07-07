// 增强的 ES 模块实现，改进 top-level await 支持

// 扩展现有的 ES 模块功能
console.log('Checking jssh availability:', typeof jssh !== 'undefined');
console.log('Checking jssh.loadESModule:', typeof jssh?.loadESModule);

if (typeof jssh !== 'undefined' && jssh.loadESModule) {
  console.log('✅ Conditions met, proceeding with enhancement...');
  
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
      
      // 对所有 .mjs 文件使用增强逻辑，不仅仅是包含 top-level await 的
      jssh.log.debug("Processing ES module: %s, hasTopLevelAwait: %s", resolvedPath, hasTopLevelAwait);
      
      // 创建一个支持增强功能的模块包装器
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
    
    // 转换 export function - 修复语法错误
    // 收集导出的函数名
    const exportedFunctions = [];
    let match;
    const exportFuncRegex = /export\\s+function\\s+([a-zA-Z_$][a-zA-Z0-9_$]*)/g;
    while ((match = exportFuncRegex.exec(transformedContent)) !== null) {
      exportedFunctions.push(match[1]);
    }
    
    // 移除 export 关键字
    transformedContent = transformedContent.replace(
      /export\\s+function\\s+/g,
      'function '
    );
    
    // 在转换后的内容末尾添加导出调用
    if (exportedFunctions.length > 0) {
      const exportCalls = exportedFunctions.map(name => 
        'exportNamed("' + name + '", ' + name + ');'
      ).join('\\n');
      transformedContent += '\\n' + exportCalls;
    }
    
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
    
    // 对于其他情况，使用原始的加载函数
    return await originalLoadESModule(specifier, referrer);
  };
  
  // 增强的同步模块加载函数
  const enhancedLoadESModuleSync = (specifier, referrer = __dirname) => {
    jssh.log.debug("Enhanced ES module load sync: specifier=%s, referrer=%s", specifier, referrer);
    
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
      
      // 对于同步版本，创建一个简化的模块包装器（不支持 top-level await）
      const moduleWrapper = `
(function() {
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
    
    // 移除 await 语句，因为同步版本不支持
    transformedContent = transformedContent.replace(/await\\s+/g, '');
    
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
    const exportedFunctions = [];
    let match;
    const exportFuncRegex = /export\\s+function\\s+([a-zA-Z_$][a-zA-Z0-9_$]*)/g;
    while ((match = exportFuncRegex.exec(transformedContent)) !== null) {
      exportedFunctions.push(match[1]);
    }
    
    transformedContent = transformedContent.replace(
      /export\\s+function\\s+/g,
      'function '
    );
    
    if (exportedFunctions.length > 0) {
      const exportCalls = exportedFunctions.map(name => 
        'exportNamed("' + name + '", ' + name + ');'
      ).join('\\n');
      transformedContent += '\\n' + exportCalls;
    }
    
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
    const moduleFunction = new Function('exportDefault', 'exportNamed', 'exports', 'module', transformedContent);
    moduleFunction(exportDefault, exportNamed, exports, module);
    
    return exports;
    
  } finally {
    // 恢复原始的 export
    globalThis.export = originalExport;
  }
})()`;

      try {
        const result = eval(moduleWrapper);
        return result;
      } catch (error) {
        jssh.log.error("Failed to execute enhanced ES module sync: %s", error.message);
        throw error;
      }
    }
    
    // 对于其他情况，使用原始的加载函数
    return originalLoadESModuleSync(specifier, referrer);
  };

  // 替换原始函数
  console.log('🔄 Replacing functions...');
  console.log('Before replacement - loadESModuleSync type:', typeof jssh.loadESModuleSync);
  
  jssh.loadESModule = enhancedLoadESModule;
  jssh.loadESModuleSync = enhancedLoadESModuleSync;
  jssh.import = enhancedLoadESModule;
  jssh.importSync = enhancedLoadESModuleSync;
  
  console.log('After replacement - loadESModuleSync type:', typeof jssh.loadESModuleSync);
  console.log('Functions replaced successfully');
  
  // 同样为动态 import 函数添加支持
  if (typeof globalThis.import !== 'undefined') {
    globalThis.import = enhancedLoadESModule;
  }
  
  console.log('✅ Enhanced ES Module with top-level await support loaded');
}