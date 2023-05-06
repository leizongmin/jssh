{
  const readfilebytes = fs.readfilebytes;
  fs.readfilebytes = function (filename) {
    return Uint8Array.from(readfilebytes(filename));
  };
}

Object.freeze(fs);
