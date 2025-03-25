package readline

// Readline 定义了 readline 功能的抽象接口
type Readline interface {
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
type ReadlineFactory func() Readline
