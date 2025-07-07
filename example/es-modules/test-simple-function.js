#!/usr/bin/env jssh

println('=== 测试简单函数模块 ===');

// 加载增强实现
const enhancedModuleCode = jssh.fs.readfile('./enhanced-es-module.js');
eval(enhancedModuleCode);

try {
  println('测试简单函数模块...');
  const result = jssh.loadESModuleSync('./simple-function.mjs');
  
  println('✅ 加载成功！');
  println('导出的内容:', Object.keys(result));
  println('hello 函数类型:', typeof result.hello);
  
  if (typeof result.hello === 'function') {
    println('调用 hello():', result.hello());
    println('🎉 简单函数测试成功！');
  }
  
} catch (error) {
  println('❌ 简单函数测试失败:');
  println('错误信息:', error.message);
  if (error.stack) {
    println('错误堆栈:', error.stack);
  }
}

println('=== 测试完成 ===');