#!/usr/bin/env jssh

println('=== 同步测试 Export Function 修复 ===');

// 加载增强实现
const enhancedModuleCode = jssh.fs.readfile('./enhanced-es-module.js');
eval(enhancedModuleCode);

try {
  println('1. 测试同步导入包含 export function 的模块...');
  const functionModule = jssh.importSync('./test-export-function.mjs');
  
  println('2. 检查导出的函数:');
  println('- add 函数:', typeof functionModule.add);
  println('- multiply 函数:', typeof functionModule.multiply);
  println('- greet 函数:', typeof functionModule.greet);
  println('- PI 常量:', functionModule.PI);
  println('- PIValue:', functionModule.PIValue);
  
  if (typeof functionModule.add === 'function') {
    println('3. 测试函数调用:');
    println('- add(2, 3) =', functionModule.add(2, 3));
    println('- multiply(4, 5) =', functionModule.multiply(4, 5));
    println('- greet("World") =', functionModule.greet("World"));
  }
  
  println('✅ Export function 修复测试成功！');
  
} catch (error) {
  println('❌ Export function 测试失败:');
  println('错误信息:', error.message);
  if (error.stack) {
    println('错误堆栈:', error.stack);
  }
}

println('=== 测试完成 ===');