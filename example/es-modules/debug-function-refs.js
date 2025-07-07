#!/usr/bin/env jssh

println('=== 调试函数引用 ===');

// 1. 保存原始引用
const originalLoadESModuleSync = jssh.loadESModuleSync;
println('1. 原始函数保存');

// 2. 创建一个简单的测试函数
function testFunction(specifier, referrer) {
  println('🟢 测试函数被调用:', specifier);
  return { test: 'success' };
}

// 3. 尝试替换
println('2. 尝试替换函数...');
jssh.loadESModuleSync = testFunction;

// 4. 检查替换是否成功
println('3. 检查替换结果:');
println('- 函数相同?', jssh.loadESModuleSync === originalLoadESModuleSync);
println('- 函数是测试函数?', jssh.loadESModuleSync === testFunction);
println('- 当前函数:', jssh.loadESModuleSync.toString().substring(0, 100));

// 5. 尝试调用
println('4. 尝试调用函数...');
try {
  const result = jssh.loadESModuleSync('./simple-function.mjs');
  println('- 调用结果:', JSON.stringify(result));
} catch (error) {
  println('- 调用失败:', error.message);
}

// 6. 尝试通过 importSync 调用
println('5. 尝试通过 importSync 调用...');
try {
  const result2 = jssh.importSync('./simple-function.mjs');
  println('- importSync 结果:', JSON.stringify(result2));
} catch (error) {
  println('- importSync 失败:', error.message);
}

println('=== 调试完成 ===');