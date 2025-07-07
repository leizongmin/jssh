
console.log('ES module executing...');
export const message = "Hello from ES module";
export function greet(name) {
  return "Hello, " + name + "!";
}
export default {
  name: "test-module",
  version: "1.0.0"
};
console.log('ES module executed');
