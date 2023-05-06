# jssh ![CodeQL](https://github.com/leizongmin/jssh/workflows/CodeQL/badge.svg) [![PkgGoDev](https://pkg.go.dev/badge/leizongmin/jssh)](https://pkg.go.dev/github.com/leizongmin/jssh) [![Go Report Card](https://goreportcard.com/badge/leizongmin/jssh)](https://goreportcard.com/report/github.com/leizongmin/jssh) ![GitHub](https://img.shields.io/github/license/leizongmin/jssh)
[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2Fleizongmin%2Fjssh.svg?type=shield)](https://app.fossa.com/projects/git%2Bgithub.com%2Fleizongmin%2Fjssh?ref=badge_shield)

一款使用 JS 编写命令行小工具的神器。

-----

**提示：此项目正在初期开发阶段，接口可能会有较大的调整，且代码测试覆盖率较低。仅供学习研究之用，请谨慎用于生产环境。**

## 简介

- 基于 [QuickJS](https://github.com/bellard/quickjs) 引擎，**支持 [ES2020](https://tc39.github.io/ecma262/) 语言特性**；
- 脚本执行引擎仅需一个**约 7 MB** 的二进制文件，无需安装其他依赖；
- 标准库支持基本的文件操作、执行系统命令、SSH 远程操作、HTTP、Socket、命令行参数解析、日志输出等操作，满足大部分的命令行工具需求；
- 所有操作均为阻塞函数，无异步操作，**简单易于编写**；
- **内存占用小，启动速度快（约 10ms）**；

## 安装

安装预构建的版本（仅支持 Linux amd64、macOS amd64 和 Windows amd64 三种版本）：

[Releases 页面](https://github.com/leizongmin/jssh/releases)

其他系统通过 Go 命令行工具安装最新版本：

```bash
go install github.com/leizongmin/jssh@latest
```

**若要卸载 jssh，直接删除 jssh 二进制文件即可**。可尝试执行以下命令：

```bash
rm $(which jssh)
```


## 命令行工具使用

- 执行脚本：`jssh file.js`；
- 进入 REPL：`jssh repl`；
- 执行命令行参数指定的脚本代码：`jssh exec "js code"`；
- 执行命令行参数指定的脚本代码，并将结果作为字符串输出：`jssh eval "js code"`；
- 将 JS 源码编译为二进制可执行文件（仅限当前操作系统）：`jssh build file.js`；
- 显示帮助信息：`jssh help`；

## 全局配置

- REPL 执行历史默认存储于`~/.jssh_history`文件中，可通过环境变量`JSSH_HISTORY=path`文件设置为新到文件路径，设置`JSSH_HISTORY=0`关闭历史记录；
- jssh 启动时会尝试加载`~/.jssh_global.js`文件和`~/.jssh_global.d`目录下的所有`.js`文件并执行，可通过环境变量`JSSH_GLOBAL=0`关闭；

## 应用示例

- **jssh 构建脚本**：[jssh/build.js](https://github.com/leizongmin/jssh/blob/main/build.js)；
- **实现一个玩具 Docker**：[toy-docker](https://github.com/leizongmin/toy-docker)；

## 参考文档

TypeScript 类型定义参考文件 [jssh.d.ts](https://github.com/leizongmin/jssh/blob/main/jssh.d.ts)，可将其导入编辑器以获得代码自动提示支持。

#### 全局变量列表

- `global`：全局变量；
- `globalThis`：全局变量，同 `global`；
- `__selfcontained`：是否 self-contained 方式执行；
- `__cpucount`：CPU 数量；
- `__os`：操作系统类型，如 `darwin`、`freebsd`、`linux`、`windows`；
- `__arch`：处理器架构，如 `386`、`amd64`、`arm`、`s390x`；
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
- `__output`：上一个命令输出的结果，`sh.exec()` 且 `mode=1` 或 `mode=2` 时有效；
- `__outputbytes`：上一个命令输出结果的字节数；
- `__code`：上一个命令结束时的状态码；
- `__globalfiles`：已自动加载的全局脚本文件列表；

#### 全局函数列表

- `require(modulename)`：加载 CommonJS 模块（**注意：仍然有 Bug，慎用**）；
- `evalfile(filename, content?)`：以 `eval` 方式执行指定脚本文件（本地文件路径或 URL），若指定参数 `content` 则不需要实际读取文件内容；
- `print(template, ...args)`：在 `stdout` 格式化字符串并输出，格式同`format()`函数；
- `println(template, ...args)`：在 `stdout` 格式化字符串并输出，末尾加换行符，格式同 `format()` 函数；
- `eprint(template, ...args)`：在 `stderr` 格式化字符串并输出，格式同`format()`函数；
- `eprintln(template, ...args)`：在 `stderr` 格式化字符串并输出，末尾加换行符，格式同 `format()` 函数；
- `stdoutlog(message)`：在 `stdout` 中输出一行日志；
- `stderrlog(message)`：在 `stderr` 中输出一行日志；
- `readline()`：从控制台获取用户一行的字符串输入；
- `sleep(ms)`：等待指定毫秒时间；
- `exit(code)`：结束进程；
- `loadconfig(filename, format?)`：加载配置文件，支持 JSON、YAML、TOML 格式；
- `networkinterfaces()`：获得网络接口信息；
- `bytesize(data)`：获取指定字符串的字节长度；
- `stdin()`：读取标准输入的内容，返回字符串；
- `stdinbytes()`：读取标准输入的内容，返回 `Uint8Array`；

#### 字符串操作

- `format(template, ...args)`：格式化字符串，如 `format("a=%d, b=%s", 123, "xxx")`；
- `randomstring(size, chars?)`：生成随机字符串；
- `formatdate(format, timestamp?)`：格式化日期时间，格式参考PHP的 `date()` 函数，文档参考 https://locutus.io/php/datetime/date/ ；
- `deepmerge(target, src)`：深度合并两个对象；

#### 编码解码操作

- `base64encode(data)`：Base64 编码字符串；
- `base64decode(data)`：Base64 解码字符串；
- `md5(data)`：MD5 编码字符串；
- `sha1(data)`：SAH1 编码字符串；
- `sha256(data)`：SHA256 编码字符串；

#### Shell 操作

- `setenv(name, value)`：设置环境变量；
- `exec(cmd, env?, mode?)`：阻塞执行指定命令：
  - `mode=0`表示直接执行，命令输出直接管道输出到标准输出 `stdout`（默认）；
  - `mode=1`表示等待命令执行后返回输出结果；
  - `mode=2`表示输出 Pipe 到标准输出并且在执行完毕后返回输出结果；
- `exec1(cmd, env?)`：阻塞执行指定命令，`mode=1`
- `exec2(cmd, env?)`：阻塞执行指定命令，`mode=2`
- `bgexec(cmd, env?, mode?)`：在后台执行指定命令（非阻塞）；
- `pty(cmd, env?)`：模拟 TTY 设备执行指定命令；
- `chdir(dir)`或`cd(dir)`：切换工作目录；
- `cwd(dir)`或`pwd(dir)`：取得当前工作目录；

#### SSH 操作

- `ssh.set(name, value)`：设置参数：
  - `name=user`：设置连接用户名，默认为当前用户；
  - `name=port`：设置端口号，默认为 `22`；
  - `name=auth`：设置授权方式，`key` 表示使用公钥（默认），`password` 表示密码；
  - `name=password`：密码，默认空；
  - `name=key`：私钥文件路径，默认为 `~/.ssh/id_rsa`；
  - `name=timeout`：连接超时毫秒时间，默认 `60000`；
- `ssh.open(host)`：连接到指定主机；
- `ssh.close()`：关闭连接；
- `ssh.setenv(name, value)`：设置环境变量；
- `ssh.exec(cmd, env?, mode?)`：执行命令；
- `ssh.exec1(cmd, env?)`：执行命令，`mode=1`；
- `ssh.exec2(cmd, env?)`：执行命令，`mode=2`；

#### 文件操作

- `fs.readdir(dir)`：读取指定目录下的文件列表；
- `fs.readfile(filename)`：读取文件内容；
- `fs.readfilebytes(filename)`：读取文件内容，返回 `Uint8Array`；
- `fs.stat(filepath)`：读取文件属性信息；
- `fs.exist(filepath)`：判断文件是否存在；
- `fs.writefile(filename, data)`：覆盖写入文件，`data`可以为字符串或 `Uint8Array`；
- `fs.appendfile(filename, data)`：追加内容到文件末尾，`data` 可以为字符串或 `Uint8Array`；
- `fs.mkdir(filepath)`：创建目录，如果目录已存在将返回 `false`；
- `fs.mkdirp(filepath)`：创建目录，如果父目录不存在将一并创建，如果目录已存在将返回 `false`；

#### 路径操作

- `path.join(...paths)`：拼接多个子路径；
- `path.abs(filepath)`：取得绝对路径；
- `path.base(filepath)`：取得文件名；
- `path.ext(filename)`：取得文件扩展名；
- `path.dir(filepath)`：取得路径所在的目录名；

#### 命令行参数操作

- `cli.get(flagname)`：获取指定命令行选项值，支持 `-name=value`、`--name=value` 两种方式；
- `cli.get(index)`：获取指定索引的命令行参数，从 `0` 开始；
- `cli.bool(flagname)`：获取指定命令行选项的布尔值，当为 `-name=0`、`-name=false`、`-name=f` 或不存在时结果为 `false`；
- `cli.args()`：获取所有命令行参数；
- `cli.opts()`：获取所有命令行选项；
- `cli.prompt(message?)`：获取用户输入的内容，按 `[Enter]` 键结束输入；
- `cli.subcommand(name, callback)`：注册子命令处理函数，当 `name=*` 表示其他情况；
- `cli.subcommandstart()`：开始解析执行子命令；

#### HTTP 操作

- `http.timeout(ms)`：设置操作超时毫秒时间，默认为 `60000`；
- `http.request(method, url, headers?, body?)`：发送 HTTP 请求，返回结果包括 `{ status, url, headers, body }`；
- `http.download(url, filename?)`：通过 HTTP 下载文件；
- `http.get(url, headers?)`：发送 HTTP GET 请求，返回结果格式同 `http.request`；
- `http.head(url, headers?)`：发送 HTTP HEAD 请求，返回结果格式同 `http.request`；
- `http.options(url, headers?)`：发送 HTTP OPTIONS 请求，返回结果格式同 `http.request`；
- `http.post(url, headers?, body?)`：发送 HTTP POST 请求，返回结果格式同 `http.request`；
- `http.put(url, headers?, body?)`：发送 HTTP PUT 请求，返回结果格式同 `http.request`；
- `http.delete(url, headers?, body?)`：发送 HTTP DELETE 请求，返回结果格式同 `http.request`；

#### 日志输出操作

可以通过环境变量 `JSSH_LOG=<DEBUG|INFO|ERROR>` 来控制日志输出等级，默认 `JSS_LOG=INFO`。

- `log.debug(template, ...args)`：输出 DEBUG 信息（绿色文字），格式同 `format()` 函数；
- `log.info(template, ...args)`：输出 INFO 信息（绿色文字），格式同 `format()` 函数；
- `log.warn(template, ...args)`：输出 WARN 信息（黄色文字），格式同 `format()` 函数；
- `log.error(template, ...args)`：输出 ERROR 信息（红色文字），格式同 `format()` 函数；
- `log.fatal(template, ...args)`：输出 FATAL 信息（红色文字）并结束进程，格式同 `format()` 函数；

#### 网络连接操作

- `socket.timeout(ms)`：设置操作超时毫秒时间，默认为 `60000`；
- `socket.tcpsend(host, port, data)`：往指定主机端口发送一段数据，并返回结果；
- `socket.tcptest(host, port)`：测试指定主机端口是否可连接；


#### 断言及测试操作

- `assert(ok, message?)`：简单断言；

## 开发

需要安装以下环境：

- Node.js >= 1.16
- Go >= 1.18

执行以下命令构建项目：

```bash
make all
```

## 致谢

- 感谢 [JetBrains](https://www.jetbrains.com/?from=jssh) 提供免费的 Goland 开发工具；
- 感谢 [bellard](https://bellard.org/) 开发的 [QuickJS](https://github.com/bellard/quickjs) 引擎；

## License

```text
MIT License

Copyright (c) 2020-2023 Zongmin Lei <leizongmin@gmail.com>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
```


[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2Fleizongmin%2Fjssh.svg?type=large)](https://app.fossa.com/projects/git%2Bgithub.com%2Fleizongmin%2Fjssh?ref=badge_large)
