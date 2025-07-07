#!/usr/bin/env jssh

println('=== 调试测试开始 ===');

// 检查 jssh 基本功能
println('1. 检查 jssh 对象:');
println('typeof jssh:', typeof jssh);

if (typeof jssh !== 'undefined') {
  println('2. 检查 ES 模块相关函数:');
  println('typeof jssh.import:', typeof jssh.import);
  println('typeof jssh.loadESModule:', typeof jssh.loadESModule);
  println('typeof jssh.isESModule:', typeof jssh.isESModule);
}

// 测试读取文件
println('3. 测试读取 math.mjs 文件:');
try {
  const content = jssh.fs.readfile('./math.mjs');
  println('文件内容长度:', content.length);
  println('是否为 ES 模块:', jssh.isESModule(content));
} catch (error) {
  println('读取文件失败:', error.message);
}

// 测试异步导入
println('4. 测试异步导入:');
jssh.import('./math.mjs').then(module => {
  println('导入成功! PI 值:', module.PI);
}).catch(error => {
  println('导入失败:', error.message);
});

println('=== 调试测试完成 ===');