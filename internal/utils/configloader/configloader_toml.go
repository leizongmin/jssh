package configloader

import (
	"github.com/BurntSushi/toml"
)

// 支持toml格式的配置文件
func init() {
	RegisterExtensionHandler("toml", tomlLoader)
}

func tomlLoader(data []byte, v interface{}) error {
	return toml.Unmarshal(data, v)
}
