#!/usr/bin/env jssh

println('=== 检查属性描述符 ===');

// 检查 jssh 对象的属性
println('1. jssh 对象信息:');
println('- jssh 类型:', typeof jssh);

// 检查 loadESModuleSync 属性
println('2. loadESModuleSync 属性:');
try {
  const descriptor = Object.getOwnPropertyDescriptor(jssh, 'loadESModuleSync');
  if (descriptor) {
    println('- writable:', descriptor.writable);
    println('- enumerable:', descriptor.enumerable);
    println('- configurable:', descriptor.configurable);
    println('- value 类型:', typeof descriptor.value);
  } else {
    println('- 属性描述符不存在');
  }
} catch (error) {
  println('- 获取属性描述符失败:', error.message);
}

// 检查是否可以定义新属性
println('3. 尝试定义新属性:');
try {
  Object.defineProperty(jssh, 'testProperty', {
    value: function() { return 'test'; },
    writable: true,
    enumerable: true,
    configurable: true
  });
  println('- 定义新属性成功');
  println('- testProperty 类型:', typeof jssh.testProperty);
} catch (error) {
  println('- 定义新属性失败:', error.message);
}

// 尝试强制重新定义 loadESModuleSync
println('4. 尝试强制重新定义:');
try {
  Object.defineProperty(jssh, 'loadESModuleSync', {
    value: function(specifier, referrer) {
      println('🟢 重新定义的函数被调用:', specifier);
      return { redefined: true };
    },
    writable: true,
    enumerable: true,
    configurable: true
  });
  println('- 重新定义成功');
  
  const result = jssh.loadESModuleSync('./simple-function.mjs');
  println('- 调用结果:', JSON.stringify(result));
  
} catch (error) {
  println('- 重新定义失败:', error.message);
}

println('=== 检查完成 ===');