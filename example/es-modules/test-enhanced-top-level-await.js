#!/usr/bin/env jssh

println('=== 测试增强的 Top-level Await 实现 ===');

// 首先加载增强的 ES 模块支持
try {
  println('1. 加载增强的 ES 模块实现...');
  const enhancedModuleCode = jssh.fs.readfile('./enhanced-es-module.js');
  eval(enhancedModuleCode);
  println('✅ 增强模块加载成功');
} catch (error) {
  println('❌ 增强模块加载失败:', error.message);
  process.exit(1);
}

// 测试增强的实现
async function testEnhancedModule() {
  try {
    println('2. 测试增强的导入功能...');
    
    // 先测试简单模块
    const simpleModule = await jssh.import('./very-simple.mjs');
    println('简单模块导入成功:', simpleModule.message);
    
    // 测试包含 top-level await 的模块
    println('3. 测试 top-level await 模块...');
    const awaitModule = await jssh.import('./top-level-await-basic.mjs');
    println('Top-level await 模块导入成功!');
    println('结果:', awaitModule.result);
    
    println('✅ 所有测试通过！');
    
  } catch (error) {
    println('❌ 测试失败:');
    println('错误信息:', error.message);
    if (error.stack) {
      println('错误堆栈:', error.stack);
    }
  }
}

// 运行测试
testEnhancedModule().then(() => {
  println('=== 测试完成 ===');
}).catch(error => {
  println('❌ 测试执行失败:', error.message);
});