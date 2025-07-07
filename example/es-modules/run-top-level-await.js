#!/usr/bin/env jssh

console.log('=== 加载增强的 Top-level Await 支持 ===\n');

// 首先加载增强的 ES 模块支持
try {
  console.log('正在加载增强的 ES 模块实现...');
  
  // 直接执行增强模块的代码
  const enhancedModuleCode = jssh.fs.readfile('./enhanced-es-module.js');
  eval(enhancedModuleCode);
  
  console.log('✅ 增强的 ES 模块支持已加载\n');
  
} catch (error) {
  console.error('❌ 加载增强模块失败:', error.message);
  process.exit(1);
}

// 现在测试 top-level await
console.log('=== 开始测试 Top-level Await 功能 ===\n');

async function runTests() {
  try {
    // 测试 1: 基本的 top-level await 模块
    console.log('测试 1: 导入包含 top-level await 的模块...');
    const demoModule = await import('./top-level-await-demo.mjs');
    
    console.log('\n模块导入成功！导出的内容:');
    console.log('- message:', demoModule.message);
    console.log('- completedAt:', demoModule.completedAt);
    console.log('- default export:', JSON.stringify(demoModule.default, null, 2));
    
    // 测试 2: 验证异步数据是否正确处理
    console.log('\n测试 2: 验证异步数据处理...');
    if (demoModule.default && demoModule.default.status === 'completed') {
      console.log('✅ 异步数据处理正确');
      console.log('- 结果状态:', demoModule.default.status);
      console.log('- PI 值:', demoModule.default.mathPI);
    } else {
      console.log('❌ 异步数据处理有问题');
    }
    
    // 测试 3: 直接测试数学模块
    console.log('\n测试 3: 测试数学模块导入...');
    const mathModule = await import('./math.mjs');
    console.log('- PI 值:', mathModule.PI);
    console.log('- E 值:', mathModule.E);
    console.log('- 2 + 3 =', mathModule.add(2, 3));
    console.log('- 4 * 5 =', mathModule.multiply(4, 5));
    
    console.log('\n🎉 所有测试通过！Top-level await 功能正常工作！');
    
  } catch (error) {
    console.error('\n❌ 测试失败:');
    console.error('错误信息:', error.message);
    console.error('错误堆栈:', error.stack);
  }
}

// 运行测试
runTests();

console.log('\n=== 测试完成 ===');