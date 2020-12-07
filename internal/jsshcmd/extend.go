package jsshcmd

import "github.com/leizongmin/go/typeutil"

var registeredGlobal typeutil.H

func init() {
	registeredGlobal = make(typeutil.H)
}

// 注册全局变量
func RegisterGlobal(name string, value interface{}) {
	registeredGlobal[name] = value
}
