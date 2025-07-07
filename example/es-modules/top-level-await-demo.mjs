// Top-level await 演示文件

console.log('开始执行 top-level await 演示...');

// 1. 基本的 top-level await
const delay = (ms) => new Promise(resolve => setTimeout(resolve, ms));

console.log('等待 1 秒...');
await delay(1000);
console.log('1 秒等待完成！');

// 2. 使用 top-level await 导入其他模块
console.log('动态导入模块...');
const mathModule = await import('./math.mjs');
console.log('导入的 PI 值:', mathModule.PI);

// 3. 使用 top-level await 处理异步操作
const fetchData = async () => {
  return new Promise(resolve => {
    setTimeout(() => {
      resolve({ data: 'Hello from async function', timestamp: Date.now() });
    }, 500);
  });
};

console.log('获取异步数据...');
const result = await fetchData();
console.log('获取到的数据:', result);

// 4. 导出一些值供其他模块使用
export const message = 'This module used top-level await!';
export const completedAt = new Date().toISOString();

// 5. 默认导出一个经过异步处理的对象
const processedData = await (async () => {
  await delay(200);
  return {
    status: 'completed',
    result: result,
    mathPI: mathModule.PI
  };
})();

export default processedData;

console.log('top-level await 演示完成！');