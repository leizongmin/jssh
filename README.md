# jssh
使用 JavaScript 编写运维脚本。

## 安装

通过 Go 命令行工具安装：

```bash
export GOPROXY=https://goproxy.cn
go get -u github.com/leizongmin/jssh
```

试用v0.1版本：

- macOS(amd64): https://cdn.itoutiao.co/jssh/v0.1/osx/jssh
- Linux(amd64): https://cdn.itoutiao.co/jssh/v0.1/linux/jssh

```bash
curl -o /usr/local/bin/jssh https://cdn.itoutiao.co/jssh/v0.1/linux/jssh && chmod +x /usr/local/bin/jssh
```

## 参考文档

TypeScript types 定义参考文件 [jssh.d.ts](https://github.com/leizongmin/jssh/blob/main/jssh.d.ts)。

## 示例

参考文件 [example/nslookup.js](https://github.com/leizongmin/jssh/blob/main/example/nslookup.js)。

```javascript
#!/usr/bin/env jssh

const host = cli.get(0);
const dns = cli.get(1)

if (!host) {
    println(`使用方法：./nslookup.js <域名> [DNS服务器]`);
    exit(1);
}

const {code, output} = sh.exec(dns ? `nslookup ${host} ${dns}` : `nslookup ${host}`, {}, true);
if (code !== 0) {
    println(output);
    exit(code);
}

let name = ""
let address = ""
const results = []
output.split("\n").slice(2).forEach(line => {
    if (line.startsWith("Name:")) {
        name = line.slice(5).trim();
    } else if (line.startsWith("Address:")) {
        address = line.slice(8).trim();
        if (name) {
            results.push({name, address})
            name = ""
            address = ""
        }
    }
});

println(JSON.stringify(results));
```

## License

AGPL-3.0
