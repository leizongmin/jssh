#!/usr/bin/env jssh

console.log('=== Synchronous ES Module Test ===');

// Create a simple ES module file
const moduleContent = `
console.log('ES module executing...');
export const message = "Hello from ES module";
export function getValue() {
  return 42;
}
export default {
  name: "test-module",
  version: "1.0.0"
};
console.log('ES module executed');
`;

jssh.fs.writefile('sync-test.mjs', moduleContent);
console.log('Created sync-test.mjs');

// Test the ES module detection
console.log('Is ES module:', jssh.isESModule(moduleContent));

// Try to load the module
console.log('\nTrying to load ES module...');
try {
  // Since the function is async but we're calling it synchronously,
  // let's see what happens
  const result = jssh.loadESModule('./sync-test.mjs', __dirname);
  console.log('loadESModule returned:', typeof result);
  
  // The result might be a Promise, let's check
  if (result && typeof result.then === 'function') {
    console.log('Result is a Promise');
    // We can't properly await in this context, but let's try some workarounds
  } else {
    console.log('Result is not a Promise');
    console.log('Result keys:', Object.keys(result || {}));
    if (result) {
      console.log('Result.message:', result.message);
      console.log('Result.default:', result.default);
    }
  }
} catch (error) {
  console.error('Error loading ES module:', error.message);
  console.error('Error stack:', error.stack);
}

console.log('\nSync test completed.');