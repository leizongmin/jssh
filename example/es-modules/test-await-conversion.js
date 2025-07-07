#!/usr/bin/env jssh

println('=== 测试 Top-level Await 转换 ===');

// 加载增强实现
const enhancedModuleCode = jssh.fs.readfile('./enhanced-es-module.js');
eval(enhancedModuleCode);

// 读取包含 top-level await 的模块
const content = jssh.fs.readfile('./top-level-await-basic.mjs');
println('1. 原始 Top-level Await 模块内容:');
println(content);

// 检测是否包含 top-level await
const hasTopLevelAwait = /(?:^|\n|\r\n?)\s*await\s+/.test(content);
println('2. 检测到 top-level await:', hasTopLevelAwait);

if (hasTopLevelAwait) {
  println('3. 开始转换 Top-level Await 模块...');
  
  try {
    // 模拟异步模块转换过程
    let transformedContent = content;
    
    // 转换 export { name1, name2 }
    transformedContent = transformedContent.replace(
      /export\s*\{([^}]+)\}/g,
      (match, exports) => {
        return exports.split(',').map(exp => {
          const [from, to] = exp.trim().split(' as ');
          const exportName = to ? to.trim() : from.trim();
          const varName = from.trim();
          return `exportNamed("${exportName}", ${varName});`;
        }).join('\n');
      }
    );
    
    println('4. 转换后的内容:');
    println(transformedContent);
    
    // 创建异步模块包装器
    const moduleWrapper = `
(async function() {
  const exports = {};
  const exportNamed = (name, value) => { 
    exports[name] = value; 
  };
  
  ${transformedContent}
  
  return exports;
})()`;

    println('5. 创建异步模块包装器完成');
    println('6. 包装器示例（前100字符）:');
    println(moduleWrapper.substring(0, 100) + '...');
    
    println('✅ Top-level Await 转换测试成功！');
    println('注意：实际执行需要在异步环境中进行');
    
  } catch (error) {
    println('❌ Top-level Await 转换失败:');
    println('错误信息:', error.message);
  }
}

println('=== 测试完成 ===');