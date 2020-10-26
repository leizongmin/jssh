package jsshcmd

import "github.com/leizongmin/go/typeutil"

var registeredGlobal typeutil.H

func init() {
	registeredGlobal = make(typeutil.H)
}

func RegisterGlobal(name string, value interface{}) {
	registeredGlobal[name] = value
}
