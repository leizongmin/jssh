{const readfilebytes=jssh.fs.readfilebytes;jssh.fs.readfilebytes=function(filename){return Uint8Array.from(readfilebytes(filename))},Object.freeze(jssh.fs)}
