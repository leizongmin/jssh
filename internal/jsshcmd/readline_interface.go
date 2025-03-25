package jsshcmd

import (
	"os"
)

// ReadlineInterface 定义了 readline 功能的抽象接口
type ReadlineInterface interface {
	// Init 初始化 readline 实例
	Init(prompt string, historyFile string, completer interface{}) error
	// Close 关闭 readline 实例
	Close() error
	// Readline 读取一行输入
	Readline() (string, error)
	// SaveHistory 保存历史记录
	SaveHistory(line string) error
	// SetPrompt 设置提示符
	SetPrompt(prompt string)
}

// ReadlineFactory 创建 ReadlineInterface 实例的工厂函数类型
type ReadlineFactory func() ReadlineInterface

// 全局 readline 工厂函数，默认使用 chzyer/readline 实现
var DefaultReadlineFactory ReadlineFactory

// 初始化默认的 readline 工厂函数
func init() {
	// 默认使用ChzyerReadline实现
	DefaultReadlineFactory = NewChzyerReadline

	// 如果环境变量JSSH_READLINE=standard，则使用StandardReadline实现
	if os.Getenv("JSSH_READLINE") == "standard" {
		DefaultReadlineFactory = NewStandardReadline
	}
}
