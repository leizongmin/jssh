#!/usr/bin/env jssh

println('=== 手动测试模块转换 ===');

// 加载增强实现
const enhancedModuleCode = jssh.fs.readfile('./enhanced-es-module.js');
eval(enhancedModuleCode);

// 读取简单模块内容
const content = jssh.fs.readfile('./very-simple.mjs');
println('1. 原始模块内容:');
println(content);

// 手动测试转换逻辑
println('2. 开始转换测试...');

try {
  // 模拟转换过程
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
  
  println('3. 转换后的内容:');
  println(transformedContent);
  
  // 创建模块包装器
  const moduleWrapper = `
(function() {
  const exports = {};
  const exportNamed = (name, value) => { 
    exports[name] = value; 
  };
  
  ${transformedContent}
  
  return exports;
})()`;

  println('4. 执行模块...');
  const result = eval(moduleWrapper);
  
  println('5. 模块执行结果:');
  println('- message:', result.message);
  
  println('✅ 手动转换测试成功！');

} catch (error) {
  println('❌ 手动转换测试失败:');
  println('错误信息:', error.message);
  if (error.stack) {
    println('错误堆栈:', error.stack);
  }
}

println('=== 测试完成 ===');