// const xss = require("https://raw.githubusercontent.com/leizongmin/js-xss/master/dist/xss.js");
// console.log(Object.keys(xss));
// console.log(Object.keys(global));

const mod = require("./test_mod");
mod.hello("hh");

const xss = require("xss");
console.log(xss("<script>alert(1)</script>"));

const cssfilter = require("cssfilter");
console.log(cssfilter("a:1; b:2; width:10;"));

const lodash = require("lodash");
console.log(lodash.partition([1, 2, 3, 4], (n) => n % 2));

// console.log("abs", path.abs("https://unpkg.com/xss@1.0.9/lib/index.js"));
// console.log("base", path.base("https://unpkg.com/xss@1.0.9/lib/index.js"));
// console.log("join", path.join("https://unpkg.com/xss@1.0.9/lib", "index.js"));
// console.log("dir", path.dir("https://unpkg.com/xss@1.0.9/lib/index.js"));
const xss2 = require("https://unpkg.com/xss@1.0.9/lib/index.js");
console.log(xss2("<script>alert(1)</script>"));

const cssfilter2 = require("https://unpkg.com/cssfilter@0.0.10/lib/index.js");
console.log(cssfilter2("a:1; b:2; width:10;"));

const lodash2 = require("https://unpkg.com/lodash@4.17.21/lodash.js");
console.log(lodash2.partition([1, 2, 3, 4], (n) => n % 2));

const lodash3 = require("https://unpkg.com/lodash");
console.log(lodash3.partition([1, 2, 3, 4], (n) => n % 2));
