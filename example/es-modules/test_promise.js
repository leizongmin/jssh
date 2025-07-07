#!/usr/bin/env jssh

console.log('=== Testing Promise and Async Support ===');

// Create a test module
jssh.fs.writefile('test.mjs', 'export const value = 123; export default "hello";');

// Test 1: Basic Promise
console.log('Testing basic Promise...');
const promise = new Promise((resolve) => {
  resolve('Promise works!');
});

promise.then(result => {
  console.log('Promise result:', result);
}).catch(err => {
  console.error('Promise error:', err);
});

// Test 2: Import function
console.log('Testing import()...');
const importPromise = import('./test.mjs');

importPromise.then(module => {
  console.log('Import successful!');
  console.log('module.value =', module.value);
  console.log('module.default =', module.default);
}).catch(err => {
  console.error('Import failed:', err.message);
  console.error('Error details:', err);
});

// Test 3: Async/await syntax
console.log('Testing async/await...');
(async function() {
  try {
    const mod = await import('./test.mjs');
    console.log('Async import successful!');
    console.log('async module.value =', mod.value);
  } catch (error) {
    console.error('Async import failed:', error.message);
  }
})();

console.log('All tests started...');