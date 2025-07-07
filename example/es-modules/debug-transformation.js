#!/usr/bin/env jssh

println('=== 调试转换过程 ===');

// 读取文件内容
const content = jssh.fs.readfile('./test-export-function.mjs');
println('1. 原始内容:');
println(content);

// 模拟转换过程
let transformedContent = content;

println('2. 开始转换...');

// 移除 await 语句
transformedContent = transformedContent.replace(/await\s+/g, '');

// 转换 export default
transformedContent = transformedContent.replace(
  /export\s+default\s+([^;\n]+)/g, 
  'exportDefault($1)'
);

// 转换 named exports
transformedContent = transformedContent.replace(
  /export\s+(?:const|let|var)\s+([a-zA-Z_$][a-zA-Z0-9_$]*)\s*=\s*([^;\n]+)/g,
  'const $1 = $2; exportNamed("$1", $1)'
);

// 转换 export function
const exportedFunctions = [];
let match;
const exportFuncRegex = /export\s+function\s+([a-zA-Z_$][a-zA-Z0-9_$]*)/g;
while ((match = exportFuncRegex.exec(transformedContent)) !== null) {
  exportedFunctions.push(match[1]);
}

println('导出的函数:', exportedFunctions.join(', '));

transformedContent = transformedContent.replace(
  /export\s+function\s+/g,
  'function '
);

if (exportedFunctions.length > 0) {
  const exportCalls = exportedFunctions.map(name => 
    'exportNamed("' + name + '", ' + name + ');'
  ).join('\n');
  transformedContent += '\n' + exportCalls;
}

// 转换 export { ... }
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

println('3. 转换后内容:');
println(transformedContent);

// 测试执行
try {
  const exports = {};
  const exportNamed = (name, value) => { 
    exports[name] = value; 
    println('导出:', name, typeof value);
  };
  const exportDefault = (value) => { 
    exports.default = value; 
  };
  
  eval(transformedContent);
  
  println('4. 最终导出:', JSON.stringify(exports));
  
} catch (error) {
  println('❌ 执行失败:', error.message);
  if (error.stack) {
    println('错误堆栈:', error.stack);
  }
}

println('=== 调试完成 ===');