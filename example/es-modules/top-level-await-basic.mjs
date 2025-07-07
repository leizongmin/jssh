// Top-level await 基本演示

console.log('开始执行 top-level await...');

// 创建一个简单的 Promise
const delay = (ms) => new Promise(resolve => setTimeout(resolve, ms));

// 使用 top-level await
await delay(1000);
console.log('等待 1 秒完成！');

// 异步获取数据
const getData = async () => {
  await delay(500);
  return { message: 'Hello from top-level await!', timestamp: Date.now() };
};

const data = await getData();
console.log('获取的数据:', data.message);

// 导出结果
const result = data;
export { result };