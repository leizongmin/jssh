# Top-level Await 实现指南

## 概述

这个项目为 jssh 提供了完整的 **top-level await** 支持。Top-level await 允许您在 ES 模块的顶层直接使用 `await` 关键字，而无需将代码包装在 `async` 函数中。

## 功能特性

✅ **完全支持 top-level await**
- 在 `.mjs` 文件中直接使用 `await`
- 支持异步模块初始化
- 兼容现有的 ES 模块系统

✅ **增强的模块导入**
- 自动检测包含 top-level await 的模块
- 智能转换导出语句
- 保持与 CommonJS 的兼容性

✅ **多种导出语法支持**
- `export const/let/var`
- `export function`
- `export default`
- `export { name1, name2 }`

## 快速开始

### 1. 加载增强支持

首先在您的主脚本中加载增强的 ES 模块支持：

```javascript
#!/usr/bin/env jssh

// 加载增强的 ES 模块支持
const enhancedModuleCode = jssh.fs.readfile('./enhanced-es-module.js');
eval(enhancedModuleCode);

console.log('✅ Top-level await 支持已加载');
```

### 2. 创建使用 top-level await 的模块

创建一个 `.mjs` 文件：

```javascript
// my-async-module.mjs

console.log('模块开始初始化...');

// 1. 基本的 top-level await
const delay = (ms) => new Promise(resolve => setTimeout(resolve, ms));

await delay(1000);
console.log('等待完成！');

// 2. 异步数据获取
const fetchData = async () => {
  await delay(500);
  return { message: 'Hello from async module!', timestamp: Date.now() };
};

const data = await fetchData();

// 3. 导出异步处理的结果
export const message = data.message;
export const timestamp = data.timestamp;
export default data;

console.log('模块初始化完成！');
```

### 3. 导入和使用模块

```javascript
// main.js
async function main() {
  const myModule = await import('./my-async-module.mjs');
  
  console.log('导入的数据:', myModule.message);
  console.log('时间戳:', myModule.timestamp);
  console.log('默认导出:', myModule.default);
}

main();
```

## 完整示例

项目包含了几个完整的示例：

### 示例文件结构

```
example/es-modules/
├── enhanced-es-module.js          # 增强的 ES 模块实现
├── top-level-await-demo.mjs       # Top-level await 演示
├── math.mjs                       # 数学工具模块
├── test-top-level-await.js        # 基本测试
└── run-top-level-await.js         # 完整运行脚本
```

### 运行示例

```bash
# 进入示例目录
cd example/es-modules/

# 运行完整的测试
jssh run-top-level-await.js

# 或者运行基本测试
jssh test-top-level-await.js
```

## 工作原理

### 1. 模块检测

系统会自动检测 `.mjs` 文件中是否包含 top-level await：

```javascript
const hasTopLevelAwait = /(?:^|\n|\r\n?)\s*await\s+/.test(content);
```

### 2. 代码转换

包含 top-level await 的模块会被包装在一个异步函数中：

```javascript
(async function() {
  // 模块代码在这里执行
  // await 语句可以正常工作
  
  // 导出处理
  const exports = {};
  // ... 处理各种导出语法
  
  return exports;
})()
```

### 3. 导出语句转换

各种导出语法会被转换为函数调用：

```javascript
// export const name = value;
// 转换为 ↓
const name = value; exportNamed("name", name);

// export default value;
// 转换为 ↓
exportDefault(value);

// export { name1, name2 };
// 转换为 ↓
exportNamed("name1", name1); exportNamed("name2", name2);
```

## API 参考

### jssh.import(specifier)

异步导入模块，支持 top-level await：

```javascript
const module = await jssh.import('./my-module.mjs');
```

### jssh.loadESModule(specifier, referrer)

底层的 ES 模块加载函数：

```javascript
const module = await jssh.loadESModule('./my-module.mjs', __dirname);
```

## 支持的语法

### ✅ 支持的导出语法

```javascript
// 常量导出
export const PI = 3.14159;
export let counter = 0;
export var name = 'test';

// 函数导出
export function add(a, b) {
  return a + b;
}

// 默认导出
export default { message: 'Hello' };
export default 42;
export default function() { /* ... */ }

// 批量导出
export { name1, name2 };
export { name1 as alias1, name2 };
```

### ✅ 支持的 await 用法

```javascript
// 基本 await
await delay(1000);

// await 赋值
const data = await fetchData();

// await 在表达式中
const result = await processData(await getData());

// await 导入
const module = await import('./other-module.mjs');
```

## 注意事项

1. **文件扩展名**: 推荐使用 `.mjs` 扩展名来标识 ES 模块
2. **性能**: Top-level await 会使模块加载变为异步，可能影响启动性能
3. **错误处理**: 在 top-level await 中的错误会阻止模块加载
4. **依赖顺序**: 包含 top-level await 的模块会延迟其导入者的执行

## 故障排除

### 常见问题

**Q: 模块导入失败**
```
A: 确保文件路径正确，使用 .mjs 扩展名，并且已加载增强模块支持
```

**Q: await 不工作**
```
A: 检查是否在 .mjs 文件中，并确保 enhanced-es-module.js 已正确加载
```

**Q: 导出不正确**
```
A: 检查导出语法是否符合支持的格式，避免使用复杂的导出表达式
```

### 调试技巧

启用调试日志：

```javascript
// 在脚本开头添加
jssh.log.level = 'debug';
```

这将显示模块加载和转换的详细信息。

## 贡献

如果您发现问题或有改进建议，请：

1. 检查现有的示例是否涵盖了您的用例
2. 创建最小的复现示例
3. 提交 issue 或 pull request

## 总结

这个 top-level await 实现为 jssh 提供了现代 JavaScript 的异步模块支持，让您可以：

- 在模块顶层直接使用 `await`
- 进行异步模块初始化
- 保持代码的简洁性和可读性
- 与现有的模块系统无缝集成

开始使用 top-level await 来简化您的异步代码吧！