# jssh ES Modules Examples

这个目录包含了jssh ES模块功能的示例和测试文件。

## 主要文件说明

### 📖 文档
- `ES_MODULE_IMPLEMENTATION_SUMMARY.md` - ES模块实现的详细文档

### 🎯 演示文件
- `es_module_demo.js` - ES模块功能的完整演示

### 🧪 测试文件
- `test_esmodule.js` - ES模块功能的综合测试
- `simple_es_test.js` - 简单的ES模块测试
- `debug_test.js` - 调试ES模块功能
- `test_sync.js` - 同步ES模块测试
- `sync_test.js` - 同步功能测试
- `test_promise.js` - Promise和异步功能测试
- `simple_test.js` - 基础功能测试
- `direct_test.js` - 直接调用测试

### 📦 示例模块
- `math.mjs` - 数学函数ES模块示例
- `utils.js` - CommonJS工具模块示例
- `test-module.mjs` - 测试用ES模块
- `test-commonjs.js` - 测试用CommonJS模块
- `simple.mjs` / `simple.js` - 简单的测试模块
- `sync-test.mjs` - 同步测试模块
- `direct.mjs` - 直接测试模块
- `test.mjs` - 基础测试模块
- `debug-es.mjs` / `debug-js.js` - 调试用模块

## 运行示例

```bash
# 运行主要演示
./release/jssh example/es-modules/es_module_demo.js

# 运行具体测试
./release/jssh example/es-modules/test_esmodule.js
./release/jssh example/es-modules/simple_es_test.js
```

## 功能特性

✅ **已实现的功能**
- ES模块检测和识别
- 动态import()函数支持
- 同步和异步模块加载
- CommonJS模块兼容性
- 模块缓存机制
- 多种路径解析支持

📝 **使用建议**
- 推荐使用 `import()` 或 `jssh.import()` 进行模块导入
- ES模块文件建议使用 `.mjs` 扩展名
- 可以在ES模块中导入CommonJS模块
- 支持相对路径、绝对路径和HTTP URL

更多详细信息请查看 `ES_MODULE_IMPLEMENTATION_SUMMARY.md` 文档。