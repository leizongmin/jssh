#!/usr/bin/env jssh

println('=== 测试函数转换逻辑 ===');

// 测试原始的 export function 代码
const testCode = `export function add(a, b) {
  return a + b;
}

export function greet(name) {
  return "Hello, " + name;
}`;

println('原始代码:');
println(testCode);

// 模拟转换过程
let transformedContent = testCode;

// 应用我们的新转换规则
// 收集导出的函数名
const exportedFunctions = [];
let match;
const exportFuncRegex = /export\s+function\s+([a-zA-Z_$][a-zA-Z0-9_$]*)/g;
while ((match = exportFuncRegex.exec(transformedContent)) !== null) {
  exportedFunctions.push(match[1]);
}

// 移除 export 关键字
transformedContent = transformedContent.replace(
  /export\s+function\s+/g,
  'function '
);

// 在转换后的内容末尾添加导出调用
if (exportedFunctions.length > 0) {
  const exportCalls = exportedFunctions.map(name => 
    'exportNamed("' + name + '", ' + name + ');'
  ).join('\n');
  transformedContent += '\n' + exportCalls;
}

println('\n转换后代码:');
println(transformedContent);

// 测试执行
try {
  const exports = {};
  const exportNamed = (name, value) => { 
    exports[name] = value; 
    println('导出函数:', name);
  };
  
  eval(transformedContent);
  
  println('\n导出结果:');
  println('- add 函数:', typeof exports.add);
  println('- greet 函数:', typeof exports.greet);
  
  if (typeof exports.add === 'function') {
    println('- add(2, 3) =', exports.add(2, 3));
  }
  if (typeof exports.greet === 'function') {
    println('- greet("World") =', exports.greet("World"));
  }
  
  println('✅ 函数转换测试成功！');
  
} catch (error) {
  println('❌ 函数转换测试失败:');
  println('错误信息:', error.message);
}

println('=== 测试完成 ===');