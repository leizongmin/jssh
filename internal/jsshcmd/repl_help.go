package jsshcmd

import (
	"strings"
)

var replApiList = []string{
	"__version", "__bin", "__pid", "__tmpdir", "__homedir", "__user", "__hostname",
	"__dirname", "__filename", "__args", "__env", "__output", "__outputbytes", "__code",
	"set(", "get(", "format(", "print(", "println(", "stdoutlog(", "stderrlog(", "readline(", "sleep(", "exit(", "loadconfig(",
	"base64encode(", "base64decode(", "md5(", "sha1(", "sha256(",
	"networkinterfaces(",
	"fs.readdir(", "fs.readfile(", "fs.stat(", "fs.exist(", "fs.writefile(", "fs.appendfile(",
	"path.join(", "path.abs(", "path.base(", "path.ext(", "path.dir(",
	"cli.get(", "cli.bool(", "cli.args(", "cli.opts(", "cli.prompt(", "cli.subcommand(", "cli.subcommandstart(",
	"http.timeout(", "http.request(", "http.download(",
	"log.debug(", "log.info(", "log.error(", "log.fatal(",
	"sh.setenv(", "sh.exec(", "sh.bgexec(", "sh.chdir(", "sh.cd(", "sh.cwd(", "sh.pwd(",
	"ssh.set(", "ssh.open(", "ssh.close(", "ssh.setenv(", "ssh.exec(",
	"socket.timeout(", "socket.tcpsend(", "socket.tcptest(",
	"sql.set(", "sql.open(", "sql.close(", "sql.query(", "sql.exec(",
}

func replCompleter(line string) (c []string) {
	line = strings.ToLower(line)
	for _, n := range replApiList {
		if strings.HasPrefix(n, line) {
			c = append(c, n)
		}
	}
	return c
}
