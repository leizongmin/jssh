const fs = {};
{
  fs.readfile = function (path) {
    return jssh.op.file_read(path);
  };

  fs.readdir = function (path) {
    return jssh.op.dir_read(path);
  };

  fs.writefile = function (path, data) {
    return jssh.op.file_write(path, data);
  };

  fs.appendfile = function (path, data) {
    return jssh.op.file_append(path, data);
  };

  fs.stat = function (path) {
    return jssh.op.file_stat(path);
  };

  fs.exist = function (path) {
    return jssh.op.file_exist(path);
  };
}
Object.freeze(fs);
