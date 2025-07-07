#!/usr/bin/env jssh

console.log('=== 简单的 ES 模块功能测试 ===');

// 测试基本的 jssh 功能
console.log('jssh 版本信息:');
console.log('- import 支持:', typeof jssh.import);
console.log('- loadESModule 支持:', typeof jssh.loadESModule);
console.log('- isESModule 支持:', typeof jssh.isESModule);

// 测试一个简单的模块导入
async function testBasicImport() {
  try {
    console.log('\n测试基本模块导入...');
    const mathModule = await jssh.import('./math.mjs');
    console.log('导入成功!');
    console.log('- PI 值:', mathModule.PI);
    console.log('- E 值:', mathModule.E);
    console.log('- 2 + 3 =', mathModule.add(2, 3));
    console.log('- 4 * 5 =', mathModule.multiply(4, 5));
    console.log('✅ 基本模块导入测试通过');
  } catch (error) {
    console.log('❌ 基本模块导入测试失败:', error.message);
  }
}

testBasicImport();

console.log('\n=== 测试完成 ===');