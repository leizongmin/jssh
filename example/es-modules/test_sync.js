#!/usr/bin/env jssh

console.log('=== Testing Synchronous ES Module ===');

// Check if sync function exists
console.log('jssh.importSync exists:', typeof jssh.importSync);
console.log('jssh.loadESModuleSync exists:', typeof jssh.loadESModuleSync);

// Create a simple ES module
const moduleContent = `
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
`;

jssh.fs.writefile('sync-test.mjs', moduleContent);
console.log('Created sync-test.mjs');

// Test sync import
console.log('\nTesting sync import...');
try {
  const result = jssh.importSync('./sync-test.mjs');
  console.log('Import successful!');
  console.log('Result type:', typeof result);
  console.log('Result keys:', Object.keys(result || {}));
  
  if (result) {
    console.log('- message:', result.message);
    if (typeof result.greet === 'function') {
      console.log('- greet("World"):', result.greet("World"));
    } else {
      console.log('- greet is not a function, type:', typeof result.greet);
    }
    console.log('- default:', JSON.stringify(result.default));
  }
} catch (error) {
  console.error('Sync import failed:', error.message);
  console.error('Error stack:', error.stack);
}

console.log('\nSync test completed.');