package jsshcmd

import (
	"os"

	"github.com/leizongmin/jssh/internal/readline"
	"github.com/leizongmin/jssh/internal/readline/chzyer"
	"github.com/leizongmin/jssh/internal/readline/standard"
)

// ReadlineFactory 创建 ReadlineInterface 实例的工厂函数类型
type ReadlineFactory func() readline.Readline

// 全局 readline 工厂函数，默认使用 chzyer/readline 实现
var DefaultReadlineFactory ReadlineFactory

// InitReadlineFactory 初始化默认的 readline 工厂函数
func InitReadlineFactory() {
	// 默认使用StandardReadline实现
	DefaultReadlineFactory = func() readline.Readline {
		return standard.NewStandardReadline()
	}

	// 如果环境变量JSSH_READLINE=chzyer，则使用ChzyerReadline实现
	if os.Getenv("JSSH_READLINE") == "chzyer" {
		DefaultReadlineFactory = func() readline.Readline {
			return chzyer.NewChzyerReadline()
		}
	}
}
