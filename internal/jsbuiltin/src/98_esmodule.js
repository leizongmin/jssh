{
  // ES Module support for jssh
  // This module provides import() function and ES module loading capabilities

  const moduleCache = new Map(); // Cache for ES modules
  const pendingModules = new Map(); // Track pending module loads to avoid circular imports

  // Check if a file is an ES module by detecting import/export statements
  const isESModule = (content) => {
    // Remove comments and strings to avoid false positives
    const cleanContent = content
      .replace(/\/\*[\s\S]*?\*\//g, '')  // Remove block comments
      .replace(/\/\/.*$/gm, '')          // Remove line comments
      .replace(/"[^"\\]*(?:\\.[^"\\]*)*"/g, '""')  // Remove double quoted strings
      .replace(/'[^'\\]*(?:\\.[^'\\]*)*'/g, "''")  // Remove single quoted strings
      .replace(/`[^`\\]*(?:\\.[^`\\]*)*`/g, '``'); // Remove template literals

    // Check for ES module syntax
    const hasImport = /\bimport\s+/.test(cleanContent) || /\bimport\s*\(/.test(cleanContent);
    const hasExport = /\bexport\s+/.test(cleanContent);
    
    return hasImport || hasExport;
  };

  // Resolve module path for ES modules (similar to CommonJS but with .mjs support)
  const resolveESModulePath = (name, dir) => {
    const extension = [".mjs", ".js", ".json"];
    
    if (name === "." || name.startsWith("/") || name.startsWith("./") || name.startsWith("../")) {
      if (isHttpUrl(dir)) {
        return jssh.path.abs(jssh.path.join(dir, name));
      }
      
      // Try exact file first
      if (jssh.fs.exist(jssh.path.join(dir, name))) {
        const fullPath = jssh.path.join(dir, name);
        if (!jssh.fs.stat(fullPath).isdir) {
          return jssh.path.abs(fullPath);
        }
      }
      
      // Try with extensions
      for (const ext of extension) {
        const pathWithExt = jssh.path.join(dir, name + ext);
        if (jssh.fs.exist(pathWithExt)) {
          return jssh.path.abs(pathWithExt);
        }
      }
      
      // Try as directory with package.json
      const pkgPath = jssh.path.join(dir, name, "package.json");
      if (jssh.fs.exist(pkgPath)) {
        try {
          const pkg = JSON.parse(jssh.fs.readfile(pkgPath));
          if (pkg.module) {
            return resolveESModulePath(pkg.module, jssh.path.join(dir, name));
          }
          if (pkg.main) {
            return resolveESModulePath(pkg.main, jssh.path.join(dir, name));
          }
        } catch (e) {
          // Ignore JSON parse errors
        }
      }
      
      // Try index files
      for (const ext of extension) {
        const indexPath = jssh.path.join(dir, name, "index" + ext);
        if (jssh.fs.exist(indexPath)) {
          return jssh.path.abs(indexPath);
        }
      }
      
      return jssh.path.abs(jssh.path.join(dir, name));
    }

    if (isHttpUrl(name)) {
      return name;
    }

    // For npm modules, delegate to existing resolve logic
    return resolveModulePath(name, dir);
  };

  // Load ES module content and handle different sources (synchronous version)
  const loadESModuleContentSync = (filepath) => {
    if (isHttpUrl(filepath)) {
      const res = readUrlContent(filepath);
      if (res.status === 200) {
        return { content: res.body, resolvedPath: res.url };
      } else {
        // Try with .mjs extension
        const res2 = readUrlContent(filepath + ".mjs");
        if (res2.status === 200) {
          return { content: res2.body, resolvedPath: res2.url };
        }
        // Try with .js extension
        const res3 = readUrlContent(filepath + ".js");
        if (res3.status === 200) {
          return { content: res3.body, resolvedPath: res3.url };
        }
        throw new Error(`Failed to load module from ${filepath}: ${res.status}`);
      }
    } else {
      if (!jssh.fs.exist(filepath)) {
        throw new Error(`Module file not found: ${filepath}`);
      }
      return { 
        content: jssh.fs.readfile(filepath), 
        resolvedPath: jssh.path.abs(filepath) 
      };
    }
  };

  // Load and execute ES module (synchronous version)
  const loadESModuleSync = (specifier, referrer = __dirname) => {
    jssh.log.debug("ES module load sync: specifier=%s, referrer=%s", specifier, referrer);
    
    const resolvedPath = resolveESModulePath(specifier, referrer);
    if (!resolvedPath) {
      throw new Error(`Cannot resolve module "${specifier}" from "${referrer}"`);
    }

    // Check cache first
    if (moduleCache.has(resolvedPath)) {
      return moduleCache.get(resolvedPath);
    }

    // Check for circular imports
    if (pendingModules.has(resolvedPath)) {
      throw new Error(`Circular dependency detected: ${resolvedPath}`);
    }

    try {
      pendingModules.set(resolvedPath, true);
      
      const { content, resolvedPath: finalPath } = loadESModuleContentSync(resolvedPath);
      
      // Handle JSON modules
      if (finalPath.endsWith('.json')) {
        const moduleExports = JSON.parse(content);
        moduleCache.set(resolvedPath, { default: moduleExports, ...moduleExports });
        return moduleCache.get(resolvedPath);
      }

      // Determine if this is an ES module or CommonJS
      const isModule = isESModule(content);
      
      if (isModule) {
        // Handle ES module - create a simulated module execution environment
        const moduleExports = {};
        let defaultExport = undefined;
        
        // Create export functions
        const exportDefault = (value) => { 
          defaultExport = value; 
          moduleExports.default = value; 
        };
        const exportNamed = (name, value) => { 
          moduleExports[name] = value; 
        };
        
        // Transform the ES module code to use our export functions
        let transformedContent = content;
        
        // Transform export default statements
        transformedContent = transformedContent.replace(
          /export\s+default\s+(.+?)(?:;|$)/gm, 
          'exportDefault($1);'
        );
        
        // Transform named export declarations
        transformedContent = transformedContent.replace(
          /export\s+(?:const|let|var)\s+([a-zA-Z_$][a-zA-Z0-9_$]*)\s*=\s*(.+?)(?:;|$)/gm,
          'const $1 = $2; exportNamed("$1", $1);'
        );
        
        // Transform function exports
        transformedContent = transformedContent.replace(
          /export\s+function\s+([a-zA-Z_$][a-zA-Z0-9_$]*)/g,
          'function $1'
        );
        transformedContent = transformedContent.replace(
          /function\s+([a-zA-Z_$][a-zA-Z0-9_$]*)/g,
          (match, name) => match + `; exportNamed("${name}", ${name});`
        );
        
        // Transform export lists
        transformedContent = transformedContent.replace(
          /export\s*\{([^}]+)\}/g,
          (match, exports) => {
            return exports.split(',').map(exp => {
              const [from, to] = exp.trim().split(' as ');
              const exportName = to ? to.trim() : from.trim();
              const varName = from.trim();
              return `exportNamed("${exportName}", ${varName});`;
            }).join('\n');
          }
        );
        
        // Remove any remaining import statements (basic transformation)
        transformedContent = transformedContent.replace(
          /import\s+.*?from\s+['"][^'"]+['"];?/g,
          '// import removed for sync execution'
        );
        transformedContent = transformedContent.replace(
          /import\s+['"][^'"]+['"];?/g,
          '// import removed for sync execution'
        );
        
        // Wrap the code in a function and execute it
        const wrappedCode = `
(function() {
  const moduleExports = {};
  const exportDefault = function(value) { 
    moduleExports.default = value; 
  };
  const exportNamed = function(name, value) { 
    moduleExports[name] = value; 
  };
  
  ${transformedContent}
  
  return moduleExports;
})();
        `;
        
        try {
          const result = jssh.evalfile(finalPath, wrappedCode);
          moduleCache.set(resolvedPath, result);
          return result;
        } catch (evalError) {
          throw new Error(`Failed to execute ES module "${finalPath}": ${evalError.message}`);
        }
      } else {
        // Handle as CommonJS module using existing logic
        const moduleExports = loadModuleFromJsContent(finalPath, jssh.path.dir(finalPath), content);
        // Wrap CommonJS exports for ES module compatibility
        const wrappedExports = { default: moduleExports, ...moduleExports };
        moduleCache.set(resolvedPath, wrappedExports);
        return wrappedExports;
      }
    } finally {
      pendingModules.delete(resolvedPath);
    }
  };

  // Transform import statements to use the import() function
  const transformESModuleImports = (content) => {
    // This is a basic transformation - in a real implementation you might want to use a proper parser
    // Transform: import { x, y } from 'module' -> const { x, y } = await import('module')
    content = content.replace(
      /import\s+\{([^}]+)\}\s+from\s+['"`]([^'"`]+)['"`]/g,
      "const { $1 } = await import('$2')"
    );
    
    // Transform: import x from 'module' -> const x = (await import('module')).default
    content = content.replace(
      /import\s+([a-zA-Z_$][a-zA-Z0-9_$]*)\s+from\s+['"`]([^'"`]+)['"`]/g,
      "const $1 = (await import('$2')).default"
    );
    
    // Transform: import * as x from 'module' -> const x = await import('module')
    content = content.replace(
      /import\s+\*\s+as\s+([a-zA-Z_$][a-zA-Z0-9_$]*)\s+from\s+['"`]([^'"`]+)['"`]/g,
      "const $1 = await import('$2')"
    );
    
    // Transform: import 'module' -> await import('module')
    content = content.replace(
      /import\s+['"`]([^'"`]+)['"`]/g,
      "await import('$1')"
    );

    return content;
  };

  // Create module exports object
  const createModuleExports = () => {
    const exports = {};
    const module = { exports };
    
    // Support both CommonJS and ES module patterns
    return { exports, module };
  };

  // Load and execute ES module
  const loadESModule = async (specifier, referrer = __dirname) => {
    jssh.log.debug("ES module load: specifier=%s, referrer=%s", specifier, referrer);
    
    const resolvedPath = resolveESModulePath(specifier, referrer);
    if (!resolvedPath) {
      throw new Error(`Cannot resolve module "${specifier}" from "${referrer}"`);
    }

    // Check cache first
    if (moduleCache.has(resolvedPath)) {
      return moduleCache.get(resolvedPath);
    }

    // Check for circular imports
    if (pendingModules.has(resolvedPath)) {
      throw new Error(`Circular dependency detected: ${resolvedPath}`);
    }

    try {
      pendingModules.set(resolvedPath, true);
      
      const { content, resolvedPath: finalPath } = loadESModuleContentSync(resolvedPath);
      
      // Handle JSON modules
      if (finalPath.endsWith('.json')) {
        const moduleExports = JSON.parse(content);
        moduleCache.set(resolvedPath, { default: moduleExports, ...moduleExports });
        return moduleCache.get(resolvedPath);
      }

      // Determine if this is an ES module or CommonJS
      const isModule = isESModule(content);
      
      if (isModule) {
        // Transform and execute as ES module
        let transformedContent = transformESModuleImports(content);
        
        // Wrap in async function to handle top-level await
        const moduleWrapper = `
(async function() {
  const exports = {};
  let defaultExport = undefined;
  
  // Override export functionality
  const exportDefault = (value) => { defaultExport = value; exports.default = value; };
  const exportNamed = (name, value) => { exports[name] = value; };
  
  // Global export function
  globalThis.export = { default: exportDefault, named: exportNamed };
  
  // Transform export statements
  ${transformedContent.replace(/export\s+default\s+/g, 'exportDefault(').replace(/export\s*\{([^}]+)\}/g, (match, exports) => {
    return exports.split(',').map(exp => {
      const [from, to] = exp.trim().split(' as ');
      return `exportNamed('${to || from.trim()}', ${from.trim()});`;
    }).join('\n');
  }).replace(/export\s+(?:const|let|var|function|class)\s+([a-zA-Z_$][a-zA-Z0-9_$]*)/g, (match, name) => {
    return match.replace('export ', '') + `; exportNamed('${name}', ${name});`;
  })}
  
  return exports;
})()`;

        const moduleExports = await jssh.evalfile(finalPath, moduleWrapper);
        moduleCache.set(resolvedPath, moduleExports);
        return moduleExports;
      } else {
        // Handle as CommonJS module using existing logic
        const moduleExports = loadModuleFromJsContent(finalPath, jssh.path.dir(finalPath), content);
        // Wrap CommonJS exports for ES module compatibility
        const wrappedExports = { default: moduleExports, ...moduleExports };
        moduleCache.set(resolvedPath, wrappedExports);
        return wrappedExports;
      }
    } finally {
      pendingModules.delete(resolvedPath);
    }
  };

  // Implement dynamic import() function
  const dynamicImport = async (specifier) => {
    // Get current module directory from call stack or use global __dirname
    const referrer = __dirname; // In a real implementation, you'd extract this from the call stack
    return await loadESModule(specifier, referrer);
  };

  // Implement synchronous import function
  const syncImport = (specifier, referrer = __dirname) => {
    return loadESModuleSync(specifier, referrer);
  };

  // Make import() available globally
  globalThis.import = dynamicImport;
  
  // Also add to jssh namespace for consistency
  jssh.import = dynamicImport;
  jssh.importSync = syncImport;
  jssh.loadESModule = loadESModule;
  jssh.loadESModuleSync = loadESModuleSync;
  jssh.isESModule = isESModule;

  // Helper functions from CommonJS module system
  const { removeShebangLine, isHttpUrl, readUrlContent, resolveModulePath, loadModuleFromJsContent } = (() => {
    // Import necessary functions from the CommonJS module system
    // These should be available from 00_module.js
    return {
      removeShebangLine: globalThis.removeShebangLine || ((data) => {
        if (typeof data !== "string") throw new Error(`unexpected input: ${data}`);
        if (!data.startsWith("#!")) return data;
        return data.replace(/^#![^\n]*/, "");
      }),
      isHttpUrl: globalThis.isHttpUrl || ((s) => /^https?:\/\//gi.test(s)),
      readUrlContent: globalThis.readUrlContent || ((url) => {
        jssh.log.debug("ES module: read content from %s", url);
        return jssh.http.get(url);
      }),
      resolveModulePath: globalThis.resolveModulePath || (typeof require !== 'undefined' && require.resolve) || (() => {
        throw new Error("resolveModulePath not available");
      }),
      loadModuleFromJsContent: globalThis.loadModuleFromJsContent || (() => {
        throw new Error("loadModuleFromJsContent not available");
      })
    };
  })();

  jssh.log.debug("ES Module support initialized");
}