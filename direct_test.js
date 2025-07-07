#!/usr/bin/env jssh

console.log('=== Direct Import Test ===');

// Create a test module
jssh.fs.writefile('direct.mjs', 'export const name = "direct test"; export default 999;');

console.log('Created direct.mjs');

// Test direct call
console.log('Testing jssh.import directly...');

try {
  const result = jssh.import('./direct.mjs');
  console.log('jssh.import result type:', typeof result);
  console.log('jssh.import result:', result);
} catch (error) {
  console.error('jssh.import failed:', error.message);
  console.error('Error stack:', error.stack);
}

// Test if loadESModule exists
console.log('jssh.loadESModule exists:', typeof jssh.loadESModule);

if (typeof jssh.loadESModule === 'function') {
  try {
    const result2 = jssh.loadESModule('./direct.mjs');
    console.log('jssh.loadESModule result type:', typeof result2);
    console.log('jssh.loadESModule result:', result2);
  } catch (error) {
    console.error('jssh.loadESModule failed:', error.message);
  }
}

console.log('Direct test completed.');