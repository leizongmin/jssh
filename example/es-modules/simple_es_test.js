#!/usr/bin/env jssh

console.log('=== Simple ES Module Test ===');

// Create a very simple ES module
const simpleModule = `export const value = 42;`;

jssh.fs.writefile('simple.mjs', simpleModule);
console.log('Created simple ES module');

// Test the detection
console.log('Is ES module:', jssh.isESModule(simpleModule));

// Test sync import
try {
  const result = jssh.loadESModuleSync('./simple.mjs');
  console.log('Load successful!');
  console.log('Result:', JSON.stringify(result));
} catch (error) {
  console.error('Load failed:', error.message);
}

// Test with CommonJS for comparison
const cjsModule = `module.exports = { value: 42 };`;
jssh.fs.writefile('simple.js', cjsModule);

try {
  const result2 = jssh.loadESModuleSync('./simple.js');
  console.log('CommonJS load successful!');
  console.log('Result:', JSON.stringify(result2));
} catch (error) {
  console.error('CommonJS load failed:', error.message);
}

console.log('Simple test completed.');