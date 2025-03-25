package jsshcmd

import (
	"github.com/chzyer/readline"
)

// ChzyerReadline 是基于 github.com/chzyer/readline 的实现
type ChzyerReadline struct {
	instance *readline.Instance
}

// NewChzyerReadline 创建一个新的 ChzyerReadline 实例
func NewChzyerReadline() ReadlineInterface {
	return &ChzyerReadline{}
}

// Init 初始化 readline 实例
func (r *ChzyerReadline) Init(prompt string, historyFile string, completer interface{}) error {
	config := &readline.Config{
		Prompt:          prompt,
		HistoryFile:     historyFile,
		InterruptPrompt: "^C",
		EOFPrompt:       "exit",
	}

	// 设置自动完成
	if completer != nil {
		if c, ok := completer.(readline.PrefixCompleterInterface); ok {
			config.AutoComplete = c
		}
	}

	instance, err := readline.NewEx(config)
	if err != nil {
		return err
	}
	r.instance = instance
	return nil
}

// Close 关闭 readline 实例
func (r *ChzyerReadline) Close() error {
	if r.instance != nil {
		return r.instance.Close()
	}
	return nil
}

// Readline 读取一行输入
func (r *ChzyerReadline) Readline() (string, error) {
	return r.instance.Readline()
}

// SaveHistory 保存历史记录
func (r *ChzyerReadline) SaveHistory(line string) error {
	return r.instance.SaveHistory(line)
}

// SetPrompt 设置提示符
func (r *ChzyerReadline) SetPrompt(prompt string) {
	r.instance.SetPrompt(prompt)
}
