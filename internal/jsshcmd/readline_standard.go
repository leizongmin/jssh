package jsshcmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"sync"
)

// StandardReadline 是基于Go标准库的readline实现
type StandardReadline struct {
	prompt      string
	historyFile string
	scanner     *bufio.Scanner
	history     []string
	historyMu   sync.Mutex
	maxHistory  int
}

// NewStandardReadline 创建一个新的StandardReadline实例
func NewStandardReadline() ReadlineInterface {
	return &StandardReadline{
		maxHistory: 1000, // 默认最多保存1000条历史记录
	}
}

// Init 初始化readline实例
func (r *StandardReadline) Init(prompt string, historyFile string, completer interface{}) error {
	r.prompt = prompt
	r.historyFile = historyFile
	r.scanner = bufio.NewScanner(os.Stdin)

	// 如果提供了历史文件路径，则加载历史记录
	if historyFile != "" {
		if err := r.loadHistory(); err != nil {
			// 仅记录错误，不影响继续使用
			fmt.Fprintf(os.Stderr, "Failed to load history: %v\n", err)
		}
	}

	// 注意：标准库实现不支持自动完成功能
	if completer != nil {
		fmt.Fprintf(os.Stderr, "Warning: StandardReadline does not support auto-completion\n")
	}

	return nil
}

// Close 关闭readline实例
func (r *StandardReadline) Close() error {
	// 保存历史记录
	if r.historyFile != "" {
		return r.saveHistoryToFile()
	}
	return nil
}

// Readline 读取一行输入
func (r *StandardReadline) Readline() (string, error) {
	// 打印提示符
	fmt.Print(r.prompt)

	// 读取一行输入
	if r.scanner.Scan() {
		line := r.scanner.Text()
		return line, nil
	}

	// 检查是否有错误
	if err := r.scanner.Err(); err != nil {
		return "", err
	}

	// 如果没有错误但也没有输入，可能是EOF（Ctrl+D）
	return "", fmt.Errorf("EOF")
}

// SaveHistory 保存历史记录
func (r *StandardReadline) SaveHistory(line string) error {
	if line = strings.TrimSpace(line); line == "" {
		return nil // 不保存空行
	}

	r.historyMu.Lock()
	defer r.historyMu.Unlock()

	// 检查是否与最后一条历史记录重复
	if len(r.history) > 0 && r.history[len(r.history)-1] == line {
		return nil
	}

	// 添加到历史记录
	r.history = append(r.history, line)

	// 如果超过最大历史记录数，删除最早的记录
	if len(r.history) > r.maxHistory {
		r.history = r.history[len(r.history)-r.maxHistory:]
	}

	// 如果设置了历史文件，则保存到文件
	if r.historyFile != "" {
		return r.saveHistoryToFile()
	}

	return nil
}

// SetPrompt 设置提示符
func (r *StandardReadline) SetPrompt(prompt string) {
	r.prompt = prompt
}

// 加载历史记录文件
func (r *StandardReadline) loadHistory() error {
	r.historyMu.Lock()
	defer r.historyMu.Unlock()

	// 检查文件是否存在
	file, err := os.Open(r.historyFile)
	if err != nil {
		if os.IsNotExist(err) {
			return nil // 文件不存在，不是错误
		}
		return err
	}
	defer file.Close()

	// 读取历史记录
	scanner := bufio.NewScanner(file)
	r.history = make([]string, 0, r.maxHistory)
	for scanner.Scan() {
		line := scanner.Text()
		if line = strings.TrimSpace(line); line != "" {
			r.history = append(r.history, line)
		}

		// 如果超过最大历史记录数，停止读取
		if len(r.history) >= r.maxHistory {
			break
		}
	}

	return scanner.Err()
}

// 保存历史记录到文件
func (r *StandardReadline) saveHistoryToFile() error {
	// 创建目录（如果不存在）
	dir := strings.TrimSuffix(r.historyFile, "/")
	dir = dir[:strings.LastIndex(dir, "/")]
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	// 创建或打开文件
	file, err := os.Create(r.historyFile)
	if err != nil {
		return err
	}
	defer file.Close()

	// 写入历史记录
	for _, line := range r.history {
		if _, err := fmt.Fprintln(file, line); err != nil {
			return err
		}
	}

	return nil
}
