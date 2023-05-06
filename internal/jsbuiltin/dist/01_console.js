const console={};{const printVars=(...vars)=>{println(vars.map(v=>String(v)).join(" "))};console.log=function(...args){printVars(...args)},console.error=function(...args){printVars(...args)}}
