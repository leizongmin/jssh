package completer

// 这个文件提供了与chzyer/readline包兼容的接口和函数
// 用于在不直接依赖chzyer/readline包的情况下提供自动完成功能

// PcItemDynamic 创建一个动态前缀自动完成项
// 与chzyer/readline包的PcItemDynamic函数兼容
func PcItemDynamic(callback func(string) []string) *dynamicPrefixCompleterItem {
	return &dynamicPrefixCompleterItem{Callback: callback}
}

// dynamicPrefixCompleterItem 是动态自动完成项的实现
type dynamicPrefixCompleterItem struct {
	Callback func(string) []string
}

// Do 执行自动完成，返回候选项列表
func (p *dynamicPrefixCompleterItem) Do(line []rune, pos int) (newLine [][]rune, length int) {
	lineStr := string(line[:pos])

	// 调用回调函数获取候选项
	candidates := p.Callback(lineStr)

	// 转换候选项为[][]rune
	for _, candidate := range candidates {
		newLine = append(newLine, []rune(candidate))
	}

	return newLine, len(lineStr)
}
