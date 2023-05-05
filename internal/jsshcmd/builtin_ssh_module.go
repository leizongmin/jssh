package jsshcmd

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"golang.org/x/crypto/ssh"

	"github.com/leizongmin/jssh/internal/jsexecutor"
	"github.com/leizongmin/jssh/internal/utils"
)

var globalSshClient *ssh.Client
var globalSshEnv map[string]string
var globalSshConfig utils.H

func init() {
	globalSshConfig = utils.H{
		"user":    mustGetCurrentUsername(),
		"auth":    "key",
		"key":     filepath.Join(mustGetHomeDir(), ".ssh/id_rsa"),
		"keypass": "",
		"port":    float64(22),
		"timeout": float64(60_000),
	}
}

func jsFnSshSet(global utils.H) jsexecutor.JSFunction {
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

		if name == "user" && !value.IsString() {
			return ctx.ThrowTypeError("ssh.set: [user] expected string type")
		}
		if name == "auth" && !value.IsString() {
			return ctx.ThrowTypeError("ssh.set: [auth] expected string type")
		}
		if name == "key" && !value.IsString() {
			return ctx.ThrowTypeError("ssh.set: [key] expected string type")
		}
		if name == "keypass" && !value.IsString() {
			return ctx.ThrowTypeError("ssh.set: [keypass] expected string type")
		}
		if name == "port" && !value.IsNumber() {
			return ctx.ThrowTypeError("ssh.set: [port] expected string type")
		}
		if name == "timeout" && !value.IsNumber() {
			return ctx.ThrowTypeError("ssh.set: [timeout] expected string type")
		}

		v, err := jsexecutor.JSValueToAny(value)
		if err != nil {
			return ctx.ThrowError(err)
		}
		globalSshConfig[name] = v

		return ctx.Bool(true)
	}
}

func jsFnSshOpen(global utils.H) jsexecutor.JSFunction {
	return func(ctx *jsexecutor.JSContext, this jsexecutor.JSValue, args []jsexecutor.JSValue) jsexecutor.JSValue {
		if len(args) < 1 {
			return ctx.ThrowSyntaxError("ssh.open: missing host")
		}
		if !args[0].IsString() {
			return ctx.ThrowTypeError("ssh.open: first argument expected string type")
		}
		host := args[0].String()

		if globalSshClient != nil {
			return ctx.ThrowInternalError("ssh.open: please close the previous connection")
		}

		conf := ssh.ClientConfig{
			Timeout: time.Second * 60,
			HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
				return nil
			},
			BannerCallback: func(message string) error {
				return nil
			},
		}

		if user, ok := globalSshConfig["user"].(string); ok {
			conf.User = user
		}
		if auth, ok := globalSshConfig["auth"].(string); ok {
			if auth == "key" {
				key, ok := globalSshConfig["key"].(string)
				if !ok {
					return ctx.ThrowInternalError("ssh.open: [key] missing")
				}
				b, err := ioutil.ReadFile(key)
				if err != nil {
					return ctx.ThrowInternalError("ssh.open: read private key from [key] fail: %s", err)
				}
				if keypass, ok := globalSshConfig["keypass"].(string); ok && len(keypass) > 0 {
					signer, err := ssh.ParsePrivateKeyWithPassphrase(b, []byte(keypass))
					if err != nil {
						return ctx.ThrowInternalError("ssh.open: parse private key from [key] fail: %s", err)
					}
					conf.Auth = []ssh.AuthMethod{ssh.PublicKeys(signer)}
				} else {
					signer, err := ssh.ParsePrivateKey(b)
					if err != nil {
						return ctx.ThrowInternalError("ssh.open: parse private key from [key] fail: %s", err)
					}
					conf.Auth = []ssh.AuthMethod{ssh.PublicKeys(signer)}
				}
			} else if auth == "password" {
				password, ok := globalSshConfig["password"].(string)
				if !ok {
					return ctx.ThrowInternalError("ssh.open: [password] missing")
				}
				conf.Auth = []ssh.AuthMethod{ssh.Password(password)}
			} else {
				return ctx.ThrowInternalError("ssh.open: [auth] only supported one of key,password")
			}
		}
		if timeout, ok := globalSshConfig["timeout"].(float64); ok {
			conf.Timeout = time.Millisecond * time.Duration(int64(timeout))
		}
		port := 22
		if p, ok := globalSshConfig["port"].(float64); ok {
			port = int(p)
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

func jsFnSshClose(global utils.H) jsexecutor.JSFunction {
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

func jsFnSshSetenv(global utils.H) jsexecutor.JSFunction {
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

func jsFnSshExec(global utils.H) jsexecutor.JSFunction {
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
				env2, ok := second.(utils.H)
				if !ok {
					return ctx.ThrowTypeError("ssh.exec: second argument expected an object")
				}
				for n, v := range env2 {
					env[n] = fmt.Sprintf("%v", v)
				}
			}
		}

		saveOutput := false
		pipeOutput := true
		if len(args) >= 3 {
			if !args[2].IsNumber() {
				return ctx.ThrowTypeError("ssh.exec: third argument expected number type")
			}
			mode := args[2].Int32()
			switch mode {
			case 0:
				saveOutput = false
				pipeOutput = true
			case 1:
				saveOutput = true
				pipeOutput = false
			case 2:
				saveOutput = true
				pipeOutput = true
			default:
				return ctx.ThrowTypeError("ssh.exec: mode expected one of 0,1,2")
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
		var code int
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

			var saveBuffer *bytes.Buffer
			if saveOutput {
				saveBuffer = bytes.NewBuffer(nil)
			}

			go func() {
				if _, err := pipeBufferAndSave(os.Stdout, stdout, saveBuffer); err != nil {
					if err != os.ErrClosed {
						stdLog.Printf("ssh.exec: [stdout] %s", err)
					}
				}
				wg.Done()
			}()
			go func() {
				if _, err := pipeBufferAndSave(os.Stderr, stderr, saveBuffer); err != nil {
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
				if err2, ok := err.(*ssh.ExitError); ok {
					code = err2.ExitStatus()
				} else {
					stdLog.Printf("ssh.exec: %s", err)
				}
			}

			if saveBuffer != nil {
				output = saveBuffer.Bytes()
			}
		} else {

			out, err := session.CombinedOutput(cmd)
			if err != nil {
				if err2, ok := err.(*ssh.ExitError); ok {
					code = err2.ExitStatus()
				} else {
					stdLog.Printf("ssh.exec: %s", err)
				}
			}
			output = out
		}

		if err := session.Close(); err != nil {
			if err != io.EOF {
				errLog.Printf("ssh.exec: close session fail: %s", err)
			}
		}

		return jsexecutor.AnyToJSValue(ctx, utils.H{
			"code":        code,
			"output":      string(output),
			"outputbytes": len(output),
		})
	}
}
