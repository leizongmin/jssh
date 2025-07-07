#!/usr/bin/env jssh

println('=== 调试增强函数调用 ===');

// 检查原始函数
println('1. 原始函数:');
println('- jssh.importSync 类型:', typeof jssh.importSync);
println('- jssh.loadESModuleSync 类型:', typeof jssh.loadESModuleSync);

// 加载增强实现
println('2. 加载增强实现...');
const enhancedModuleCode = jssh.fs.readfile('./enhanced-es-module.js');
eval(enhancedModuleCode);

// 检查增强后的函数
println('3. 增强后的函数:');
println('- jssh.importSync 类型:', typeof jssh.importSync);
println('- jssh.loadESModuleSync 类型:', typeof jssh.loadESModuleSync);

// 测试函数是否被替换
println('4. 测试函数调用...');
try {
  // 直接调用我们的增强函数
  println('调用 jssh.loadESModuleSync...');
  const result = jssh.loadESModuleSync('./very-simple.mjs');
  println('结果:', JSON.stringify(result));
} catch (error) {
  println('错误:', error.message);
}

println('=== 调试完成 ===');