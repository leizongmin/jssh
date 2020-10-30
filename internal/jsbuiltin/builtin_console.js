const console = {};

{
  function printVars(...vars) {
    println(vars.map((v) => String(v)).join(" "));
  }

  console.log = function log(...args) {
    printVars(...args);
  };

  console.error = function error(...args) {
    printVars(...args);
  };
}
