#!/usr/bin/env jssh

console.log('=== Debug ES Module Test ===');

// Check what functions are available
console.log('Available ES module functions:');
console.log('- jssh.import:', typeof jssh.import);
console.log('- jssh.loadESModule:', typeof jssh.loadESModule); 
console.log('- jssh.isESModule:', typeof jssh.isESModule);

// Test ES module detection
console.log('\nTesting ES module detection:');
const esContent = 'export const test = 42;';
const jsContent = 'const test = 42;';

if (typeof jssh.isESModule === 'function') {
  console.log('ES content is ES module:', jssh.isESModule(esContent));
  console.log('JS content is ES module:', jssh.isESModule(jsContent));
} else {
  console.log('jssh.isESModule not available');
}

// Create test files
jssh.fs.writefile('debug-es.mjs', esContent);
jssh.fs.writefile('debug-js.js', jsContent);

console.log('\nCreated test files');

// Test file existence
console.log('debug-es.mjs exists:', jssh.fs.exist('debug-es.mjs'));
console.log('debug-js.js exists:', jssh.fs.exist('debug-js.js'));

// Read back content
console.log('\nFile contents:');
console.log('debug-es.mjs:', jssh.fs.readfile('debug-es.mjs'));
console.log('debug-js.js:', jssh.fs.readfile('debug-js.js'));

// Test current directory
console.log('\nCurrent directory info:');
console.log('__dirname:', __dirname);
console.log('Current working directory files:');
const files = jssh.fs.readdir('.');
for (const file of files) {
  if (file.name.endsWith('.mjs') || file.name.endsWith('.js')) {
    console.log('- ' + file.name + ' (' + (file.isdir ? 'dir' : 'file') + ')');
  }
}

console.log('\nDebug test completed.');