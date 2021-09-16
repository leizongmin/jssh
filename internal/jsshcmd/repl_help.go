package jsshcmd

import (
	"strings"
)

var replApiList = []string{
	"fs.readdir(", "fs.readfile(", "fs.stat(", "fs.exist(", "fs.writefile(", "fs.appendfile(", "fs.readfilebytes(",
	"path.join(", "path.abs(", "path.base(", "path.ext(", "path.dir(",
	"cli.get(", "cli.bool(", "cli.args(", "cli.opts(", "cli.prompt(", "cli.subcommand(", "cli.subcommandstart(",
	"http.timeout(", "http.request(", "http.download(", "http.get(", "http.head(", "http.options(", "http.post(", "http.put(", "http.delete(",
	"log.debug(", "log.info(", "log.warn(", "log.error(", "log.fatal(",
	"setenv(", "exec(", "exec1(", "exec2", "bgexec(", "chdir(", "cd(", "cwd(", "pwd(",
	"ssh.set(", "ssh.open(", "ssh.close(", "ssh.setenv(", "ssh.exec(", "ssh.exec1(", "ssh.exec2(",
	"socket.timeout(", "socket.tcpsend(", "socket.tcptest(",
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
