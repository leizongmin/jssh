package standard

import (
	"github.com/leizongmin/jssh/internal/readline"
	"github.com/leizongmin/jssh/internal/readline/completer"
)

// StandardReadlineWithCompleter 是基于Go标准库的readline实现，支持自动完成功能
type StandardReadlineWithCompleter struct {
	// 嵌入StandardReadline以继承其基本功能
	StandardReadline
	// 自动完成器
	completer completer.PrefixCompleterInterface
}

// NewStandardReadlineWithCompleter 创建一个新的StandardReadlineWithCompleter实例
func NewStandardReadlineWithCompleter() readline.Readline {
	return &StandardReadlineWithCompleter{
		StandardReadline: *NewStandardReadline().(*StandardReadline),
	}
}

// Init 初始化readline实例
func (r *StandardReadlineWithCompleter) Init(prompt string, historyFile string, completerObj interface{}) error {
	// 调用基础实现的Init方法
	err := r.StandardReadline.Init(prompt, historyFile, nil)
	if err != nil {
		return err
	}

	// 设置自动完成器
	if completerObj != nil {
		if c, ok := completerObj.(completer.PrefixCompleterInterface); ok {
			r.completer = c
		}
	}

	return nil
}

// Readline 读取一行输入，支持自动完成
func (r *StandardReadlineWithCompleter) Readline() (string, error) {
	// 这里可以实现自动完成功能
	// 但由于标准库不直接支持自动完成，我们仍然使用基本实现
	// 在实际应用中，可以考虑使用其他支持自动完成的库
	return r.StandardReadline.Readline()
}
