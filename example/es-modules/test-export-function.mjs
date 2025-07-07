// 测试 export function 语法修复

export function add(a, b) {
  return a + b;
}

export function multiply(x, y) {
  return x * y;
}

export function greet(name) {
  return `Hello, ${name}!`;
}

// 测试混合导出
export const PI = 3.14159;
export { PI as PIValue };