package chzyer

import (
	"github.com/leizongmin/jssh/internal/readline/completer"
)

// CompleterAdapter 将内部的completer.PrefixCompleterInterface适配为chzyer_readline.PrefixCompleterInterface
type CompleterAdapter struct {
	completer completer.PrefixCompleterInterface
}

// NewCompleterAdapter 创建一个新的CompleterAdapter
func NewCompleterAdapter(completer completer.PrefixCompleterInterface) *CompleterAdapter {
	return &CompleterAdapter{completer: completer}
}

// Do 实现chzyer_readline.PrefixCompleterInterface接口
func (a *CompleterAdapter) Do(line []rune, pos int) (newLine [][]rune, length int) {
	return a.completer.Do(line, pos)
}
