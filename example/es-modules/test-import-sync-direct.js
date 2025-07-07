#!/usr/bin/env jssh

println('=== 直接测试 jssh.importSync ===');

// 加载增强实现
const enhancedModuleCode = jssh.fs.readfile('./enhanced-es-module.js');
eval(enhancedModuleCode);

println('1. 测试 jssh.loadESModuleSync...');
try {
  const result1 = jssh.loadESModuleSync('./test-export-function.mjs');
  println('loadESModuleSync 成功:', JSON.stringify(result1));
} catch (error) {
  println('loadESModuleSync 失败:', error.message);
}

println('2. 测试 jssh.importSync...');
try {
  const result2 = jssh.importSync('./test-export-function.mjs');
  println('importSync 成功:', JSON.stringify(result2));
} catch (error) {
  println('importSync 失败:', error.message);
  if (error.stack) {
    println('错误堆栈:', error.stack);
  }
}

println('=== 测试完成 ===');