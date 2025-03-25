//go:build chzyer

package jsshcmd

import (
	"github.com/leizongmin/jssh/internal/readline"
	"github.com/leizongmin/jssh/internal/readline/chzyer"
)

// InitReadlineFactory 初始化默认的 readline 工厂函数
func InitReadlineFactory() {
	// 使用ChzyerReadline实现
	DefaultReadlineFactory = func() readline.Readline {
		return chzyer.NewChzyerReadline()
	}
}
