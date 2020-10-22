package jsshcmd

import (
	"fmt"
	"github.com/leizongmin/go/typeutil"
	"github.com/leizongmin/jssh/internal/jsexecutor"
	"golang.org/x/crypto/ssh"
	"io"
	"io/ioutil"
	"net"
	"os"
	"strings"
	"sync"
	"time"
)

var globalSshClient *ssh.Client
var globalSshEnv map[string]string

func JsFnSshSet(global typeutil.H) jsexecutor.JSFunction {
	return func(ctx *jsexecutor.JSContext, this jsexecutor.JSValue, args []jsexecutor.JSValue) jsexecutor.JSValue {
		if len(args) < 1 {
			return ctx.ThrowSyntaxError("ssh.set: missing name")
		}
		if !args[0].IsString() {
			return ctx.ThrowTypeError("ssh.set: first argument expected string type")
		}
		name := args[0].String()

		if len(args) < 2 {
			return ctx.ThrowSyntaxError("ssh.set: missing value")
		}
		value := args[1]

		config := ctx.Globals().Get("__ssh_config")
		if !config.IsObject() || config.IsNull() {
			return ctx.ThrowInternalError("ssh.set: __ssh_config expected an object")
		}

		config.Set(name, value)
		ctx.Globals().Set("__ssh_config", config)
		return ctx.Bool(true)
	}
}

func JsFnSshOpen(global typeutil.H) jsexecutor.JSFunction {
	return func(ctx *jsexecutor.JSContext, this jsexecutor.JSValue, args []jsexecutor.JSValue) jsexecutor.JSValue {
		if len(args) < 1 {
			return ctx.ThrowSyntaxError("ssh.open: missing host")
		}
		if !args[0].IsString() {
			return ctx.ThrowTypeError("ssh.open: first argument expected string type")
		}
		host := args[0].String()

		conf := ssh.ClientConfig{
			Timeout: time.Second * 60,
			HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
				return nil
			},
			BannerCallback: func(message string) error {
				return nil
			},
		}

		sshConfig := ctx.Globals().Get("__ssh_config")
		if !sshConfig.IsObject() || sshConfig.IsNull() {
			return ctx.ThrowInternalError("ssh.open: __ssh_config expected an object")
		}
		if sshConfig.Get("user").IsString() {
			conf.User = sshConfig.Get("user").String()
		}
		if sshConfig.Get("auth").IsString() {
			auth := strings.ToLower(sshConfig.Get("auth").String())
			if auth == "key" {
				if !sshConfig.Get("key").IsString() {
					return ctx.ThrowInternalError("ssh.open: __ssh_config.key missing")
				}
				key := sshConfig.Get("key").String()
				b, err := ioutil.ReadFile(key)
				if err != nil {
					return ctx.ThrowInternalError("ssh.open: read private key from __ssh_config.key fail: %s", err)
				}
				if sshConfig.Get("keypass").IsString() && sshConfig.Get("keypass").Len() > 0 {
					keypass := sshConfig.Get("keypass").String()
					signer, err := ssh.ParsePrivateKeyWithPassphrase(b, []byte(keypass))
					if err != nil {
						return ctx.ThrowInternalError("ssh.open: parse private key from __ssh_config.key fail: %s", err)
					}
					conf.Auth = []ssh.AuthMethod{ssh.PublicKeys(signer)}
				} else {
					signer, err := ssh.ParsePrivateKey(b)
					if err != nil {
						return ctx.ThrowInternalError("ssh.open: parse private key from __ssh_config.key fail: %s", err)
					}
					conf.Auth = []ssh.AuthMethod{ssh.PublicKeys(signer)}
				}
			} else if auth == "password" {
				if !sshConfig.Get("password").IsString() {
					return ctx.ThrowInternalError("ssh.open: __ssh_config.password missing")
				}
				password := sshConfig.Get("password").String()
				conf.Auth = []ssh.AuthMethod{ssh.Password(password)}
			} else {
				return ctx.ThrowInternalError("ssh.open: __ssh_config.auth only supported one of key,password")
			}
		}
		if sshConfig.Get("timeout").IsNumber() {
			timeout := sshConfig.Get("timeout").Int64()
			conf.Timeout = time.Millisecond * time.Duration(timeout)
		}
		port := 22
		if sshConfig.Get("port").IsNumber() {
			port = int(sshConfig.Get("port").Int32())
		}

		addr := fmt.Sprintf("%s:%d", host, port)
		client, err := ssh.Dial("tcp", addr, &conf)
		if err != nil {
			return ctx.ThrowError(err)
		}
		globalSshClient = client
		globalSshEnv = make(map[string]string)

		return ctx.Bool(true)
	}
}

