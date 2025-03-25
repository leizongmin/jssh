package completer

// PrefixCompleterInterface 是自动完成接口的抽象
type PrefixCompleterInterface interface {
	// Do 执行自动完成，返回候选项列表
	Do(line []rune, pos int) (newLine [][]rune, length int)
}

// PrefixCompleter 是自动完成的实现
type PrefixCompleter struct {
	Children []PrefixCompleterInterface
}

// Do 执行自动完成，返回候选项列表
func (p *PrefixCompleter) Do(line []rune, pos int) (newLine [][]rune, length int) {
	var candidates [][]rune

	// 遍历所有子项，查找匹配的候选项
	for _, child := range p.Children {
		childNewLine, childLength := child.Do(line, pos)
		if childLength > length {
			length = childLength
			candidates = childNewLine
		} else if childLength == length {
			candidates = append(candidates, childNewLine...)
		}
	}

	return candidates, length
}

// NewPrefixCompleter 创建一个新的前缀自动完成器
func NewPrefixCompleter(items ...PrefixCompleterInterface) *PrefixCompleter {
	return &PrefixCompleter{Children: items}
}

// PcItem 创建一个前缀自动完成项
func PcItem(prefix string) *prefixCompleterItem {
	return &prefixCompleterItem{Prefix: []rune(prefix)}
}

// prefixCompleterItem 是自动完成项的实现
type prefixCompleterItem struct {
	Prefix   []rune
	Children []PrefixCompleterInterface
}

// Do 执行自动完成，返回候选项列表
func (p *prefixCompleterItem) Do(line []rune, pos int) (newLine [][]rune, length int) {
	lineStr := string(line[:pos])
	prefixStr := string(p.Prefix)

	// 如果行前缀匹配当前项
	if len(lineStr) <= len(prefixStr) && prefixStr[:len(lineStr)] == lineStr {
		return [][]rune{p.Prefix}, len(lineStr)
	}

	// 如果行前缀完全匹配当前项，并且有子项，则查找子项
	if len(lineStr) >= len(prefixStr) && lineStr[:len(prefixStr)] == prefixStr {
		for _, child := range p.Children {
			childNewLine, childLength := child.Do(line, pos)
			if childLength > length {
				length = childLength
				newLine = childNewLine
			} else if childLength == length {
				newLine = append(newLine, childNewLine...)
			}
		}
	}

	return
}
