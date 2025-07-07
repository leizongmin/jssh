#!/usr/bin/env jssh

// Test ES Module functionality

console.log('Testing ES Module support in jssh...');

// Test 1: Check if import() function is available
console.log('1. Testing import() function availability:');
console.log('import function available:', typeof globalThis.import === 'function');

// Test 2: Create a simple module file for testing
const testModuleContent = `
export const message = "Hello from ES Module!";
export function greet(name) {
  return \`Hello, \${name}!\`;
}
export default {
  version: "1.0.0",
  name: "test-module"
};
`;

jssh.fs.writefile('test-module.mjs', testModuleContent);
console.log('2. Created test module file: test-module.mjs');

// Test 3: Test dynamic import
(async () => {
  try {
    console.log('3. Testing dynamic import...');
    const module = await import('./test-module.mjs');
    console.log('Import successful!');
    console.log('- module.message:', module.message);
    console.log('- module.greet("World"):', module.greet("World"));
    console.log('- module.default:', JSON.stringify(module.default));
  } catch (error) {
    console.error('Dynamic import failed:', error.message);
  }
})();

// Test 4: Test CommonJS compatibility
console.log('4. Testing CommonJS compatibility...');
const testCommonJSContent = `
module.exports = {
  name: "CommonJS Module",
  getValue: () => 42
};
`;

jssh.fs.writefile('test-commonjs.js', testCommonJSContent);

(async () => {
  try {
    const cjsModule = await import('./test-commonjs.js');
    console.log('CommonJS import successful!');
    console.log('- cjsModule.default.name:', cjsModule.default.name);
    console.log('- cjsModule.default.getValue():', cjsModule.default.getValue());
  } catch (error) {
    console.error('CommonJS import failed:', error.message);
  }
})();

console.log('ES Module tests completed!');