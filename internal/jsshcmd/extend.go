package jsshcmd

import "github.com/leizongmin/jssh/internal/utils"

var registeredGlobal utils.H

func init() {
	registeredGlobal = make(utils.H)
}

// RegisterGlobal 注册全局变量
func RegisterGlobal(name string, value interface{}) {
	registeredGlobal[name] = value
}
