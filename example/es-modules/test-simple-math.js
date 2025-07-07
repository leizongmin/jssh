#!/usr/bin/env jssh

println('=== 测试简单数学模块 ===');

try {
  println('1. 使用同步导入加载 simple-math.mjs...');
  const mathModule = jssh.importSync('./simple-math.mjs');
  
  println('2. 导入成功！检查导出内容:');
  println('- PI 值:', mathModule.PI);
  println('- E 值:', mathModule.E);
  println('- add 函数:', typeof mathModule.add);
  println('- multiply 函数:', typeof mathModule.multiply);
  
  if (typeof mathModule.add === 'function') {
    println('3. 测试函数调用:');
    println('- 2 + 3 =', mathModule.add(2, 3));
    println('- 4 * 5 =', mathModule.multiply(4, 5));
  }
  
  println('✅ 简单模块测试成功！');
  
} catch (error) {
  println('❌ 简单模块测试失败:');
  println('错误信息:', error.message);
  if (error.stack) {
    println('错误堆栈:', error.stack);
  }
}

println('=== 测试完成 ===');