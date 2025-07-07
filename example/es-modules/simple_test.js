#!/usr/bin/env jssh

console.log('=== Testing ES Module Support ===');

// Check if import function exists
console.log('import function exists:', !!globalThis.import);
console.log('jssh.import exists:', !!jssh.import);

// Create a simple ES module
const moduleContent = 'export const hello = "world"; export default 42;';
jssh.fs.writefile('simple.mjs', moduleContent);
console.log('Created simple.mjs');

// Test import function
try {
  const result = import('./simple.mjs');
  console.log('import() returned:', typeof result);
  
  if (result && typeof result.then === 'function') {
    console.log('import() returned a Promise - this is correct!');
    
    result.then(module => {
      console.log('Module loaded successfully:');
      console.log('- hello:', module.hello);
      console.log('- default:', module.default);
    }).catch(err => {
      console.error('Module loading failed:', err.message);
    });
  } else {
    console.log('import() did not return a Promise');
  }
} catch (error) {
  console.error('import() failed:', error.message);
}

console.log('Test completed.');