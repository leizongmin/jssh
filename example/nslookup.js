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
