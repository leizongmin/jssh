// 数学工具模块

export const PI = 3.14159265359;
export const E = 2.71828182846;

export function add(a, b) {
  return a + b;
}

export function multiply(a, b) {
  return a * b;
}

export function power(base, exponent) {
  return Math.pow(base, exponent);
}

export default {
  PI,
  E,
  add,
  multiply,
  power
};
