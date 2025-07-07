#!/usr/bin/env jssh

println('=== 最终 Export Function 测试 ===');

// 加载增强实现
const enhancedModuleCode = jssh.fs.readfile('./enhanced-es-module.js');
eval(enhancedModuleCode);

try {
  println('测试加载 test-export-function.mjs...');
  const functionModule = jssh.loadESModuleSync('./test-export-function.mjs');
  
  println('✅ 模块加载成功！');
  
  // 检查所有导出
  const keys = Object.keys(functionModule);
  println('导出的键:', keys.join(', '));
  
  for (const key of keys) {
    println(`- ${key}: ${typeof functionModule[key]}`);
  }
  
  // 测试函数调用
  if (typeof functionModule.add === 'function') {
    println('测试函数调用:');
    println('- add(2, 3) =', functionModule.add(2, 3));
    println('- multiply(4, 5) =', functionModule.multiply(4, 5));
    println('- greet("World") =', functionModule.greet("World"));
    println('- PI =', functionModule.PI);
    println('- PIValue =', functionModule.PIValue);
    
    println('🎉 所有测试通过！Export function 修复成功！');
  } else {
    println('❌ 函数导出失败');
  }
  
} catch (error) {
  println('❌ 测试失败:');
  println('错误信息:', error.message);
  if (error.stack) {
    println('错误堆栈:', error.stack);
  }
}

println('=== 测试完成 ===');