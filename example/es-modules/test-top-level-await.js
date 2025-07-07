#!/usr/bin/env jssh

console.log('=== 测试 Top-level Await 功能 ===\n');

// 测试函数
async function testTopLevelAwait() {
  try {
    console.log('正在导入 top-level-await-demo.mjs...');
    const demoModule = await import('./top-level-await-demo.mjs');
    
    console.log('\n导入成功！模块导出的内容:');
    console.log('- message:', demoModule.message);
    console.log('- completedAt:', demoModule.completedAt);
    console.log('- default export:', JSON.stringify(demoModule.default, null, 2));
    
    console.log('\n✅ Top-level await 测试通过！');
    
  } catch (error) {
    console.error('\n❌ Top-level await 测试失败:');
    console.error('错误信息:', error.message);
    console.error('错误详情:', error);
  }
}

// 执行测试
testTopLevelAwait();

console.log('\n=== 测试完成 ===');