func JsFnSshClose(global typeutil.H) jsexecutor.JSFunction {
	return func(ctx *jsexecutor.JSContext, this jsexecutor.JSValue, args []jsexecutor.JSValue) jsexecutor.JSValue {
		if globalSshClient != nil {
			if err := globalSshClient.Close(); err != nil {
				errLog.Printf("ssh.close: close ssh client fail: %s", err)
			}
		}
		globalSshClient = nil

		return ctx.Bool(true)
	}
}

func JsFnSshSetenv(global typeutil.H) jsexecutor.JSFunction {
	return func(ctx *jsexecutor.JSContext, this jsexecutor.JSValue, args []jsexecutor.JSValue) jsexecutor.JSValue {
		if len(args) < 1 {
			return ctx.ThrowSyntaxError("ssh.setenv: missing env name")
		}
		if !args[0].IsString() {
			return ctx.ThrowTypeError("ssh.setenv: first argument expected string type")
		}
		name := args[0].String()

		if len(args) < 2 {
			return ctx.ThrowSyntaxError("ssh.setenv: missing env value")
		}
		if !args[1].IsString() {
			return ctx.ThrowTypeError("ssh.setenv: second argument expected string type")
		}
		value := args[1].String()

		if globalSshEnv == nil {
			return ctx.ThrowInternalError("ssh.exec: please open a session")
		}

		globalSshEnv[name] = value
		return ctx.Bool(true)
	}
}

func JsFnSshExec(global typeutil.H) jsexecutor.JSFunction {
	return func(ctx *jsexecutor.JSContext, this jsexecutor.JSValue, args []jsexecutor.JSValue) jsexecutor.JSValue {
		if len(args) < 1 {
			return ctx.ThrowSyntaxError("ssh.exec: missing exec command")
		}
		if !args[0].IsString() {
			return ctx.ThrowTypeError("ssh.exec: first argument expected string type")
		}
		cmd := args[0].String()

		env := make(map[string]string)
		for n, v := range globalSshEnv {
			env[n] = v
		}
		if len(args) >= 2 {
			if args[1].IsNull() || args[1].IsUndefined() {
			} else {
				if !args[1].IsObject() {
					return ctx.ThrowTypeError("ssh.exec: second argument expected an object")
				}
				second, err := jsexecutor.JSValueToAny(args[1])
				if err != nil {
					return ctx.ThrowError(err)
				}
				env2, ok := second.(typeutil.H)
				if !ok {
					return ctx.ThrowTypeError("ssh.exec: second argument expected an object")
				}
				for n, v := range env2 {
					env[n] = fmt.Sprintf("%v", v)
				}
			}
		}

		pipeOutput := true
		if len(args) >= 3 {
			if !args[2].IsBool() {
				return ctx.ThrowTypeError("ssh.exec: third argument expected boolean type")
			}
			if args[2].Bool() {
				pipeOutput = false
			}
		}

		if globalSshClient == nil {
			return ctx.ThrowInternalError("ssh.exec: please open a connection")
		}

		session, err := globalSshClient.NewSession()
		if err != nil {
			return ctx.ThrowError(err)
		}
		envLines := make([]string, 0)
		for n, v := range env {
			envLines = append(envLines, fmt.Sprintf("export %s=%s", n, v))
		}
		cmd = strings.Join(envLines, "\n") + "\n" + cmd

		var output []byte
		if pipeOutput {
			session.Stdin = os.Stdin
			stdout, err := session.StdoutPipe()
			if err != nil {
				return ctx.ThrowError(err)
			}
			stderr, err := session.StderrPipe()
			if err != nil {
				return ctx.ThrowError(err)
			}

			var wg sync.WaitGroup
			wg.Add(2)

			go func() {
				if _, err := io.Copy(os.Stdout, stdout); err != nil {
					if err != os.ErrClosed {
						stdLog.Printf("ssh.exec: [stdout] %s", err)
					}
				}
				wg.Done()
			}()
			go func() {
				if _, err := io.Copy(os.Stderr, stderr); err != nil {
					if err != os.ErrClosed {
						stdLog.Printf("ssh.exec: [stderr] %s", err)
					}
				}
				wg.Done()
			}()

			if err := session.Start(cmd); err != nil {
				return ctx.ThrowError(err)
			}
			wg.Wait()

			if err := session.Wait(); err != nil {
				stdLog.Printf("ssh.exec: %s", err)
			}
		} else {

			out, err := session.CombinedOutput(cmd)
			if err != nil {
				stdLog.Printf("ssh.exec: %s", err)
			}
			output = out
		}

		if err := session.Close(); err != nil {
			if err != io.EOF {
				errLog.Printf("ssh.exec: close session fail: %s", err)
			}
		}

		return jsexecutor.AnyToJSValue(ctx, typeutil.H{
			"output":      string(output),
			"outputbytes": len(output),
		})
	}
}
