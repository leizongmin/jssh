#!/usr/bin/env jssh

println('=== 测试函数替换 ===');

// 检查原始函数
println('1. 原始函数是否存在:');
println('- jssh.loadESModuleSync:', typeof jssh.loadESModuleSync);

// 保存原始函数引用
const originalFunc = jssh.loadESModuleSync;

// 加载增强实现
println('2. 加载增强实现...');
const enhancedModuleCode = jssh.fs.readfile('./enhanced-es-module.js');
eval(enhancedModuleCode);

// 检查函数是否被替换
println('3. 检查函数是否被替换:');
println('- 函数相同?', jssh.loadESModuleSync === originalFunc);
println('- 新函数类型:', typeof jssh.loadESModuleSync);

// 尝试直接创建一个简单的替换来测试
println('4. 创建简单的测试替换...');
jssh.loadESModuleSync = function(specifier, referrer) {
  println('🔍 自定义函数被调用:', specifier);
  return { test: 'custom function called' };
};

try {
  const result = jssh.loadESModuleSync('./simple-function.mjs');
  println('5. 测试结果:', JSON.stringify(result));
} catch (error) {
  println('5. 测试失败:', error.message);
}

println('=== 测试完成 ===');