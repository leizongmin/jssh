package jsshcmd

import (
	"github.com/leizongmin/jssh/internal/readline"
)

// DefaultReadlineFactory 是创建 readline 实例的工厂函数
var DefaultReadlineFactory func() readline.Readline
