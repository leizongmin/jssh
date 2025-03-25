//go:build !standard && !chzyer

package jsshcmd

import (
	"github.com/leizongmin/jssh/internal/readline"
	"github.com/leizongmin/jssh/internal/readline/standard"
)

// InitReadlineFactory 初始化默认的 readline 工厂函数
func InitReadlineFactory() {
	// 默认使用StandardReadline实现
	DefaultReadlineFactory = func() readline.Readline {
		return standard.NewStandardReadline()
	}
}
