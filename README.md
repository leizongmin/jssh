# jssh

使用 JavaScript 编写运维脚本。

## 安装

安装 Linux 版本：

```bash
curl -O -L https://cdn.itoutiao.co/jssh/v0.1/jssh-linux.tar.gz && \
tar -xvf jssh-linux.tar.gz && \
cp jssh /usr/local/bin/jssh && \
chmod +x /usr/local/bin/jssh
```

安装 macOS 版本：

```bash
curl -O -L https://cdn.itoutiao.co/jssh/v0.1/jssh-osx.tar.gz && \
tar -xvf jssh-osx.tar.gz && \
cp jssh /usr/local/bin/jssh && \
chmod +x /usr/local/bin/jssh
```

通过 Go 命令行工具安装：

```bash
GOPROXY=https://goproxy.cn go get -u github.com/leizongmin/jssh
```

下载压缩包：

- macOS(amd64): https://cdn.itoutiao.co/jssh/v0.1/jssh-osx.tar.gz
- Linux(amd64): https://cdn.itoutiao.co/jssh/v0.1/jssh-linux.tar.gz

## 命令行工具使用

- 执行脚本：`jssh file.js`；
- 进入 REPL：`jssh -i`；
- 执行命令行参数指定的脚本代码：`jssh -c "js code"`；
- 执行命令行参数指定的脚本代码，并将结果作为字符串输出：`jssh -x "js code"`；

## 参考文档

TypeScript types 定义参考文件 [jssh.d.ts](https://github.com/leizongmin/jssh/blob/main/jssh.d.ts)。

#### 全局变量列表

- `__version`：jssh 版本号；
- `__bin`：jssh 二进制文件路径；
- `__pid`：当前进程 PID；
- `__tmpdir`：临时目录；
- `__homedir`：用户 Home 目录；
- `__user`：当前用户名；
- `__hostname`：当前主机名；
- `__dirname`：当前脚本文件的目录；
- `__filename`：当前脚本完整文件名；
- `__args`：当前进程命令行参数；
- `__env`：环境变量；
- `__output`：上一个命令输出的结果，`sh.exec()`且`mode=1`或`mode=2`时有效；
- `__outputbytes`：上一个命令输出结果的字节数；
- `__code`：上一个命令结束时的状态码；

#### 全局函数列表

- `set(name, value)`：设置全局变量；
- `get(name)`：获取全局变量；
- `format(template, ...args)`：格式化字符串，如`format("a=%d, b=%s", 123, "xxx")`；
- `print(template, ...args)`：格式化字符串并输出；
- `println(template, ...args)`：格式化字符串并输出，末尾加换行符；
- `readline()`：从控制台获取用户一行的字符串输入；
- `sleep(ms)`：等待指定毫秒时间；
- `exit(code)`：结束进程；
- `loadconfig(filename, format?)`：加载配置文件，支持 JSON、YAML、TOML 格式；
- `base64encode`：Base64 编码字符串；
- `base64decode`：Base64 解码字符串；
- `md5`：MD5 编码字符串；
- `sha1`：SAH1 编码字符串；
- `sha256`：SHA256 编码字符串；

#### Shell 操作

- `sh.setenv(name, value)`：设置环境变量；
- `sh.exec(cmd, env?, mode?)`：阻塞执行指定命令：
  - `mode=0`表示直接执行，命令输出直接 Pipe 到标准输出（默认）；
  - `mode=1`表示等待命令执行后返回输出结果；
  - `mode=2`表示输出 Pipe 到标准输出并且在执行完毕后返回输出结果；
- `sh.bgexec(cmd, env?, mode?)`：在后台执行指定命令（非阻塞）；
- `sh.chdir(dir)`或`sh.cd(dir)`：切换工作目录；
- `sh.cwd(dir)`或`sh.pwd(dir)`：取得当前工作目录；

#### SSH 操作

- `ssh.set(name, value)`：设置参数：
  - `name=user`：设置连接用户名，默认为当前用户；
  - `name=port`：设置端口号，默认为`22`；
  - `name=auth`：设置授权方式，`key`表示使用公钥（默认），`password`表示密码；
  - `name=password`：密码，默认空；
  - `name=key`：私钥文件路径，默认为`~/.ssh/id_rsa`；
  - `name=timeout`：连接超时毫秒时间，默认`60000`；
- `ssh.open(host)`：连接到指定主机；
- `ssh.close()`：关闭连接；
- `ssh.setenv(name, value)`：设置环境变量；
- `ssh.exec(cmd, env?, mode?)`：执行命令；

#### 文件操作

- `fs.readdir(dir)`：读取指定目录下的文件列表；
- `fs.readfile(filename)`：读取文件内容；
- `fs.stat(filepath)`：读取文件属性信息；
- `fs.exist(filepath)`：判断文件是否存在；
- `fs.writefile(filename, data)`：覆盖写入文件；
- `fs.appendfile(filename, data)`：追加内容到文件末尾；

#### 路径操作

- `path.join(...paths)`：拼接多个子路径；
- `path.abs(filepath)`：取得绝对路径；
- `path.base(filepath)`：取得文件名；
- `path.ext(filename)`：取得文件扩展名；
- `path.dir(filepath)`：取得路径所在的目录名；

#### 命令行参数操作

- `cli.get(flagname)`：获取指定命令行选项值，支持`-name=value`、`--name=value`两种方式；
- `cli.get(index)`：获取指定索引的命令行参数，从`0`开始；
- `cli.bool(flagname)`：获取指定命令行选项的布尔值，当为`-name=0`、`-name=false`、`-name=f`或不存在时结果为`false`；
- `cli.args()`：获取所有命令行参数；
- `cli.opts()`：获取所有命令行选项；
- `cli.prompt(message?)`：获取用户输入的内容，按`[Enter]`结束输入；

#### HTTP 操作

- `http.timeout(ms)`：设置操作超时毫秒时间，默认为`60000`；
- `http.request(method, url, headers?, body?)`：发送 HTTP 请求；
- `http.download(url, filename?)`：通过 HTTP 下载文件；

#### 日志输出操作

- `log.info(template, ...args)`：输出 INFO 信息（绿色文字）；
- `log.error(template, ...args)`：输出 ERROR 信息（红色文字）；
- `log.fatal(template, ...args)`：输出 FATAL 信息（红色文字）并结束进程；

#### 网络连接操作

- `socket.timeout(ms)`：设置操作超时毫秒时间，默认为`60000`；
- `socket.tcpsend(host, port, data)`：往指定主机端口发送一段数据，并返回结果；
- `socket.tcptest(host, port)`：测试指定主机端口是否可连接；

#### SQL 连接操作

- `sql.set(name, value)`：设置 SQL 连接配置：
  - `name=connMaxLifetime`：最长非活跃毫秒时间，默认`60000`；
- `sql.open(driverName, dataSourceName)`：打开连接：
  - 当`driverName=mysql`时，`dataSourceName`格式：`user:password@tcp(host:port)/database?params`；
  - 暂不支持其他 driver；
- `sql.query(sql, ...args)`：执行查询，并返回结果：
  - 当`driverName=mysql`时，`dataSourceName`需要增加参数`interpolateParams=true`来开启模板参数替换，`args`才生效；
- `sql.exec(sql, ...args)`：执行查询，返回`lastInsertId`和`rowsAffected`；
- `sql.close()`：关闭连接；

## 示例

- **jssh 构建脚本**：[build.js](https://github.com/leizongmin/jssh/blob/main/build.js)；
- **nslookup 包装**：[example/nslookup.js](https://github.com/leizongmin/jssh/blob/main/example/nslookup.js)；

## License

AGPL-3.0
