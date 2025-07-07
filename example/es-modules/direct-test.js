#!/usr/bin/env jssh

println('=== 直接测试增强实现 ===');

// 加载增强实现
const enhancedModuleCode = jssh.fs.readfile('./enhanced-es-module.js');
eval(enhancedModuleCode);

println('1. 增强实现已加载');

// 测试文件检测
const simpleContent = jssh.fs.readfile('./very-simple.mjs');
println('2. 简单模块内容长度:', simpleContent.length);

const awaitContent = jssh.fs.readfile('./top-level-await-basic.mjs');
println('3. Await 模块内容长度:', awaitContent.length);

// 检测 top-level await
const hasAwait = /(?:^|\n|\r\n?)\s*await\s+/.test(awaitContent);
println('4. 检测到 top-level await:', hasAwait);

// 尝试手动调用增强的导入
println('5. 尝试导入简单模块...');
jssh.loadESModule('./very-simple.mjs', __dirname).then(result => {
  println('简单模块导入结果:', JSON.stringify(result));
}).catch(error => {
  println('简单模块导入失败:', error.message);
});

println('=== 测试完成 ===');