package jsshcmd

import (
	"github.com/leizongmin/jssh/internal/readline/completer"
	"github.com/leizongmin/jssh/quickjs"
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

func replCompleter(jsGlobals quickjs.Value) completer.PrefixCompleterInterface {
	var items []completer.PrefixCompleterInterface
	for _, n := range replApiList {
		items = append(items, completer.PcItem(n))
	}
	if names, err := jsGlobals.PropertyNames(); err == nil {
		for _, n := range names {
			a := jsGlobals.GetByAtom(n.Atom)
			s := n.String()
			if a.IsFunction() {
				s += "("
			}
			items = append(items, completer.PcItem(s))
			a.Free()
		}
	}
	return completer.NewPrefixCompleter(items...)
}
